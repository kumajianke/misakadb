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
- Windows环境中，强制落盘目录导致无权限的Bug;
- Windows环境中，部分场景下，删除数据库无权限的提醒;

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

### 下个版本代办
功能	状态
正常删除目录	部分完成，实际是改名
正常失败回滚	基本有
删除后的物理清理	未完成
alldb.list 同步删除	未完成
程序中断后加载任务	已完成
程序启动后自动继续任务	未完成
删除动作中断后的准确恢复	未完成
回滚失败后的任务保留	未完成
成功任务清理	未完成
Windows 目录同步	已处理
:::