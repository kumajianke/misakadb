# 07. 鉴权、安全与 misaka-tools

## 关键文件与安全边界

MisakaDB 的鉴权与敏感信息依赖 `profiles/`：

- `profiles/master.mikey`：32 字节 AES 密钥（权限：0400）
- `profiles/user.dat`：加密后的用户表（权限：0600）

服务端在初始化 RegisterCenter 时会强依赖 `profiles/master.mikey`，缺失会 panic：

- [NewRegisterCenter](file:///workspace/network/RegisterCenter/core.go#L38-L43)

## 安全模块 safe（AES-GCM）

实现：[AES256.go](file:///workspace/safe/AES256.go)

- `InitPassword()`：生成随机 32 字节密钥并写入 `profiles/master.mikey`
- `EncryptByte/DecryptByte()`：通过 `RegisterCenter.MasterKey` 获取 key 来加解密

## 用户模块 miusers

实现：[userManager.go](file:///workspace/miusers/userManager.go)

机制：

- 明文密码不落盘
- 用户文件中保存 bcrypt hash（`bcrypt.GenerateFromPassword`）
- 整个用户表 JSON 再通过 AES-GCM 加密后写入 `profiles/user.dat`

常用能力：

- `VerifyPassword(username, password)`
- `VerifyRole(username, role)`
- `ChangePassword/ChangeRole/RemoveUser`

## misaka-tools（开发/运维工具）

入口：[misaka-tools.go](file:///workspace/tools/misaka-tools.go)  
命令实现：[CommandExecute](file:///workspace/command/ToolsCommands/CommandExecute.go#L17-L167)

### sys-init（首次初始化）

用途：初始化 `profiles/master.mikey` 与 `profiles/user.dat`，并创建 root 用户。

启动示例：

```bash
go run ./tools/misaka-tools.go sys-init
```

输出会给出 root 初始随机密码（只出现一次），随后建议立即修改 root 密码：

```bash
go run ./tools/misaka-tools.go chpwd root
```

### 用户管理相关

- 添加用户：
  - `go run ./tools/misaka-tools.go add-user <username> <password>`
- 修改某用户密码（需要 root 鉴权）：
  - `go run ./tools/misaka-tools.go chpwd <username>`
- 修改角色（需要 root 鉴权）：
  - `go run ./tools/misaka-tools.go chmod <username> <role>`
- 删除用户（不允许删除 root，需要 root 鉴权）：
  - `go run ./tools/misaka-tools.go remove <username>`
- 设置远程登录标记（需要 root 鉴权）：
  - `go run ./tools/misaka-tools.go remote <username> true|false`
- 进入 admin-cli（需要 root 鉴权）：
  - `go run ./tools/misaka-tools.go admin-cli`

## 与服务端 login 的关系

服务端 `login` 命令只负责：

- 校验用户名密码
- 将 `ServiceConnContext.LoginUser` 写入连接上下文

实现：[ImpLogin](file:///workspace/command/Func.go#L64-L87)

因此：

- 用户/角色的创建与修改应通过 `misaka-tools`（落盘到 `profiles/user.dat`）
- 业务命令执行前，客户端通过 `login <u> <p>` 来建立会话级身份

