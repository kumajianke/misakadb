<template><div><h1 id="misakadb" tabindex="-1"><a class="header-anchor" href="#misakadb"><span>MisakaDB</span></a></h1>
<p><img src="@source/logos.png" alt=""></p>
<blockquote>
<p>便捷简单的文档数据库，始终如一。</p>
</blockquote>
<div class="hint-container warning">
<p class="hint-container-title">免责声明</p>
<p><strong>当前项目处于<code v-pre>α版本测试阶段</code>, 不建议直接用于生产环境。</strong></p>
</div>
<div class="hint-container tip">
<p class="hint-container-title">misaka的诞生</p>
<p>misaka的命名来自动漫《魔法禁书目录》、《某科学的超电池炮》的炮姐，御坂美琴中的御坂(misaka)及衍生的御坂网络，在后期我们也会仿造misaka来实现对应的分布式技术。</p>
</div>
<p><code v-pre>MisakaDB</code> 是一个轻量级的 JSON 文档数据库，支持 <code v-pre>JSON</code> 格式的内容字段存储。在 0.0.3 版本，仅支持单机模式。在0.17版本之后，支持插件模式，开发者可以在Misaka源码中开发自定义的插件，从而提高CRUD的效率。</p>
<h2 id="快速开始" tabindex="-1"><a class="header-anchor" href="#快速开始"><span>快速开始</span></a></h2>
<h3 id="客户端" tabindex="-1"><a class="header-anchor" href="#客户端"><span>客户端</span></a></h3>
<p>执行 client 目录下的 <code v-pre>misaka-cli</code> 程序或者利用源码执行 <code v-pre>uv run client/misaka-cli.py</code> 启动客户端。</p>
<div class="language-text line-numbers-mode" data-highlighter="prismjs" data-ext="text"><pre v-pre><code><span class="line">misaka-cli --mode shell</span>
<span class="line"></span></code></pre>
<div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0"><div class="line-number"></div></div></div><h4 id="支持参数" tabindex="-1"><a class="header-anchor" href="#支持参数"><span>支持参数</span></a></h4>
<ul>
<li><code v-pre>address</code> type=str default=&quot;127.0.0.1&quot; help=&quot;服务地址, 默认本地&quot;</li>
<li><code v-pre>port</code> type=int default=10032 help=&quot;服务端口, 默认10032&quot;</li>
<li><code v-pre>mode</code> type=str default=&quot;default&quot; help=&quot;运行模式: shell[cli调试模式]、onlyConn[仅连接]&quot;</li>
</ul>
<h3 id="miql" tabindex="-1"><a class="header-anchor" href="#miql"><span>miql</span></a></h3>
<div class="hint-container tip">
<p class="hint-container-title">miql</p>
<p>MIQL是用于增删改查数据库的语法, 适用于客户端SDK\Shell模式对数据库、数据表的内容操作。miql的指令由mq.开头。mq 必须登录用户执行。</p>
</div>
<h2 id="许可证" tabindex="-1"><a class="header-anchor" href="#许可证"><span>许可证</span></a></h2>
<p>本项目基于 <code v-pre>木兰宽松许可证</code> 许可证开源。详见 <a href="./LICENSE">LICENSE</a> 文件。</p>
</div></template>


