import comp from "G:/project-pipline/misakadb/.docs/vuepress-starter/docs/.vuepress/.temp/pages/posts/article3.html.vue"
const data = JSON.parse("{\"path\":\"/posts/article3.html\",\"title\":\"插件开发\",\"lang\":\"en-US\",\"frontmatter\":{\"date\":\"2026-07-02T00:00:00.000Z\",\"title\":\"插件开发\",\"category\":[\"插件\"],\"tag\":[\"插件开发\"],\"description\":\"插件开发规范\"},\"headers\":[{\"level\":2,\"title\":\"插件开发\",\"slug\":\"插件开发\",\"link\":\"#插件开发\",\"children\":[{\"level\":3,\"title\":\"快捷启动\",\"slug\":\"快捷启动\",\"link\":\"#快捷启动\",\"children\":[]},{\"level\":3,\"title\":\"插件的安装\",\"slug\":\"插件的安装\",\"link\":\"#插件的安装\",\"children\":[]},{\"level\":3,\"title\":\"常见插件函数\",\"slug\":\"常见插件函数\",\"link\":\"#常见插件函数\",\"children\":[]},{\"level\":3,\"title\":\"常见插件术语\",\"slug\":\"常见插件术语\",\"link\":\"#常见插件术语\",\"children\":[]}]}],\"git\":{\"updatedTime\":1786240905000,\"contributors\":[{\"name\":\"hello_kuma\",\"username\":\"\",\"email\":\"2791528600@qq.com\",\"commits\":2},{\"name\":\"库玛\",\"username\":\"\",\"email\":\"2791528600@qq.com\",\"commits\":5}],\"changelog\":[{\"hash\":\"83bb2a6d309049c38bd1f678081e3169d712184a\",\"time\":1786240905000,\"email\":\"2791528600@qq.com\",\"author\":\"库玛\",\"message\":\"v0.1.8  更新 WIP\"},{\"hash\":\"0b629a822e3705d0b16d359e177ceba3b8e43d3a\",\"time\":1785333064000,\"email\":\"2791528600@qq.com\",\"author\":\"库玛\",\"message\":\"fixed:[pluginsX]统一命名为pluginsX，将Alias管理交给pluginsX update:[docs]追加内容解释\"},{\"hash\":\"cc9f1dd792633a672d57671a08ad371887bd35bc\",\"time\":1785280856000,\"email\":\"2791528600@qq.com\",\"author\":\"库玛\",\"message\":\"编译部分内容\"},{\"hash\":\"131e06e99c5e80006e8afca62830be8acbf8e718\",\"time\":1785198545000,\"email\":\"2791528600@qq.com\",\"author\":\"库玛\",\"message\":\"修复了Windows环境下，部分tui的问题。\"},{\"hash\":\"bd0c16649dab7ef395e0c77a85eb8a00ba6a14bf\",\"time\":1785119984000,\"email\":\"2791528600@qq.com\",\"author\":\"库玛\",\"message\":\"update\"},{\"hash\":\"d9d5d8ad9e7e5383775bbf79e05751db206759d0\",\"time\":1784778255000,\"email\":\"2791528600@qq.com\",\"author\":\"hello_kuma\",\"message\":\"refactor(plugin-system): 重构插件注册系统并完善插件开发文档\"},{\"hash\":\"6d0d5738737eded2c5745c8bfe17c84049e71350\",\"time\":1782989277000,\"email\":\"2791528600@qq.com\",\"author\":\"hello_kuma\",\"message\":\"refactor(atomic-work-center): 重构任务中心字段命名并完善文档\"}]},\"filePathRelative\":\"posts/article3.md\",\"excerpt\":\"<div class=\\\"hint-container tip\\\">\\n<p class=\\\"hint-container-title\\\">Tips</p>\\n<p>misaka支持天然插件加载和开发，用户可以在 misaka 开发自己的插件，官方默认加载了一个misaka basic mode的插件供以参考，请勿删除misaka basic mode插件，这会导致misaka不可用。</p>\\n</div>\\n<h2>插件开发</h2>\\n<p>Misaka支持插件的开发和使用。</p>\\n<p><strong>模组原理</strong>：开发者创建函数，然后按照指定的type约束传入注册机，如：<code>func ([]tasktype.TaskType) []tasktype.TaskType</code>, 注册机会根据类型在程序启动之后，在相应的地方调用函数。</p>\"}")
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
