<template><div><div class="hint-container tip">
<p class="hint-container-title">Tips</p>
<p>misaka支持天然插件加载和开发，用户可以在 misaka 开发自己的插件，官方默认加载了一个misaka basic mode的插件供以参考，请勿删除misaka basic mode插件，这会导致misaka不可用。</p>
</div>
<h2 id="插件开发" tabindex="-1"><a class="header-anchor" href="#插件开发"><span>插件开发</span></a></h2>
<p>Misaka支持插件的开发和使用。</p>
<div class="hint-container warning">
<p class="hint-container-title">Warning</p>
<p>以下设计尚在WIP中。</p>
</div>
<p><strong>模组原理</strong>：开发者创建函数，然后按照指定的type约束传入注册机，如：<code v-pre>func ([]tasktype.TaskType) []tasktype.TaskType</code>, 注册机会根据类型在程序启动之后，在相应的地方调用函数。</p>
<h3 id="快捷启动" tabindex="-1"><a class="header-anchor" href="#快捷启动"><span>快捷启动</span></a></h3>
<p>我们在源代码中创建自己的包并编写代码，开发者可以参考Misaka官方编写的<code v-pre>base_unloaded</code>插件。开发者需要创建一个<code v-pre>plugin.yaml</code>配置文件来配置我们的插件信息，参考代码如下：</p>
<div class="language-yaml line-numbers-mode" data-highlighter="prismjs" data-ext="yml"><pre v-pre><code><span class="line"><span class="token key atrule">name</span><span class="token punctuation">:</span> <span class="token string">"基础加载模块 @misaka"</span></span>
<span class="line"><span class="token key atrule">boot</span><span class="token punctuation">:</span> <span class="token string">"./tasks.go/Register()"</span></span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div><div class="line-number"></div></div></div><p>其中属性name表示的是我们的插件名称，建议开发者在名称后面加上自己的签名这样方便区分相同的插件重名。boot是我们启动之后运行的函数。</p>
<p>接着编写插件代码，贴上参考代码便于启动:</p>
<div class="language-go line-numbers-mode" data-highlighter="prismjs" data-ext="go"><pre v-pre><code><span class="line"><span class="token keyword">package</span> base_unloader</span>
<span class="line"></span>
<span class="line"><span class="token keyword">import</span> <span class="token punctuation">(</span></span>
<span class="line">	<span class="token string">"fmt"</span></span>
<span class="line">	tasktype <span class="token string">"misakadb/atomic/atomicWorkCenter/TaskType"</span></span>
<span class="line">	<span class="token string">"misakadb/clilog"</span></span>
<span class="line">	pluginsloader <span class="token string">"misakadb/plugins/pluginsLoader"</span></span>
<span class="line"><span class="token punctuation">)</span></span>
<span class="line"></span>
<span class="line"><span class="token keyword">const</span> <span class="token punctuation">(</span></span>
<span class="line">	TaskRemoveFolder tasktype<span class="token punctuation">.</span>TaskType <span class="token operator">=</span> <span class="token string">"remove_folder"</span></span>
<span class="line"><span class="token punctuation">)</span></span>
<span class="line"></span>
<span class="line"><span class="token keyword">func</span> <span class="token function">AddTaskType</span><span class="token punctuation">(</span>allTaskTypelst <span class="token punctuation">[</span><span class="token punctuation">]</span>tasktype<span class="token punctuation">.</span>TaskType<span class="token punctuation">)</span> <span class="token punctuation">(</span><span class="token punctuation">[</span><span class="token punctuation">]</span>tasktype<span class="token punctuation">.</span>TaskType<span class="token punctuation">,</span> <span class="token builtin">error</span><span class="token punctuation">)</span> <span class="token punctuation">{</span></span>
<span class="line">	allTaskTypelst <span class="token operator">=</span> <span class="token function">append</span><span class="token punctuation">(</span>allTaskTypelst<span class="token punctuation">,</span> TaskRemoveFolder<span class="token punctuation">)</span></span>
<span class="line">	<span class="token keyword">return</span> allTaskTypelst<span class="token punctuation">,</span> <span class="token boolean">nil</span></span>
<span class="line"><span class="token punctuation">}</span></span>
<span class="line"></span>
<span class="line"><span class="token keyword">func</span> <span class="token function">OnRemoveFile</span><span class="token punctuation">(</span>taskType tasktype<span class="token punctuation">.</span>TaskType<span class="token punctuation">,</span> params <span class="token punctuation">[</span><span class="token punctuation">]</span><span class="token builtin">string</span><span class="token punctuation">)</span> <span class="token builtin">error</span> <span class="token punctuation">{</span></span>
<span class="line">	<span class="token comment">// 实现对文件的删除操作</span></span>
<span class="line">	<span class="token keyword">return</span> <span class="token boolean">nil</span></span>
<span class="line"></span>
<span class="line"><span class="token punctuation">}</span></span>
<span class="line"></span>
<span class="line"><span class="token keyword">func</span> <span class="token function">RollRemoveFile</span><span class="token punctuation">(</span>taskType tasktype<span class="token punctuation">.</span>TaskType<span class="token punctuation">,</span> params <span class="token punctuation">[</span><span class="token punctuation">]</span><span class="token builtin">string</span><span class="token punctuation">)</span> <span class="token builtin">error</span> <span class="token punctuation">{</span></span>
<span class="line">	<span class="token comment">// 回滚删除操作</span></span>
<span class="line">	<span class="token keyword">return</span> <span class="token boolean">nil</span></span>
<span class="line"><span class="token punctuation">}</span></span>
<span class="line"></span>
<span class="line"><span class="token keyword">func</span> <span class="token function">init</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span></span>
<span class="line">	pluginsloader<span class="token punctuation">.</span><span class="token function">RegisterBuiltinPlugin</span><span class="token punctuation">(</span>modName<span class="token punctuation">,</span> Register<span class="token punctuation">)</span></span>
<span class="line"><span class="token punctuation">}</span></span>
<span class="line"></span>
<span class="line"><span class="token keyword">func</span> <span class="token function">Register</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token builtin">error</span> <span class="token punctuation">{</span></span>
<span class="line">	<span class="token keyword">if</span> err <span class="token operator">:=</span> pluginsloader<span class="token punctuation">.</span><span class="token function">RegisterPluginsActionsInTaskType</span><span class="token punctuation">(</span>modName<span class="token punctuation">,</span> AddTaskType<span class="token punctuation">)</span><span class="token punctuation">;</span> err <span class="token operator">!=</span> <span class="token boolean">nil</span> <span class="token punctuation">{</span></span>
<span class="line">		<span class="token keyword">return</span> fmt<span class="token punctuation">.</span><span class="token function">Errorf</span><span class="token punctuation">(</span><span class="token string">"register task types failed: %w"</span><span class="token punctuation">,</span> err<span class="token punctuation">)</span></span>
<span class="line">	<span class="token punctuation">}</span></span>
<span class="line"></span>
<span class="line">	<span class="token keyword">if</span> err <span class="token operator">:=</span> pluginsloader<span class="token punctuation">.</span><span class="token function">RegisterPluginsActionsInTaskTypeAction</span><span class="token punctuation">(</span>modName<span class="token punctuation">,</span> TaskRemoveFile<span class="token punctuation">,</span> OnRemoveFile<span class="token punctuation">)</span><span class="token punctuation">;</span> err <span class="token operator">!=</span> <span class="token boolean">nil</span> <span class="token punctuation">{</span></span>
<span class="line">		<span class="token keyword">return</span> fmt<span class="token punctuation">.</span><span class="token function">Errorf</span><span class="token punctuation">(</span><span class="token string">"register task action failed: %w"</span><span class="token punctuation">,</span> err<span class="token punctuation">)</span></span>
<span class="line">	<span class="token punctuation">}</span></span>
<span class="line">	clilog<span class="token punctuation">.</span><span class="token function">Success</span><span class="token punctuation">(</span><span class="token string">"基础插件加载完毕."</span><span class="token punctuation">)</span></span>
<span class="line">	<span class="token keyword">return</span> <span class="token boolean">nil</span></span>
<span class="line"><span class="token punctuation">}</span></span>
<span class="line"></span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h3 id="插件的安装" tabindex="-1"><a class="header-anchor" href="#插件的安装"><span>插件的安装</span></a></h3>
<p>安装插件的方式是使用misaka-tools, misaka-tools提供几个插件管理的函数:</p>
<table>
<thead>
<tr>
<th>参数</th>
<th>作用</th>
</tr>
</thead>
<tbody>
<tr>
<td>misaka-tools plu-add [插件的目录]</td>
<td>添加插件到misaka</td>
</tr>
<tr>
<td>misaka-tools plu-remove [插件的名称]</td>
<td>删除指定的插件</td>
</tr>
<tr>
<td>misaka-tools plu-list</td>
<td>获取所有的插件信息</td>
</tr>
</tbody>
</table>
<div class="hint-container warning">
<p class="hint-container-title">提示</p>
<ul>
<li><code v-pre>misaka-tools</code>卸载或者添加一个插件, 都需要重新编译一次misaka，<strong>这需要用户环境配置misaka所需的Go语言环境</strong>。</li>
<li>默认分发的Misaka都会自动安装一个叫做<code v-pre>misaka basic mode</code>的插件，这个插件提供了很多基础的功能，所以不建议删除或者覆盖。</li>
</ul>
</div>
<p>执行完毕之后插件的信息可以通过misaka-tools查看，也可以执行编译好的misaka进行查看，</p>
<h3 id="任务模块" tabindex="-1"><a class="header-anchor" href="#任务模块"><span>任务模块</span></a></h3>
<h4 id="addtasktype" tabindex="-1"><a class="header-anchor" href="#addtasktype"><span>AddTaskType</span></a></h4>
<ul>
<li>参数: <code v-pre>AddTaskType(allTaskTypelst []tasktype.TaskType) </code></li>
<li>返回值： <code v-pre>[]tasktype.TaskType</code></li>
<li>作用： 添加原子任务类型</li>
</ul>
<blockquote>
<p>原子任务类型: 原子任务中心会根据原子任务类型来调用对应的函数。</p>
</blockquote>
<h3 id="" tabindex="-1"><a class="header-anchor" href="#"><span></span></a></h3>
</div></template>


