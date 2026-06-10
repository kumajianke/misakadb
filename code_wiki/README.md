# MisakaDB Code Wiki

这份 Code Wiki 面向需要继续开发 MisakaDB（当前 WIP）的同学，目标是把“从哪里读 → 怎么跑起来 → 主链路怎么走 → 模块怎么协作 → 目前缺什么”讲清楚。

## 目录

- [01-Overall.md](./01-Overall.md)：项目定位、核心概念、目录结构速览
- [02-Architecture.md](./02-Architecture.md)：整体架构与关键数据流（含时序/依赖）
- [03-Modules.md](./03-Modules.md)：主要模块职责与边界
- [04-Network-Protocol.md](./04-Network-Protocol.md)：TCP 协议、心跳与收发帧格式
- [05-Commands.md](./05-Commands.md)：命令体系（普通命令/MIQL）与分发机制
- [06-Engine-TinyDB.md](./06-Engine-TinyDB.md)：引擎接口、TinyDB 的落盘结构与关键实现
- [07-Auth-And-Tools.md](./07-Auth-And-Tools.md)：用户体系、安全组件与 misaka-tools
- [08-WIP-Checklist.md](./08-WIP-Checklist.md)：WIP 现状清单与下一步开发任务（含 DROPDB）

## 推荐阅读顺序

1. 01-Overall.md（先建立全局概念）
2. 02-Architecture.md（抓主链路）
3. 05-Commands.md + 04-Network-Protocol.md（理解协议与命令分发）
4. 06-Engine-TinyDB.md（进入存储/引擎细节）
5. 08-WIP-Checklist.md（开始接手开发）

