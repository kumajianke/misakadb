<template><div><h1 id="misakadb" tabindex="-1"><a class="header-anchor" href="#misakadb"><span>MisakaDB</span></a></h1>
<p><img src="@source/logos.png" alt=""></p>
<blockquote>
<p>[!NOTE]
便捷简单的文档数据库，始终如一。</p>
</blockquote>
<p><code v-pre>MisakaDB</code> 是一个轻量级的 JSON 文档数据库，支持 <code v-pre>JSON</code> 格式的内容字段存储。在 0.0.3 版本，仅支持单机模式。</p>
<h2 id="特性" tabindex="-1"><a class="header-anchor" href="#特性"><span>特性</span></a></h2>
<ul>
<li>轻量级设计，部署简单</li>
<li>支持 JSON 格式数据存储</li>
<li>Go 语言编写的高性能服务端</li>
<li>Python 客户端 SDK</li>
<li>心跳保活机制</li>
<li>连接管理与命令分发</li>
</ul>
<h2 id="全局锁池" tabindex="-1"><a class="header-anchor" href="#全局锁池"><span>全局锁池</span></a></h2>
<p>MisakaDB 在锁管理上引入了全局锁池设计，用统一的池化方式管理业务锁对象，减少重复创建与散落管理带来的开销。整个方案围绕“双池分代 + 原子引用计数 + 热点锁升级”展开，目标是在保证并发性能的同时，为后续锁回收和长期运行时的内存稳定性打下基础。</p>
<h3 id="原理" tabindex="-1"><a class="header-anchor" href="#原理"><span>原理</span></a></h3>
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
<h3 id="设计价值" tabindex="-1"><a class="header-anchor" href="#设计价值"><span>设计价值</span></a></h3>
<ul>
<li>减少业务层重复创建锁对象的成本</li>
<li>让热点锁优先停留在高命中路径，提升锁获取效率</li>
<li>通过原子计数为后续无阻塞回收提供判断依据</li>
<li>为长期运行场景中的锁对象治理和内存回收预留扩展空间</li>
</ul>
<h2 id="项目结构" tabindex="-1"><a class="header-anchor" href="#项目结构"><span>项目结构</span></a></h2>
<div class="language-text line-numbers-mode" data-highlighter="prismjs" data-ext="text"><pre v-pre><code><span class="line">misakadb/</span>
<span class="line">├── bin/                    # 编译后的二进制文件</span>
<span class="line">├── client/                 # Python 客户端</span>
<span class="line">│   ├── apis/              # API 实现</span>
<span class="line">│   ├── interface/        # 接口定义</span>
<span class="line">│   ├── network/           # 网络通信</span>
<span class="line">│   └── usage/             # 使用示例</span>
<span class="line">├── clilog/                 # 日志模块</span>
<span class="line">├── command/                # 命令处理</span>
<span class="line">├── config/                 # 配置管理</span>
<span class="line">├── engine/                 # 数据库引擎</span>
<span class="line">├── misaka-doc/             # 文档</span>
<span class="line">├── network/                # 网络层</span>
<span class="line">├── profiles/               # 配置文件</span>
<span class="line">├── safe/                   # 安全模块</span>
<span class="line">├── tools/                  # 工具集</span>
<span class="line">└── miusers/                # 用户管理</span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h2 id="快速开始" tabindex="-1"><a class="header-anchor" href="#快速开始"><span>快速开始</span></a></h2>
<h3 id="启动服务端" tabindex="-1"><a class="header-anchor" href="#启动服务端"><span>启动服务端</span></a></h3>
<h4 id="源码启动" tabindex="-1"><a class="header-anchor" href="#源码启动"><span>源码启动</span></a></h4>
<p>需要安装 Go 语言环境后执行：</p>
<div class="language-bash line-numbers-mode" data-highlighter="prismjs" data-ext="sh"><pre v-pre><code><span class="line">go run misaka.go <span class="token punctuation">[</span>参数<span class="token punctuation">]</span></span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div></div></div><h4 id="程序启动" tabindex="-1"><a class="header-anchor" href="#程序启动"><span>程序启动</span></a></h4>
<p>首先编译：</p>
<div class="language-bash line-numbers-mode" data-highlighter="prismjs" data-ext="sh"><pre v-pre><code><span class="line">go build misaka</span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div></div></div><p>然后启动编译后的程序即可。</p>
<h4 id="配置文件启动" tabindex="-1"><a class="header-anchor" href="#配置文件启动"><span>配置文件启动</span></a></h4>
<p>配置文件 <code v-pre>misaka.yaml</code> 内容示例：</p>
<div class="language-yaml line-numbers-mode" data-highlighter="prismjs" data-ext="yml"><pre v-pre><code><span class="line"><span class="token key atrule">network</span><span class="token punctuation">:</span></span>
<span class="line">  <span class="token key atrule">port</span><span class="token punctuation">:</span> <span class="token number">8080</span></span>
<span class="line">  <span class="token key atrule">address</span><span class="token punctuation">:</span> 0.0.0.0</span>
<span class="line">  <span class="token key atrule">max_conn</span><span class="token punctuation">:</span> <span class="token number">100000</span></span>
<span class="line">  <span class="token key atrule">retry_count</span><span class="token punctuation">:</span> <span class="token number">3</span></span>
<span class="line">  <span class="token key atrule">retry_delay</span><span class="token punctuation">:</span> <span class="token number">30</span></span>
<span class="line"></span>
<span class="line"><span class="token key atrule">service</span><span class="token punctuation">:</span></span>
<span class="line">  <span class="token key atrule">version</span><span class="token punctuation">:</span> <span class="token string">"0.0.3"</span></span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p><strong>配置说明：</strong></p>
<table>
<thead>
<tr>
<th>配置项</th>
<th>说明</th>
</tr>
</thead>
<tbody>
<tr>
<td>network.port</td>
<td>服务端口</td>
</tr>
<tr>
<td>network.address</td>
<td>监听地址</td>
</tr>
<tr>
<td>network.max_conn</td>
<td>最大连接数</td>
</tr>
<tr>
<td>network.retry_count</td>
<td>重试次数</td>
</tr>
<tr>
<td>network.retry_delay</td>
<td>心跳/重试超时时间(秒)</td>
</tr>
<tr>
<td>service.version</td>
<td>版本号</td>
</tr>
</tbody>
</table>
<h4 id="命令行参数" tabindex="-1"><a class="header-anchor" href="#命令行参数"><span>命令行参数</span></a></h4>
<div class="language-text line-numbers-mode" data-highlighter="prismjs" data-ext="text"><pre v-pre><code><span class="line">-address    启动的地址簇</span>
<span class="line">-port       启动的端口</span>
<span class="line">-configs    配置文件路径【默认 misaka.yaml】</span>
<span class="line">-debug      调试模式，上线前确保为 false，会在 6060 端口启动 pprof 服务</span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h3 id="使用客户端" tabindex="-1"><a class="header-anchor" href="#使用客户端"><span>使用客户端</span></a></h3>
<h4 id="源码运行" tabindex="-1"><a class="header-anchor" href="#源码运行"><span>源码运行</span></a></h4>
<div class="language-bash line-numbers-mode" data-highlighter="prismjs" data-ext="sh"><pre v-pre><code><span class="line">uv run main.py</span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div></div></div><h4 id="使用打包好的客户端" tabindex="-1"><a class="header-anchor" href="#使用打包好的客户端"><span>使用打包好的客户端</span></a></h4>
<div class="language-bash line-numbers-mode" data-highlighter="prismjs" data-ext="sh"><pre v-pre><code><span class="line">main.exe</span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div></div></div><h4 id="基础用法示例" tabindex="-1"><a class="header-anchor" href="#基础用法示例"><span>基础用法示例</span></a></h4>
<div class="language-python line-numbers-mode" data-highlighter="prismjs" data-ext="py"><pre v-pre><code><span class="line"><span class="token keyword">from</span> apis <span class="token keyword">import</span> connect</span>
<span class="line"></span>
<span class="line"><span class="token comment"># 创建并连接客户端</span></span>
<span class="line">client <span class="token operator">=</span> connect<span class="token punctuation">(</span><span class="token string">"127.0.0.1"</span><span class="token punctuation">,</span> <span class="token number">10032</span><span class="token punctuation">)</span></span>
<span class="line"></span>
<span class="line"><span class="token comment"># 执行命令</span></span>
<span class="line">result <span class="token operator">=</span> client<span class="token punctuation">.</span>execute_command<span class="token punctuation">(</span><span class="token string">"your_command"</span><span class="token punctuation">)</span></span>
<span class="line"></span>
<span class="line"><span class="token comment"># 关闭连接</span></span>
<span class="line">client<span class="token punctuation">.</span>close<span class="token punctuation">(</span><span class="token punctuation">)</span></span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h2 id="api-参考" tabindex="-1"><a class="header-anchor" href="#api-参考"><span>API 参考</span></a></h2>
<h3 id="misakadbclient" tabindex="-1"><a class="header-anchor" href="#misakadbclient"><span>MisakaDBClient</span></a></h3>
<p>主要方法：</p>
<table>
<thead>
<tr>
<th>方法</th>
<th>说明</th>
</tr>
</thead>
<tbody>
<tr>
<td><code v-pre>connect(retries, retry_delay)</code></td>
<td>连接到服务器</td>
</tr>
<tr>
<td><code v-pre>close()</code></td>
<td>关闭连接</td>
</tr>
<tr>
<td><code v-pre>get_service_info()</code></td>
<td>获取服务信息</td>
</tr>
<tr>
<td><code v-pre>execute_command(command)</code></td>
<td>执行命令</td>
</tr>
<tr>
<td><code v-pre>ping()</code></td>
<td>心跳检测</td>
</tr>
<tr>
<td><code v-pre>get_server_version()</code></td>
<td>获取服务器版本</td>
</tr>
<tr>
<td><code v-pre>is_command_allowed(command)</code></td>
<td>检查命令是否允许</td>
</tr>
<tr>
<td><code v-pre>get_network_config()</code></td>
<td>获取网络配置</td>
</tr>
</tbody>
</table>
<h3 id="heartbeatcontroller" tabindex="-1"><a class="header-anchor" href="#heartbeatcontroller"><span>HeartbeatController</span></a></h3>
<p>心跳控制器，用于维护与服务器的连接状态。</p>
<table>
<thead>
<tr>
<th>属性</th>
<th>说明</th>
</tr>
</thead>
<tbody>
<tr>
<td><code v-pre>running</code></td>
<td>是否运行中</td>
</tr>
<tr>
<td><code v-pre>count</code></td>
<td>心跳总次数</td>
</tr>
<tr>
<td><code v-pre>success_count</code></td>
<td>成功次数</td>
</tr>
<tr>
<td><code v-pre>failure_count</code></td>
<td>失败次数</td>
</tr>
<tr>
<td><code v-pre>loss_rate</code></td>
<td>丢包率</td>
</tr>
</tbody>
</table>
<h2 id="许可证" tabindex="-1"><a class="header-anchor" href="#许可证"><span>许可证</span></a></h2>
<p>本项目基于 <code v-pre>木兰宽松许可证</code> 许可证开源。详见 <a href="./LICENSE">LICENSE</a> 文件。</p>
</div></template>


