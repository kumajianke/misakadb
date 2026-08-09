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
misaka支持天然插件加载和开发，用户可以在 misaka 开发自己的插件，官方默认加载了一个misaka basic mode的插件供以参考，请勿删除misaka basic mode插件，这会导致misaka不可用。
:::


## 插件开发
Misaka支持插件的开发和使用。

**模组原理**：开发者创建函数，然后按照指定的type约束传入注册机，如：`func ([]tasktype.TaskType) []tasktype.TaskType`, 注册机会根据类型在程序启动之后，在相应的地方调用函数。

### 快捷启动
我们在源代码中创建自己的包并编写代码，开发者可以参考Misaka官方编写的`base_unloaded`插件。开发者需要创建一个`plugin.yaml`配置文件来配置我们的插件信息，参考代码如下：
```yaml
name: "基础加载模块 @misaka"
boot: "./tasks.go/Register()"
```
其中属性name表示的是我们的插件名称，建议开发者在名称后面加上自己的签名这样方便区分相同的插件重名。boot是我们启动之后运行的函数。

接着编写插件代码，贴上参考代码便于启动:
```go
package base_unloader

import (
	"fmt"
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	"misakadb/clilog"
	pluginsloader "misakadb/plugins/pluginsLoader"
)

const (
	TaskRemoveFile   tasktype.TaskType = "remove_file"
	TaskModFile      tasktype.TaskType = "mod_file"
	TaskRemoveFolder tasktype.TaskType = "remove_folder"
)

func OnRemoveFile(taskType tasktype.TaskType, params []string) error {
	// 实现对文件的删除操作
	return nil
}

func RollRemoveFile(taskType tasktype.TaskType, params []string) error {
	// 回滚删除操作
	return nil
}

/*
添加TaskType
*/
func AddTaskType() error {
	if err := pluginsloader.RegisterPluginTaskTypeWithAlias(modName, "misaka.removefile@用于删除文件的tasktype", TaskRemoveFile); err != nil {
		return fmt.Errorf("register alias %s failed: %w", "misaka.removefile", err)
	}
	if err := pluginsloader.RegisterPluginTaskTypeWithAlias(modName, "misaka.removefolder@用于删目录的tasktype", TaskRemoveFolder); err != nil {
		return fmt.Errorf("register alias %s failed: %w", "misaka.removefolder", err)
	}
	return nil
}

/*
添加TaskTypeAction
*/
func AddTaskTypeAction() error {
	if err := pluginsloader.RegisterPluginsActionsInTaskTypeAction(modName, TaskRemoveFile, OnRemoveFile); err != nil {
		return fmt.Errorf("register task action failed: %w", err)
	}
	return nil
}

func Register() error {
	if err := AddTaskType(); err != nil {
		return err
	}

	if err := AddTaskTypeAction(); err != nil {
		return err
	}

	clilog.Success("基础插件加载完毕.")
	return nil
}


```

### 插件的安装
安装插件的方式是使用misaka-tools, misaka-tools提供几个插件管理的函数:

|参数|作用|
|--|--|
|misaka-tools plu-add [插件的目录]|添加插件到misaka|
|misaka-tools plu-remove [插件的名称]|删除指定的插件|
|misaka-tools plu-list|获取所有的插件信息|

:::warning 提示
- `misaka-tools`卸载或者添加一个插件, 都需要重新编译一次misaka，**这需要用户环境配置misaka所需的Go语言环境**。
- 默认分发的Misaka都会自动安装一个叫做`misaka basic mode`的插件，这个插件提供了很多基础的功能，所以不建议删除或者覆盖。
:::

执行完毕之后插件的信息可以通过misaka-tools查看，也可以执行编译好的misaka进行查看，在`Debug`模式下打开Misaka，按住键盘`p`键查看添加的和加载的插件信息。

### 常见插件函数
> `TaskType`和`TaskTypeAction`是原子任务中心的一个概念，[点击跳转](article2.html)。
#### 插件追加 `TaskType`
- **函数**：`RegisterPluginTaskTypeWithAlias()`
- **参数**：`(plugin string, alias string, taskType tasktype.TaskType)`
- **作用**：传入插件名称、TaskType别名、taskType对象以注册TaskType到pluginsx。

#### 插件追加 `TaskTypeAction`
- **函数**：`RegisterPluginsActionsInTaskTypeAction()`
- **参数**：`(plugins string, taskType tasktype.TaskType, action pluginsxInterface.FuncTaskTypeAction)`
> ```go
> type FuncTaskTypeAction = func(taskType tasktype.TaskType, params []string) error
> ```
- **作用**：传入插件名称、TaskType别名、taskType对象以注册TaskType到pluginsx。

### 常见插件术语
#### pluginsx
`pluginsx` 是 `Misaka` 用于管理插件上下的包，插件注册的所有内容都放在了 `pluginsx.PluginsX{}` 对象上。其中PluginsX是单例存在的，使用函数 `GetPluginsX` 获取唯一引用。

