<template><div><h1 id="原子任务中心的作用" tabindex="-1"><a class="header-anchor" href="#原子任务中心的作用"><span>原子任务中心的作用</span></a></h1>
<p>当客户端发来miql之后，misakadb要统一进行解析，为了保持miql的任务的原子性、一致性、隔离性、持久性这四个属性，misakaDB需要一个原子任务中心来管理所有的任务。</p>
<blockquote>
<p>如当用户发送一个删除任务，这个任务拆分下来就是需要对文件、JSON等进行删除的多个子任务。这些任务需要保证要么不执行要么全部执行。</p>
</blockquote>
</div></template>


