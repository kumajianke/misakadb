# net.Conn 与 MAC 地址说明

## 结论

`net.Conn` 不能直接获取对端 MAC 地址。

你通过 `net.Conn` 通常只能拿到：

- 本地地址：`conn.LocalAddr()`
- 对端地址：`conn.RemoteAddr()`

这两个地址本质上是：

- IP
- 端口

不是 MAC 地址。

## 为什么拿不到

MAC 地址是二层链路层地址，只在本地局域网链路中有效。

而 `net.Conn` 是传输层/应用层抽象，它屏蔽了下面这些细节：

- 以太网帧
- ARP
- 网卡层信息

所以：

- 对于跨网段、跨路由的连接，服务端通常根本拿不到客户端真实 MAC
- 即使在同一局域网，也不是通过 `net.Conn` 直接获取

## 你能从 net.Conn 拿到什么

典型方式是：

```go
remoteAddr := conn.RemoteAddr().String()
localAddr := conn.LocalAddr().String()
```

如果是 TCP 连接，通常还能断言成：

```go
tcpAddr, ok := conn.RemoteAddr().(*net.TCPAddr)
if ok {
    ip := tcpAddr.IP.String()
    port := tcpAddr.Port
    _, _ = ip, port
}
```

这里拿到的是远端 IP 和端口，不是 MAC。

## 什么时候可能拿到 MAC

只有在下面这种场景才有可能间接拿到：

- 你和目标在同一个二层网络
- 本机 ARP 缓存里正好有这个 IP 对应的 MAC
- 你额外查询操作系统 ARP 表

这也不是 `net.Conn` 的能力，而是“先从连接拿 IP，再去系统网络栈查 ARP 缓存”。

## 在服务端项目里要注意什么

如果你的目标是：

- 标识客户端身份
- 防止伪造
- 做设备绑定

那不要依赖 MAC 地址，因为：

- 跨路由场景拿不到
- 很容易变化
- 也可能被伪造

更适合的方案通常是：

- 用户名/密码
- 设备 ID
- 证书
- token / session
- 应用层签名

## 对 MisakaDB 的建议

对于当前项目，更建议用这些标识客户端：

- `LoginUser`
- 连接来源 IP
- 后续可补充会话 ID
- 更进一步可补充客户端证书或设备标识

不要把 MAC 作为核心身份依据。
