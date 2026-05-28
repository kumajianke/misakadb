import json
import sys
import time
from typing import Optional, Dict, Any, List, Union


from network import sock


class MisakaDBClient:
    """MisakaDB客户端API"""
    
    def __init__(self, host: str = "127.0.0.1", port: int = 10032):
        self.host = host
        self.port = port
        self.client: Optional[sock.clientCore] = None
        self.service_info: Optional[Dict[str, Any]] = None
        self.connected = False
        
    def connect(self, retries: int = 3, retry_delay: float = 1.0) -> bool:
        """连接到MisakaDB服务器
        
        Args:
            retries: 重试次数
            retry_delay: 重试延迟（秒）
            
        Returns:
            bool: 连接是否成功
        """
        for attempt in range(retries):
            try:
                self.client = sock.clientCore(self.host, self.port)
                connect_err = self.client.connect()
                
                if connect_err is not None:
                    if attempt < retries - 1:
                        print(f"连接失败，第{attempt + 1}次重试... ({connect_err})", file=sys.stderr)
                        time.sleep(retry_delay)
                        continue
                    else:
                        print(f"连接失败: {connect_err}", file=sys.stderr)
                        return False
                        
                self.connected = True
                return True
                
            except Exception as e:
                if attempt < retries - 1:
                    print(f"连接过程中发生错误，第{attempt + 1}次重试... ({e})", file=sys.stderr)
                    time.sleep(retry_delay)
                else:
                    print(f"连接过程中发生错误: {e}", file=sys.stderr)
                    return False
        return False
    
    def get_service_info(self) -> Optional[Dict[str, Any]]:
        """获取服务器信息
        
        Returns:
            Optional[Dict[str, Any]]: 服务信息字典，如果失败返回None
        """
        if not self.connected or self.client is None:
            print("未连接到服务器，请先调用connect()", file=sys.stderr)
            return None
            
        try:
            response = self.client.send_command("get-service-info")
            self.service_info = json.loads(response)
            return self.service_info
        except json.JSONDecodeError as e:
            print(f"解析服务信息失败: {e}", file=sys.stderr)
            return None
        except Exception as e:
            print(f"获取服务信息失败: {e}", file=sys.stderr)
            return None
    
    def execute_command(self, command: str) -> Optional[Union[str, Dict, List]]:
        """执行自定义命令
        
        Args:
            command: 要执行的命令字符串
            
        Returns:
            Optional[Union[str, Dict, List]]: 命令执行结果
        """
        if not self.connected or self.client is None:
            print("未连接到服务器，请先调用connect()", file=sys.stderr)
            return None
            
        try:
            response = self.client.send_command(command)
            
            # 尝试解析为JSON
            try:
                return json.loads(response)
            except json.JSONDecodeError:
                # 如果不是JSON，返回原始字符串
                return response.decode('utf-8', errors='replace')
                
        except Exception as e:
            print(f"执行命令失败: {e}", file=sys.stderr)
            return None
    
    def ping(self) -> bool:
        """检查服务器是否可达
        
        Returns:
            bool: 服务器是否可达
        """
        if not self.connected or self.client is None:
            return False
            
        try:
            # 发送一个简单的ping命令或使用心跳
            response = self.client.send_command("ping")
            return response is not None
        except:
            return False
    
    def get_server_version(self) -> Optional[str]:
        """获取服务器版本
        
        Returns:
            Optional[str]: 服务器版本号
        """
        service_info = self.get_service_info()
        if service_info:
            return service_info.get('service', {}).get('version')
        return None
    
    def get_allowed_commands(self) -> Optional[List[str]]:
        """获取服务器允许的命令列表
        
        Returns:
            Optional[List[str]]: 允许的命令列表
        """
        service_info = self.get_service_info()
        if service_info:
            return service_info.get('service', {}).get('allow_command', [])
        return None
    
    def is_command_allowed(self, command: str) -> bool:
        """检查命令是否被服务器允许
        
        Args:
            command: 要检查的命令
            
        Returns:
            bool: 命令是否被允许
        """
        allowed_commands = self.get_allowed_commands()
        if allowed_commands:
            return command in allowed_commands
        return False
    
    def get_network_config(self) -> Optional[Dict[str, Any]]:
        """获取网络配置
        
        Returns:
            Optional[Dict[str, Any]]: 网络配置信息
        """
        service_info = self.get_service_info()
        if service_info:
            return service_info.get('network', {})
        return None
    
    def close(self) -> None:
        """关闭连接"""
        if self.client is not None and self.client.s:
            self.client.s.close()
            self.client = None
            self.service_info = None
            self.connected = False
    
    def __enter__(self):
        """上下文管理器入口"""
        if not self.connect():
            raise ConnectionError(f"无法连接到服务器 {self.host}:{self.port}")
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        """上下文管理器退出"""
        self.close()
    
    def __repr__(self) -> str:
        """对象表示"""
        status = "已连接" if self.connected else "未连接"
        return f"<MisakaDBClient {self.host}:{self.port} [{status}]>"


# 兼容旧版本的函数接口
def connect(host: str, port: int) -> Optional[MisakaDBClient]:
    """连接到MisakaDB服务器（兼容旧版本接口）
    
    Args:
        host: 服务器地址
        port: 服务器端口
        
    Returns:
        Optional[MisakaDBClient]: 客户端实例，如果连接失败返回None
    """
    client = MisakaDBClient(host, port)
    if client.connect():
        return client
    return None


# 示例使用
if __name__ == "__main__":
    # 使用上下文管理器（推荐）
    print("=== 使用上下文管理器 ===")
    try:
        with MisakaDBClient("127.0.0.1", 10032) as client:
            # 获取服务信息
            service_info = client.get_service_info()
            if service_info:
                print("服务信息获取成功:")
                print(json.dumps(service_info, indent=2, ensure_ascii=False))
                
                # 获取特定信息
                version = client.get_server_version()
                allowed_commands = client.get_allowed_commands()
                network_config = client.get_network_config()
                
                print(f"\n服务器版本: {version}")
                print(f"允许的命令: {allowed_commands}")
                print(f"网络配置: {network_config}")
                
                # 检查ping
                if client.ping():
                    print("服务器可达")
                else:
                    print("服务器不可达")
                    
                # 检查命令是否允许
                test_command = "get-service-info"
                if client.is_command_allowed(test_command):
                    print(f"命令 '{test_command}' 被允许")
                else:
                    print(f"命令 '{test_command}' 不被允许")
                    
                # 执行自定义命令
                result = client.execute_command("get-service-info")
                print(f"\n执行命令结果: {result}")
    except ConnectionError as e:
        print(f"连接失败: {e}")
    
    print("\n=== 使用传统方式 ===")
    client = MisakaDBClient("127.0.0.1", 10032)
    if client.connect():
        try:
            service_info = client.get_service_info()
            if service_info:
                print(f"服务版本: {service_info.get('service', {}).get('version', '未知')}")
                print(f"监听地址: {service_info.get('network', {}).get('address', '未知')}")
                print(f"监听端口: {service_info.get('network', {}).get('port', '未知')}")
        finally:
            client.close()
    else:
        print("连接失败")
