> [!NOTE]
> 便捷简单的文档数据库，始终如一。

`MisakaDB` 是一个轻量级的JSON文档数据库，支持`JSON`格式的内容字段存储，在0.0.3版本，仅支持单机模式。

## 启动服务

### 源码启动

服务端的代码由GO语言编写，需要安装go语言服务后，执行命令：

```python
go run misaka.go [参数] 
```

### 程序启动

程序支持多平台编译：

```python
go build misaka
```

然后在命令行启动编译后的程序即可。

### 配置文件启动

在 0.0.3 版本，yaml的内容为：

```python
network:
  port: 8080
  address: 0.0.0.0
  max_conn: 100000
  retry_count: 3
  retry_delay: 30

service:
  version: "0.0.3"

```

>  **| network**
>
> port 端口
>
> address 地址
>
> max_conn 最大服务连接数
>
> retry_count 错误连接重试次数
>
> retry_delay 心跳、重试deadline时长
>
>  **| service（非专业人士请勿修改）**
>
> all_command 服务端允许的所有的命令
>
> version 版本号

> [!NOTE]
> 在0.0.3版本，命令行模式，支持如下参数启动：
>
> - **`address`**     **启动的地址簇**
> - **`port`**            **启动的端口**
> - `configs`     **支持启动的配置【默认misaka.yaml】**
> - `debug`         **调试程序的参数 上线前确保是false 会在6060端口启动一个**`pprof`**服务**

## 客户端使用

客户端由`python`开发，源码使用命令启动：

```python
uv run main.py
```

打包的客户端启动的命令：

```python
main.exe
```