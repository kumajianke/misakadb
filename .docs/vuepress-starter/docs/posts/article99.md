---
date: 2026-08-09
title: 更新日志
category:
  - 其他
tag:
  - 更新
description: 一些更新导致的改动
---

## V0.1.8
:::tip 常规更新
### fixed
- 修复Windows环境中，强制落盘目录导致无权限的Bug;
- 修复Windows环境中，部分场景下，删除数据库提醒无权限的问题;
- 修复Windows环境中，内置coin插件`OnRemoveFolder`删除已被打开的目录导致出现的问题;

### feature
- 实现 `TaskCombo`、`AddTaskCombo` 等业务函数;
- 实现删除数据库逻辑, 并挂载给 `TinyDB` 引擎;
- 修改内置库 `atoimc` 库名称为 `atomicX`;
- `engine.NewTinyEngine`

### chore
- 加入更新日志

### refactor
- 修改获取PluginX单例的函数名: `GetPluginBus` -> `GetPluginX`;
- 修改 `PluginBus` 名称为 `PluginsX`;
- 修改 内置插件 `base_unloader` 插件名称为 `coin`;
- 将 `AtomicDropDB` 的逻辑转交给 内置插件 `coin` 执行；


:::

## V0.1.9
:::tip 常规更新
### fixed

### feature
- `judgment`清道夫功能，物理删除任务完成状态的所有任务。
- `pluginx` 的 `AfterTask` 的支持，可以在任务完成之后执行指定内容（必须执行的任务）
  - [tip] 后期开发过程中删除转为`AddAfter`的支持
- `OnRemoveFolder` 

### chore

### refactor
- 重构了关于 `TaskBooksShipBuilder` 相关的类变量的名称(旧名称为: TaskTypeBuilder、TaskTypeShipBuilder)

:::

:::

## TODO LIST
- 【coin】修改文件的`TaskType`及其相关函数
- 【coin】创建文件的`TaskType`及其相关函数

:::