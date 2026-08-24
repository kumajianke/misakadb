import comp from "G:/project-pipline/misakadb/.docs/vuepress-starter/docs/.vuepress/.temp/pages/posts/article4.html.vue"
const data = JSON.parse("{\"path\":\"/posts/article4.html\",\"title\":\"更新日志\",\"lang\":\"en-US\",\"frontmatter\":{\"date\":\"2026-08-09T00:00:00.000Z\",\"title\":\"更新日志\",\"category\":[\"其他\"],\"tag\":[\"更新\"],\"description\":\"一些更新导致的改动\"},\"headers\":[{\"level\":2,\"title\":\"V0.1.8\",\"slug\":\"v0-1-8\",\"link\":\"#v0-1-8\",\"children\":[]}],\"git\":{\"updatedTime\":1787531896000,\"contributors\":[{\"name\":\"库玛\",\"username\":\"\",\"email\":\"2791528600@qq.com\",\"commits\":3}],\"changelog\":[{\"hash\":\"a1f23ea1d32342a1d95fdc1d98aeb58470446935\",\"time\":1787531896000,\"email\":\"2791528600@qq.com\",\"author\":\"库玛\",\"message\":\"fix &amp; refactor: 修复Windows平台bug，重构数据库删除逻辑\"},{\"hash\":\"c51d27b9d35cba5be4c8115a8b269575166c3431\",\"time\":1786252207000,\"email\":\"2791528600@qq.com\",\"author\":\"库玛\",\"message\":\"TaskCombo的注册和Base插件的Combo提交\"},{\"hash\":\"83bb2a6d309049c38bd1f678081e3169d712184a\",\"time\":1786240905000,\"email\":\"2791528600@qq.com\",\"author\":\"库玛\",\"message\":\"v0.1.8  更新 WIP\"}]},\"filePathRelative\":\"posts/article4.md\",\"excerpt\":\"<h2>V0.1.8</h2>\\n<div class=\\\"hint-container tip\\\">\\n<p class=\\\"hint-container-title\\\">常规更新</p>\\n<h3>fixed</h3>\\n<ul>\\n<li>Windows环境中，强制落盘目录导致无权限的Bug;</li>\\n<li>Windows环境中，部分场景下，删除数据库无权限的提醒。</li>\\n</ul>\\n<h3>feature</h3>\\n<ul>\\n<li>实现 <code>TaskCombo</code>、<code>AddTaskCombo</code> 等业务函数;</li>\\n<li>实现删除数据库逻辑;</li>\\n</ul>\\n<h3>chore</h3>\\n<ul>\\n<li>加入更新日志</li>\\n</ul>\\n<h3>refactor</h3>\\n<ul>\\n<li>修改获取PluginX单例的函数名: <code>GetPluginBus</code> -&gt; <code>GetPluginX</code>;</li>\\n<li>修改 <code>PluginBus</code> 名称为 <code>PluginsX</code>;</li>\\n<li>修改 内置插件 <code>base_unloader</code> 插件名称为 <code>coin</code>;</li>\\n<li>将 <code>AtomicDropDB</code> 的逻辑转交给 内置插件 <code>coin</code> 执行；</li>\\n</ul>\\n<h3>下个版本代办</h3>\\n<p>功能\\t状态\\n正常删除目录\\t部分完成，实际是改名\\n正常失败回滚\\t基本有\\n删除后的物理清理\\t未完成\\nalldb.list 同步删除\\t未完成\\n程序中断后加载任务\\t已完成\\n程序启动后自动继续任务\\t未完成\\n删除动作中断后的准确恢复\\t未完成\\n回滚失败后的任务保留\\t未完成\\n成功任务清理\\t未完成\\nWindows 目录同步\\t已处理</p>\\n</div>\"}")
export { comp, data }

if (import.meta.webpackHot) {
  import.meta.webpackHot.accept()
  if (__VUE_HMR_RUNTIME__.updatePageData) {
    __VUE_HMR_RUNTIME__.updatePageData(data)
  }
}

if (import.meta.hot) {
  import.meta.hot.accept(({ data }) => {
    __VUE_HMR_RUNTIME__.updatePageData(data)
  })
}
