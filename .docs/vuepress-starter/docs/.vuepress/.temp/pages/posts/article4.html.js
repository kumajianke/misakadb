import comp from "G:/project-pipline/misakadb/.docs/vuepress-starter/docs/.vuepress/.temp/pages/posts/article4.html.vue"
const data = JSON.parse("{\"path\":\"/posts/article4.html\",\"title\":\"更新日志\",\"lang\":\"en-US\",\"frontmatter\":{\"date\":\"2026-08-09T00:00:00.000Z\",\"title\":\"更新日志\",\"category\":[\"其他\"],\"tag\":[\"更新\"],\"description\":\"一些更新导致的改动\"},\"headers\":[{\"level\":2,\"title\":\"V0.1.8\",\"slug\":\"v0-1-8\",\"link\":\"#v0-1-8\",\"children\":[]}],\"git\":{},\"filePathRelative\":\"posts/article4.md\",\"excerpt\":\"<h2>V0.1.8</h2>\\n<div class=\\\"hint-container tip\\\">\\n<p class=\\\"hint-container-title\\\">常规更新</p>\\n<h3>feature</h3>\\n<ul>\\n<li>新增<code>TaskCombo</code>;</li>\\n<li>实现删除数据库逻辑;</li>\\n</ul>\\n<h3>chore</h3>\\n<ul>\\n<li>加入更新日志</li>\\n</ul>\\n<h3>refactor</h3>\\n<ul>\\n<li>修改获取PluginX单例的函数名: <code>GetPluginBus</code> -&gt; <code>GetPluginX</code>;</li>\\n<li>修改<code>PluginBus</code>名称为 <code>PluginsX</code>;</li>\\n</ul>\\n</div>\"}")
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
