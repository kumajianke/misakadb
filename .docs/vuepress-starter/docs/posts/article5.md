---
date: 2026-08-09
title: Miql的数据库操作
category:
  - miql
tag:
  - miql
description: miql 对数据库的使用手册
---
## MIQL
:::tip miql是什么
MIQL是用于增删改查数据库的语法,由mq.开头。mq 必须登录用户执行。
:::

### 创建数据库
```miql
mq.createDB(<name>[, engine])
```

创建一个数据库，指定其name和引擎
```mson
{
	active: "cre-dat",
	name: <name>,
	engine: [engine(默认是tinydb)]
}
```

### 删除数据库
```miql
mq.dropDB(<name>)
```
::: warning 危险操作
操作不可逆，仅支持 `root` 权限用户删除数据库。
:::

```mson
{
	active: "drp-dat",
	name: <name>
}
```