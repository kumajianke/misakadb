---
date: 2026-07-02
title: 原子任务中心
category:
  - 内部实现
tag:
  - 插件开发
description: 插件开发
---

:::tip
misaka支持天然插件加载和开发，用户可以在 misaka的目录中的 `plugins/mods` 中开发自己的插件，官方默认加载了一个base的插件供以参考，请勿删除base插件，这会导致misaka不可用。
:::


## 插件开发
misaka约定了tasks的插件开发.
:::warning
以下设计尚在WIP中。
:::
**模组原理**：开发者创建函数，然后按照指定的type约束传入注册机，如：`func ([]tasktype.TaskType) []tasktype.TaskType`, 注册机会根据类型在程序启动之后，在相应的地方调用函数。

当前可以走注册机有：
### 任务模块
#### AddTaskType
- 参数: `AddTaskType(allTaskTypelst []tasktype.TaskType) `
- 返回值： `[]tasktype.TaskType`
- 作用： 添加原子任务类型
> 原子任务类型: 原子任务中心会根据原子任务类型来调用对应的函数。
