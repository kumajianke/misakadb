

# MisakaDB

![](./logos.png)

> 便捷简单的文档数据库，始终如一。
::: warning 免责声明
**当前项目处于`α版本测试阶段`, 不建议直接用于生产环境。**
:::

::: tip misaka的诞生
misaka的命名来自动漫《魔法禁书目录》、《某科学的超电池炮》的炮姐，御坂美琴中的御坂(misaka)及衍生的御坂网络，在后期我们也会仿造misaka来实现对应的分布式技术。
:::
`MisakaDB` 是一个轻量级的 JSON 文档数据库，支持 `JSON` 格式的内容字段存储。在 0.0.3 版本，仅支持单机模式。在0.17版本之后，支持插件模式，开发者可以在Misaka源码中开发自定义的插件，从而提高CRUD的效率。

## 快速开始

### 客户端
执行 client 目录下的 `misaka-cli` 程序或者利用源码执行 `uv run client/misaka-cli.py` 启动客户端。
```
misaka-cli --mode shell
```

#### 支持参数
- `address` type=str default="127.0.0.1" help="服务地址, 默认本地"
- `port` type=int default=10032 help="服务端口, 默认10032"
- `mode` type=str default="default" help="运行模式: shell[cli调试模式]、onlyConn[仅连接]"

### miql

:::tip miql
MIQL是用于增删改查数据库的语法, 适用于客户端SDK\Shell模式对数据库、数据表的内容操作。miql的指令由mq.开头。mq 必须登录用户执行。
:::



## 许可证
本项目基于 `木兰宽松许可证` 许可证开源。详见 [LICENSE](./LICENSE) 文件。
