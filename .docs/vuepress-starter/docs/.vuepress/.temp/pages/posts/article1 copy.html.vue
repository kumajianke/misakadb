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
<span class="line">lock <span class="token punctuation">,</span> unlock<span class="token punctuation">,</span> err <span class="token operator">:=</span> global_lock<span class="token punctuation">.</span><span class="token function">GetOrStoreGlobalLock</span><span class="token punctuation">(</span><span class="token string">"lock_key"</span><span class="token punctuation">)</span></span>
<span class="line"><span class="token keyword">defer</span> <span class="token function">unlock</span><span class="token punctuation">(</span><span class="token punctuation">)</span></span>
<span class="line"><span class="token comment">// 独占锁之后需要进行的操作</span></span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><div class="hint-container warning">
<p class="hint-container-title">Warning</p>
<p>我们获取的第一个返回值虽然是lock，但是不要使用</p>
<div class="language-go line-numbers-mode" data-highlighter="prismjs" data-ext="go"><pre v-pre><code><span class="line"><span class="token keyword">defer</span> lock<span class="token punctuation">.</span><span class="token function">Unlock</span><span class="token punctuation">(</span><span class="token punctuation">)</span></span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div></div></div><p>来释放锁，这样会导致锁无法被全局锁池释放!</p>
</div>
<h2 id="知其所以然" tabindex="-1"><a class="header-anchor" href="#知其所以然"><span>知其所以然</span></a></h2>
<p>有兴趣的同学可以去看看我们的博客: <a href="http://blog.6ugo.cn/?view=tech&amp;id=316efd3b-c398-4fb2-bca0-3655893a0934" target="_blank" rel="noopener noreferrer">内部实现原理</a></p>
<h2 id="源码解析" tabindex="-1"><a class="header-anchor" href="#源码解析"><span>源码解析</span></a></h2>
<p>源码在代码中的 <code v-pre>lock/global_lock/GlobalLock.go</code> 文件中。</p>
<div class="language-go line-numbers-mode" data-highlighter="prismjs" data-ext="go"><pre v-pre><code><span class="line"><span class="token comment">// 本代码引用的是0.17的源码 如果有出入还请对照参考</span></span>
<span class="line"><span class="token operator">...</span></span>
<span class="line"><span class="token keyword">var</span> once sync<span class="token punctuation">.</span>Once</span>
<span class="line"><span class="token keyword">var</span> globalLocksPool <span class="token operator">*</span>GlobalLocksPool</span>
<span class="line"></span>
<span class="line"><span class="token comment">// 全局锁的结构</span></span>
<span class="line"><span class="token keyword">type</span> GlobalLocks <span class="token keyword">struct</span> <span class="token punctuation">{</span></span>
<span class="line">  <span class="token operator">...</span></span>
<span class="line"><span class="token punctuation">}</span></span>
<span class="line"></span>
<span class="line"><span class="token comment">// 全局锁池的结构</span></span>
<span class="line"><span class="token keyword">type</span> GlobalLocksPool <span class="token keyword">struct</span> <span class="token punctuation">{</span></span>
<span class="line">	<span class="token operator">...</span></span>
<span class="line"><span class="token punctuation">}</span></span>
<span class="line"></span>
<span class="line"><span class="token comment">// 获取一个全局单例，用于存储所有的锁</span></span>
<span class="line"><span class="token keyword">func</span> <span class="token function">GetGlobalLockPool</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token operator">*</span>GlobalLocksPool <span class="token punctuation">{</span></span>
<span class="line">  <span class="token operator">...</span></span>
<span class="line"><span class="token punctuation">}</span></span>
<span class="line"></span>
<span class="line"><span class="token comment">// 获取年轻代的锁的引用快照</span></span>
<span class="line"><span class="token keyword">func</span> <span class="token punctuation">(</span>p <span class="token operator">*</span>GlobalLocksPool<span class="token punctuation">)</span> <span class="token function">GetYoungPoolSnapshot</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token keyword">map</span><span class="token punctuation">[</span><span class="token builtin">string</span><span class="token punctuation">]</span><span class="token builtin">int32</span> <span class="token punctuation">{</span></span>
<span class="line">	<span class="token operator">...</span></span>
<span class="line"><span class="token punctuation">}</span></span>
<span class="line"></span>
<span class="line"><span class="token comment">// 获取旧代的锁的引用快照</span></span>
<span class="line"><span class="token keyword">func</span> <span class="token punctuation">(</span>p <span class="token operator">*</span>GlobalLocksPool<span class="token punctuation">)</span> <span class="token function">GetOldPoolSnapshot</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token keyword">map</span><span class="token punctuation">[</span><span class="token builtin">string</span><span class="token punctuation">]</span><span class="token builtin">int32</span> <span class="token punctuation">{</span></span>
<span class="line">	<span class="token operator">...</span></span>
<span class="line"><span class="token punctuation">}</span></span>
<span class="line"></span>
<span class="line"><span class="token comment">// 获取一个锁，如果锁是第一次创建，自动注册到全局锁池</span></span>
<span class="line"><span class="token keyword">func</span> <span class="token function">GetOrStoreGlobalLock</span><span class="token punctuation">(</span>lock_name <span class="token builtin">string</span><span class="token punctuation">,</span> lock_method <span class="token builtin">string</span><span class="token punctuation">)</span> <span class="token punctuation">(</span><span class="token operator">*</span>GlobalLocks<span class="token punctuation">,</span> <span class="token keyword">func</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">,</span> <span class="token builtin">error</span><span class="token punctuation">)</span> <span class="token punctuation">{</span></span>
<span class="line">	<span class="token operator">...</span></span>
<span class="line"><span class="token punctuation">}</span></span>
<span class="line"></span>
<span class="line"><span class="token keyword">func</span> <span class="token function">lockPoolsGCThread</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span></span>
<span class="line">	<span class="token comment">//回收机制触发</span></span>
<span class="line"><span class="token punctuation">}</span></span>
<span class="line"></span>
<span class="line"><span class="token comment">// 启动全局锁池的回收线程</span></span>
<span class="line"><span class="token keyword">func</span> <span class="token function">StartLockPoolGC</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span></span>
<span class="line">	<span class="token operator">...</span></span>
<span class="line"><span class="token punctuation">}</span></span>
<span class="line"></span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div></div></template>


