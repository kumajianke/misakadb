import comp from "/Users/ahs/project-code/misaka_db/.docs/vuepress-starter/docs/.vuepress/.temp/pages/posts/article3.html.vue"
const data = JSON.parse("{\"path\":\"/posts/article3.html\",\"title\":\"原子任务中心\",\"lang\":\"en-US\",\"frontmatter\":{\"date\":\"2026-07-02T00:00:00.000Z\",\"title\":\"原子任务中心\",\"category\":[\"内部实现\"],\"tag\":[\"插件开发\"],\"description\":\"插件开发\"},\"headers\":[{\"level\":2,\"title\":\"任务插件\",\"slug\":\"任务插件\",\"link\":\"#任务插件\",\"children\":[]}],\"git\":{},\"filePathRelative\":\"posts/article3.md\",\"excerpt\":\"<div class=\\\"hint-container tip\\\">\\n<p class=\\\"hint-container-title\\\">Tips</p>\\n<p>misaka支持天然插件加载和开发，用户可以在 misaka的目录中的 <code>plugins/mods</code> 中开发自己的插件，官方默认加载了一个base的插件供以参考，请勿删除base插件，这会导致misaka不可用。</p>\\n</div>\\n<h2>任务插件</h2>\\n\"}")
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
