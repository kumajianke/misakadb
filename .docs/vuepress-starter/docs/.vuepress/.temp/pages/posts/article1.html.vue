<template><div><h1 id="全局锁池" tabindex="-1"><a class="header-anchor" href="#全局锁池"><span>全局锁池</span></a></h1>
<div class="hint-container tip">
<p class="hint-container-title">什么是全局锁池</p>
<p>全局锁池可以帮助我们快速的在项目中对某个资源、操作句柄进行上锁。</p>
</div>
<h2 id="全局锁池的实现" tabindex="-1"><a class="header-anchor" href="#全局锁池的实现"><span>全局锁池的实现</span></a></h2>
<p>全局锁池内部被划分为两个子池：</p>
<ul>
<li><code v-pre>YoungPool</code>：存储当前高频访问、命中率较高的热点锁</li>
<li><code v-pre>OldPool</code>：存储暂时降级的锁，用于后续淘汰和回收判断</li>
</ul>
<p>业务线程在获取锁时，会优先查询 <code v-pre>YoungPool</code>。如果未命中，再继续查询 <code v-pre>OldPool</code>。一旦在 <code v-pre>OldPool</code> 中找到目标锁，就会将该锁重新提升到 <code v-pre>YoungPool</code>，从而让热点锁始终尽量停留在更快命中的路径上。</p>
<p>为了降低并发争用，双池查询使用基于 <code v-pre>xsync.Map</code> 的并发 Map 实现，避免传统全局互斥锁在高并发场景下造成过多阻塞。锁对象内部还维护了一个基于原子操作的引用计数器 <code v-pre>RefCounter</code>，用于记录锁的活跃状态。每次成功加锁时增加计数，每次释放锁时减少计数，这样后续在进行分代回收时，就可以更安全地判断哪些锁已经不再被业务使用。</p>
<p>当前实现已经具备以下能力：</p>
<ul>
<li>基于全局单例方式统一维护锁池</li>
<li>基于双池结构实现热点锁读取与旧锁提升</li>
<li>使用原子引用计数跟踪锁的使用状态</li>
<li>使用并发安全的 Map 降低锁池本身的竞争开销</li>
</ul>
<p>当前版本中，后台 GC 守护协程与完整的“年轻代批量降级、旧生代按引用计数回收”机制仍在持续完善中。也就是说，现阶段已经完成了全局锁池的核心访问链路与数据结构设计，而分代淘汰和自动回收属于下一阶段重点增强能力。</p>
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


