<template><div><h1 id="misakadb" tabindex="-1"><a class="header-anchor" href="#misakadb"><span>MisakaDB</span></a></h1>
<p><img src="@source/logos.png" alt=""></p>
<blockquote>
<p>便捷简单的文档数据库，始终如一。</p>
</blockquote>
<div class="hint-container tip">
<p class="hint-container-title">Tips</p>
<p>misaka的命名来自动漫《魔法禁书目录》、《某科学的超电池炮》的炮姐，御坂美琴中的御坂(misaka)及衍生的御坂网络，在后期我们也会仿造misaka来实现对应的分布式技术。</p>
</div>
<p><code v-pre>MisakaDB</code> 是一个轻量级的 JSON 文档数据库，支持 <code v-pre>JSON</code> 格式的内容字段存储。在 0.0.3 版本，仅支持单机模式。在0.17版本之后，支持插件模式，开发者可以在Misaka源码中开发自定义的插件，从而提高CRUD的效率。</p>
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
<h3 id="设计价值" tabindex="-1"><a class="header-anchor" href="#设计价值"><span>设计价值</span></a></h3>
<ul>
<li>减少业务层重复创建锁对象的成本</li>
<li>让热点锁优先停留在高命中路径，提升锁获取效率</li>
<li>通过原子计数为后续无阻塞回收提供判断依据</li>
<li>为长期运行场景中的锁对象治理和内存回收预留扩展空间</li>
</ul>
<h2 id="项目结构" tabindex="-1"><a class="header-anchor" href="#项目结构"><span>项目结构</span></a></h2>
<div class="language-text line-numbers-mode" data-highlighter="prismjs" data-ext="text"><pre v-pre><code><span class="line">misaka_db/</span>
<span class="line">├── .docs/                  # VuePress 文档工程</span>
<span class="line">│   └── vuepress-starter/   # 文档站点源码与构建配置</span>
<span class="line">├── atomic/                 # 原子任务、事件总线与文件操作能力</span>
<span class="line">│   ├── EventBus/           # 原子任务事件总线</span>
<span class="line">│   ├── atomicFileHandler/  # 原子文件写入等操作</span>
<span class="line">│   └── atomicWorkCenter/   # 原子任务中心、序列化与任务类型</span>
<span class="line">├── bin/                    # 编译产物与运行时配置</span>
<span class="line">├── client/                 # Python 客户端 SDK 与打包脚本</span>
<span class="line">│   ├── apis/               # 客户端 API</span>
<span class="line">│   ├── interface/          # 状态与接口定义</span>
<span class="line">│   ├── mql/                # MQL 相关调用封装</span>
<span class="line">│   ├── network/            # 客户端网络通信</span>
<span class="line">│   └── usage/              # 使用示例</span>
<span class="line">├── clilog/                 # 日志输出模块</span>
<span class="line">├── command/                # 命令分发与工具命令实现</span>
<span class="line">├── config/                 # 配置加载与映射</span>
<span class="line">├── db-datas/               # 本地数据库数据目录</span>
<span class="line">├── engine/                 # 存储引擎与查询实现</span>
<span class="line">│   ├── Mson/               # MSON 解析</span>
<span class="line">│   ├── base/               # 引擎基础能力</span>
<span class="line">│   ├── db_list/            # 数据库列表管理</span>
<span class="line">│   ├── dispatch/           # 引擎分发工厂</span>
<span class="line">│   ├── share/              # MIQL 共享逻辑</span>
<span class="line">│   └── tinydb/             # TinyDB 核心、索引与数据集</span>
<span class="line">├── lock/                   # 全局锁池与锁前缀</span>
<span class="line">├── miusers/                # 用户管理</span>
<span class="line">├── network/                # 服务端网络层与注册中心</span>
<span class="line">├── profiles/               # 默认配置与用户数据</span>
<span class="line">├── safe/                   # 加密与安全模块</span>
<span class="line">├── shares/                 # 跨模块共享工具</span>
<span class="line">├── tools/                  # 工具程序入口</span>
<span class="line">├── tui/                    # 终端 UI 能力</span>
<span class="line">├── misaka.go               # 服务主入口</span>
<span class="line">└── go.mod                  # Go 模块定义</span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h2 id="快速开始" tabindex="-1"><a class="header-anchor" href="#快速开始"><span>快速开始</span></a></h2>
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


