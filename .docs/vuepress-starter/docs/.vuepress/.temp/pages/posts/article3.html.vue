<template><div><div class="hint-container tip">
<p class="hint-container-title">Tips</p>
<p>misaka支持天然插件加载和开发，用户可以在 misaka的目录中的 <code v-pre>plugins/mods</code> 中开发自己的插件，官方默认加载了一个base的插件供以参考，请勿删除base插件，这会导致misaka不可用。</p>
</div>
<h2 id="插件开发" tabindex="-1"><a class="header-anchor" href="#插件开发"><span>插件开发</span></a></h2>
<p>Misaka支持插件的开发和使用。</p>
<div class="hint-container warning">
<p class="hint-container-title">Warning</p>
<p>以下设计尚在WIP中。</p>
</div>
<p><strong>模组原理</strong>：开发者创建函数，然后按照指定的type约束传入注册机，如：<code v-pre>func ([]tasktype.TaskType) []tasktype.TaskType</code>, 注册机会根据类型在程序启动之后，在相应的地方调用函数。</p>
<h3 id="使用方式" tabindex="-1"><a class="header-anchor" href="#使用方式"><span>使用方式</span></a></h3>
<p>我们在源代码中创建自己的包并编写代码，开发者可以参考Misaka官方编写的<code v-pre>base_unloaded</code>插件。开发者需要创建一个<code v-pre>plugin.yaml</code>配置文件来配置我们的插件信息，参考代码如下：</p>
<div class="language-yaml line-numbers-mode" data-highlighter="prismjs" data-ext="yml"><pre v-pre><code><span class="line"><span class="token key atrule">name</span><span class="token punctuation">:</span> <span class="token string">"基础加载模块 @misaka"</span></span>
<span class="line"><span class="token key atrule">boot</span><span class="token punctuation">:</span> <span class="token string">"./tasks.go/Register()"</span></span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div><div class="line-number"></div></div></div><p>其中属性name表示的是我们的插件名称，建议开发者在名称后面加上自己的签名这样方便区分相同的插件重名。boot是我们启动之后运行的函数。</p>
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


