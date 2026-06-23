<template><div><h1 id="原子任务中心的作用" tabindex="-1"><a class="header-anchor" href="#原子任务中心的作用"><span>原子任务中心的作用</span></a></h1>
<p>当客户端发来<span style="color:#0099FF"><code v-pre>miql</code></span>之后，misakadb要统一进行解析，为了保持miql的任务的原子性、一致性、隔离性、持久性这四个属性，misakaDB需要一个原子任务中心来管理所有的任务。</p>
<blockquote>
<p>如当用户发送一个删除任务，这个任务拆分下来就是需要对文件、JSON等进行删除的多个子任务。这些任务需要保证要么不执行要么全部执行。</p>
</blockquote>
<p>所以我们执行多个需要多步骤的任务不能直接进行执行，而是把子任务一次性追加到我们的新任务中。</p>
<div class="language-go line-numbers-mode" data-highlighter="prismjs" data-ext="go"><pre v-pre><code><span class="line">atomic <span class="token operator">:=</span> atomic_work_center<span class="token punctuation">.</span><span class="token function">NewAtomicWorkCenter</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">;</span></span>
<span class="line"><span class="token comment">// 创建一个原子任务中心</span></span>
<span class="line">task <span class="token operator">=</span> atomic_work_center<span class="token punctuation">.</span><span class="token function">NewTask</span><span class="token punctuation">(</span><span class="token boolean">nil</span><span class="token punctuation">)</span><span class="token punctuation">;</span></span>
<span class="line"><span class="token comment">// 创建一个空任务</span></span>
<span class="line">task<span class="token punctuation">.</span>TaskBody <span class="token operator">=</span> tasktype<span class="token punctuation">.</span><span class="token function">NewShipBuilder</span><span class="token punctuation">(</span><span class="token punctuation">)</span></span>
<span class="line">	<span class="token punctuation">.</span><span class="token function">Add</span><span class="token punctuation">(</span>tasktype<span class="token punctuation">.</span>TaskModFile<span class="token punctuation">,</span> <span class="token string">"xxx.filename.txt"</span><span class="token punctuation">)</span></span>
<span class="line">	<span class="token punctuation">.</span><span class="token function">Build</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">;</span></span>
<span class="line"><span class="token comment">// 构建一个任务体，任务体包含需要原子执行的多个任务子作业</span></span>
<span class="line"></span>
<span class="line">ok<span class="token punctuation">,</span> task_id <span class="token operator">:=</span> atomic<span class="token punctuation">.</span><span class="token function">addTask</span><span class="token punctuation">(</span>task<span class="token punctuation">,</span> <span class="token number">3</span><span class="token punctuation">)</span> <span class="token comment">// 尝试添加任务到原子中心 可以重试3次</span></span>
<span class="line"><span class="token comment">// ok: 表示是否添加成功</span></span>
<span class="line"><span class="token comment">// task_id: 表示任务在中心的唯一标志</span></span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h1 id="原子日志序列器" tabindex="-1"><a class="header-anchor" href="#原子日志序列器"><span>原子日志序列器</span></a></h1>
<p><span style="color:#0099FF"><strong>AtomicWorkCenter</strong></span> 管理着所有的任务，我们将这些任务数据可以称为WAL日志。防止宕机等问题导致的数据丢失，我们会在启动之后启动对应的序列器协程。如果需要让协程启动序列化任务可以使用AtomicWorkEventBus总线进行信号通知。</p>
<div class="language-go line-numbers-mode" data-highlighter="prismjs" data-ext="go"><pre v-pre><code><span class="line"><span class="token keyword">var</span> eb <span class="token operator">*</span>eventbus<span class="token punctuation">.</span>AtomicWorkEventBus</span>
<span class="line">eb <span class="token operator">=</span> eventbus<span class="token punctuation">.</span><span class="token function">NewAtomicWorkCenterEventBus</span><span class="token punctuation">(</span><span class="token punctuation">)</span></span>
<span class="line">eb<span class="token punctuation">.</span>EventBus <span class="token operator">&lt;-</span> <span class="token string">"sync-to-local"</span> <span class="token comment">// 通知序列器协程启动序列化任务</span></span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>我们也可以加载本地日志到内存:</p>
</div></template>


