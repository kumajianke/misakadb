<template><div><h2 id="v0-1-8" tabindex="-1"><a class="header-anchor" href="#v0-1-8"><span>V0.1.8</span></a></h2>
<div class="hint-container tip">
<p class="hint-container-title">常规更新</p>
<h3 id="fixed" tabindex="-1"><a class="header-anchor" href="#fixed"><span>fixed</span></a></h3>
<ul>
<li>Windows环境中，强制落盘目录导致无权限的Bug;</li>
<li>Windows环境中，部分场景下，删除数据库无权限的提醒。</li>
</ul>
<h3 id="feature" tabindex="-1"><a class="header-anchor" href="#feature"><span>feature</span></a></h3>
<ul>
<li>实现 <code v-pre>TaskCombo</code>、<code v-pre>AddTaskCombo</code> 等业务函数;</li>
<li>实现删除数据库逻辑;</li>
</ul>
<h3 id="chore" tabindex="-1"><a class="header-anchor" href="#chore"><span>chore</span></a></h3>
<ul>
<li>加入更新日志</li>
</ul>
<h3 id="refactor" tabindex="-1"><a class="header-anchor" href="#refactor"><span>refactor</span></a></h3>
<ul>
<li>修改获取PluginX单例的函数名: <code v-pre>GetPluginBus</code> -&gt; <code v-pre>GetPluginX</code>;</li>
<li>修改 <code v-pre>PluginBus</code> 名称为 <code v-pre>PluginsX</code>;</li>
<li>修改 内置插件 <code v-pre>base_unloader</code> 插件名称为 <code v-pre>coin</code>;</li>
<li>将 <code v-pre>AtomicDropDB</code> 的逻辑转交给 内置插件 <code v-pre>coin</code> 执行；</li>
</ul>
<h3 id="下个版本代办" tabindex="-1"><a class="header-anchor" href="#下个版本代办"><span>下个版本代办</span></a></h3>
<p>功能	状态
正常删除目录	部分完成，实际是改名
正常失败回滚	基本有
删除后的物理清理	未完成
alldb.list 同步删除	未完成
程序中断后加载任务	已完成
程序启动后自动继续任务	未完成
删除动作中断后的准确恢复	未完成
回滚失败后的任务保留	未完成
成功任务清理	未完成
Windows 目录同步	已处理</p>
</div>
</div></template>


