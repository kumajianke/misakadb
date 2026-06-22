<template><div><h1 id="全局锁池" tabindex="-1"><a class="header-anchor" href="#全局锁池"><span>全局锁池</span></a></h1>
<h2 id="什么是全局锁池" tabindex="-1"><a class="header-anchor" href="#什么是全局锁池"><span>什么是全局锁池</span></a></h2>
<div class="hint-container tip">
<p class="hint-container-title">Tips</p>
<p>全局锁池可以帮助我们快速的在项目中对某个资源、操作句柄进行上锁。</p>
</div>
<h2 id="使用方法" tabindex="-1"><a class="header-anchor" href="#使用方法"><span>使用方法</span></a></h2>
<div class="language-go line-numbers-mode" data-highlighter="prismjs" data-ext="go"><pre v-pre><code><span class="line"><span class="token keyword">import</span> <span class="token punctuation">(</span></span>
<span class="line">	<span class="token string">"misakadb/lock/global_lock"</span></span>
<span class="line"><span class="token punctuation">)</span></span>
<span class="line"></span>
<span class="line">lock <span class="token punctuation">,</span> unlock<span class="token punctuation">,</span> err <span class="token operator">:=</span>global_lock<span class="token punctuation">.</span><span class="token function">GetOrStoreGlobalLock</span><span class="token punctuation">(</span><span class="token string">"lock_key"</span><span class="token punctuation">)</span></span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div></div></template>


