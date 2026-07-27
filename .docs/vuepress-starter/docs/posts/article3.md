---
date: 2026-07-02
title: 插件开发
category:
  - 插件
tag:
  - 插件开发
description: 插件开发规范
---

:::tip
misaka支持天然插件加载和开发，用户可以在 misaka的目录中的 `plugins/mods` 中开发自己的插件，官方默认加载了一个base的插件供以参考，请勿删除base插件，这会导致misaka不可用。
:::


## 插件开发
Misaka支持插件的开发和使用。
:::warning
以下设计尚在WIP中。
:::
**模组原理**：开发者创建函数，然后按照指定的type约束传入注册机，如：`func ([]tasktype.TaskType) []tasktype.TaskType`, 注册机会根据类型在程序启动之后，在相应的地方调用函数。

### 使用方式
我们在源代码中创建自己的包并编写代码，开发者可以参考Misaka官方编写的`base_unloaded`插件。开发者需要创建一个`plugin.yaml`配置文件来配置我们的插件信息，参考代码如下：
```yaml
name: "基础加载模块 @misaka"
boot: "./tasks.go/Register()"
```
其中属性name表示的是我们的插件名称，建议开发者在名称后面加上自己的签名这样方便区分相同的插件重名。boot是我们启动之后运行的函数。
### 任务模块
#### AddTaskType
- 参数: `AddTaskType(allTaskTypelst []tasktype.TaskType) `
- 返回值： `[]tasktype.TaskType`
- 作用： 添加原子任务类型
> 原子任务类型: 原子任务中心会根据原子任务类型来调用对应的函数。

### 
