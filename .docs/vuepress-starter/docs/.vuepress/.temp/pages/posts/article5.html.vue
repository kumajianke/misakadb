<template><div><h2 id="miql" tabindex="-1"><a class="header-anchor" href="#miql"><span>MIQL</span></a></h2>
<div class="hint-container tip">
<p class="hint-container-title">miql是什么</p>
<p>MIQL是用于增删改查数据库的语法,由mq.开头。mq 必须登录用户执行。</p>
</div>
<h3 id="创建数据库" tabindex="-1"><a class="header-anchor" href="#创建数据库"><span>创建数据库</span></a></h3>
<div class="language-miql line-numbers-mode" data-highlighter="prismjs" data-ext="miql"><pre v-pre><code><span class="line">mq.createDB(&lt;name&gt;[, engine])</span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div></div></div><p>创建一个数据库，指定其name和引擎</p>
<div class="language-mson line-numbers-mode" data-highlighter="prismjs" data-ext="mson"><pre v-pre><code><span class="line">{</span>
<span class="line">	active: &quot;cre-dat&quot;,</span>
<span class="line">	name: &lt;name&gt;,</span>
<span class="line">	engine: [engine(默认是tinydb)]</span>
<span class="line">}</span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h3 id="删除数据库" tabindex="-1"><a class="header-anchor" href="#删除数据库"><span>删除数据库</span></a></h3>
<div class="language-miql line-numbers-mode" data-highlighter="prismjs" data-ext="miql"><pre v-pre><code><span class="line">mq.dropDB(&lt;name&gt;)</span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div></div></div><div class="hint-container warning">
<p class="hint-container-title">危险操作</p>
<p>操作不可逆，仅支持 <code v-pre>root</code> 权限用户删除数据库。</p>
</div>
<div class="language-mson line-numbers-mode" data-highlighter="prismjs" data-ext="mson"><pre v-pre><code><span class="line">{</span>
<span class="line">	active: &quot;drp-dat&quot;,</span>
<span class="line">	name: &lt;name&gt;</span>
<span class="line">}</span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div></div></template>


