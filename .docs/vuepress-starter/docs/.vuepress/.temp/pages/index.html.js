import comp from "/Users/ahs/project-code/misaka_db/.docs/vuepress-starter/docs/.vuepress/.temp/pages/index.html.vue"
const data = JSON.parse("{\"path\":\"/\",\"title\":\"主页\",\"lang\":\"en-US\",\"frontmatter\":{\"home\":true,\"title\":\"主页\",\"heroImage\":\"/logos.png\",\"heroImageDark\":\"/logos-dark.png\",\"actions\":[{\"text\":\"快速开发\",\"link\":\"/get-started.html\",\"type\":\"primary\"},{\"text\":\"下载\",\"link\":\"https://gitee.com/kumare/misakadb\",\"type\":\"secondary\"}],\"features\":[{\"title\":\"文档数据库支持\",\"details\":\"通过MiQL实现文档数据库的支持，支持增删改查等多个操作。\"},{\"title\":\"实现语言\",\"details\":\"Golang实现服务端、Python实现客户端。\"},{\"title\":\"全局锁池\",\"details\":\"实现全局锁池，确保在高并发环境下，每个事务只能被一个线程执行。\"}],\"footer\":\"MIT Licensed | Copyright © 库码工作室  & vuepress-starter\"},\"headers\":[],\"git\":{\"updatedTime\":1782178263000,\"contributors\":[{\"name\":\"hello_kuma\",\"username\":\"\",\"email\":\"2791528600@qq.com\",\"commits\":2}],\"changelog\":[{\"hash\":\"418d98b8ab2796f538e547cbada0e4acc21f5e01\",\"time\":1782178263000,\"email\":\"2791528600@qq.com\",\"author\":\"hello_kuma\",\"message\":\"更新文档内容\"},{\"hash\":\"8626ce43170c6b2ec42c38ae5c83146f2b849150\",\"time\":1782111363000,\"email\":\"2791528600@qq.com\",\"author\":\"hello_kuma\",\"message\":\"实现内部文档挂在到 misaka.kuamre.cn\"}]},\"filePathRelative\":\"README.md\"}")
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
