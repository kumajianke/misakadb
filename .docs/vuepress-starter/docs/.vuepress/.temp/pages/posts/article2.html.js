import comp from "/Users/ahs/project-code/misaka_db/.docs/vuepress-starter/docs/.vuepress/.temp/pages/posts/article2.html.vue"
const data = JSON.parse("{\"path\":\"/posts/article2.html\",\"title\":\"原子任务中心\",\"lang\":\"en-US\",\"frontmatter\":{\"date\":\"2026-06-22T00:00:00.000Z\",\"title\":\"原子任务中心\",\"category\":[\"内部实现\"],\"tag\":[\"原子任务\"],\"description\":\"原子任务中心介绍\"},\"headers\":[],\"git\":{},\"filePathRelative\":\"posts/article2.md\",\"excerpt\":\"\\n<p>当客户端发来miql之后，misakadb要统一进行解析，为了保持miql的任务的原子性、一致性、隔离性、持久性这四个属性，misakaDB需要一个原子任务中心来管理所有的任务。</p>\\n<blockquote>\\n<p>如当用户发送一个删除任务，这个任务拆分下来就是需要对文件、JSON等进行删除的多个子任务。这些任务需要保证要么不执行要么全部执行。\\n所以我们删除任务不能直接进行删除，而是</p>\\n</blockquote>\\n\"}")
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
