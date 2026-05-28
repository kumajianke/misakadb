import json
import sys
import threading
import time
from typing import Optional, Dict, Any, List, Union


from network import sock


class HeartbeatController:
    def __init__(self, owner: "MisakaDBClient", interval: float) -> None:
        self._owner = owner
        self.interval = interval
        self._enabled = threading.Event()
        self._shutdown = threading.Event()
        self._stats_lock = threading.Lock()
        self._thread: Optional[threading.Thread] = None
        self._attempt_count = 0
        self._success_count = 0
        self._failure_count = 0
        self._last_error: Optional[str] = None
        self._last_sent_at: Optional[float] = None

    def _ensure_thread(self) -> None:
        if self._thread is not None and self._thread.is_alive():
            return

        self._shutdown.clear()
        self._thread = threading.Thread(
            target=self._loop,
            name=f"MisakaDBHeartbeat-{self._owner.host}:{self._owner.port}",
            daemon=True,
        )
        self._thread.start()

    def _loop(self) -> None:
        while not self._shutdown.is_set():
            if not self._enabled.is_set():
                if self._shutdown.wait(0.1):
                    return
                continue

            if self._owner.connected and self._owner.client is not None:
                try:
                    with self._owner._socket_lock:
                        if self._owner.connected and self._owner.client is not None:
                            self._owner.client.send_heartbeat()
                            print("send")
                            with self._stats_lock:
                                self._attempt_count += 1
                                self._success_count += 1
                                self._last_error = None
                                self._last_sent_at = time.time()
                except Exception as e:
                    with self._stats_lock:
                        self._attempt_count += 1
                        self._failure_count += 1
                        self._last_error = str(e)
                    print(f"心跳发送失败: {e}", file=sys.stderr)
                    self._owner.connected = False

            if self._shutdown.wait(self.interval):
                return

    @property
    def running(self) -> bool:
        return self._enabled.is_set()

    @property
    def count(self) -> int:
        with self._stats_lock:
            return self._attempt_count

    @property
    def success_count(self) -> int:
        with self._stats_lock:
            return self._success_count

    @property
    def failure_count(self) -> int:
        with self._stats_lock:
            return self._failure_count

    @property
    def loss_rate(self) -> float:
        with self._stats_lock:
            if self._attempt_count == 0:
                return 0.0
            return self._failure_count / self._attempt_count * 100.0

    @property
    def last_error(self) -> Optional[str]:
        with self._stats_lock:
            return self._last_error

    @property
    def last_sent_at(self) -> Optional[float]:
        with self._stats_lock:
            return self._last_sent_at

    def stats(self) -> Dict[str, Any]:
        with self._stats_lock:
            loss_rate = (
                self._failure_count / self._attempt_count * 100.0
                if self._attempt_count > 0 else 0.0
            )
            return {
                "running": self._enabled.is_set(),
                "interval": self.interval,
                "count": self._attempt_count,
                "success_count": self._success_count,
                "failure_count": self._failure_count,
                # 当前协议没有心跳回包确认，这里按发送失败率统计。
                "loss_rate": loss_rate,
                "last_error": self._last_error,
                "last_sent_at": self._last_sent_at,
            }

    def start(self) -> None:
        self._ensure_thread()
        self._enabled.set()

    def stop(self) -> None:
        self._enabled.clear()

    def shutdown(self) -> None:
        self._enabled.clear()
        self._shutdown.set()


class MisakaDBClient:
    """MisakaDB客户端API"""
    
    def __init__(
        self,
        host: str = "127.0.0.1",
        port: int = 10032,
        heartbeat_interval: float = 1
    ):
        self.host = host
        self.port = port
        self.client: Optional[sock.clientCore] = None
        self.service_info: Optional[Dict[str, Any]] = None
        self.connected = False
        self.heartbeat_interval = heartbeat_interval
        self._socket_lock = threading.Lock()
        self.heart = HeartbeatController(self, heartbeat_interval)
        self.heart.start()
        
    def connect(self, retries: int = 3, retry_delay: float = 1.0) -> bool:
        """连接到MisakaDB服务器
        
        Args:
            retries: 重试次数
            retry_delay: 重试延迟（秒）
            
        Returns:
            bool: 连接是否成功
        """
        self.heart.start()

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
            with self._socket_lock:
                if not self.connected or self.client is None:
                    return None
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
            with self._socket_lock:
                if not self.connected or self.client is None:
                    return None
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
            with self._socket_lock:
                if not self.connected or self.client is None:
                    return False
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
        self.heart.shutdown()
        with self._socket_lock:
            if self.client is not None and self.client.s:
                self.client.s.close()
            self.client = None
            self.service_info = None
            self.connected = False
        self.heart = HeartbeatController(self, self.heartbeat_interval)
    
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
