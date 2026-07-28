import comp from "G:/project-pipline/misakadb/.docs/vuepress-starter/docs/.vuepress/.temp/pages/posts/article3.html.vue"
const data = JSON.parse("{\"path\":\"/posts/article3.html\",\"title\":\"插件开发\",\"lang\":\"en-US\",\"frontmatter\":{\"date\":\"2026-07-02T00:00:00.000Z\",\"title\":\"插件开发\",\"category\":[\"插件\"],\"tag\":[\"插件开发\"],\"description\":\"插件开发规范\"},\"headers\":[{\"level\":2,\"title\":\"插件开发\",\"slug\":\"插件开发\",\"link\":\"#插件开发\",\"children\":[{\"level\":3,\"title\":\"快捷启动\",\"slug\":\"快捷启动\",\"link\":\"#快捷启动\",\"children\":[]},{\"level\":3,\"title\":\"任务模块\",\"slug\":\"任务模块\",\"link\":\"#任务模块\",\"children\":[]},{\"level\":3,\"title\":\"\",\"slug\":\"\",\"link\":\"#\",\"children\":[]}]}],\"git\":{\"updatedTime\":1785119984000,\"contributors\":[{\"name\":\"hello_kuma\",\"username\":\"\",\"email\":\"2791528600@qq.com\",\"commits\":2},{\"name\":\"库玛\",\"username\":\"\",\"email\":\"2791528600@qq.com\",\"commits\":1}],\"changelog\":[{\"hash\":\"bd0c16649dab7ef395e0c77a85eb8a00ba6a14bf\",\"time\":1785119984000,\"email\":\"2791528600@qq.com\",\"author\":\"库玛\",\"message\":\"update\"},{\"hash\":\"d9d5d8ad9e7e5383775bbf79e05751db206759d0\",\"time\":1784778255000,\"email\":\"2791528600@qq.com\",\"author\":\"hello_kuma\",\"message\":\"refactor(plugin-system): 重构插件注册系统并完善插件开发文档\"},{\"hash\":\"6d0d5738737eded2c5745c8bfe17c84049e71350\",\"time\":1782989277000,\"email\":\"2791528600@qq.com\",\"author\":\"hello_kuma\",\"message\":\"refactor(atomic-work-center): 重构任务中心字段命名并完善文档\"}]},\"filePathRelative\":\"posts/article3.md\",\"excerpt\":\"<div class=\\\"hint-container tip\\\">\\n<p class=\\\"hint-container-title\\\">Tips</p>\\n<p>misaka支持天然插件加载和开发，用户可以在 misaka的目录中的 <code>plugins/mods</code> 中开发自己的插件，官方默认加载了一个base的插件供以参考，请勿删除base插件，这会导致misaka不可用。</p>\\n</div>\\n<h2>插件开发</h2>\\n<p>Misaka支持插件的开发和使用。</p>\\n<div class=\\\"hint-container warning\\\">\\n<p class=\\\"hint-container-title\\\">Warning</p>\\n<p>以下设计尚在WIP中。</p>\\n</div>\"}")
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
