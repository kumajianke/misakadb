import comp from "G:/project-pipline/misakadb/.docs/vuepress-starter/docs/.vuepress/.temp/pages/posts/article5.html.vue"
const data = JSON.parse("{\"path\":\"/posts/article5.html\",\"title\":\"Miql的数据库操作\",\"lang\":\"en-US\",\"frontmatter\":{\"date\":\"2026-08-09T00:00:00.000Z\",\"title\":\"Miql的数据库操作\",\"category\":[\"miql\"],\"tag\":[\"miql\"],\"description\":\"miql 对数据库的使用手册\"},\"headers\":[{\"level\":2,\"title\":\"MIQL\",\"slug\":\"miql\",\"link\":\"#miql\",\"children\":[{\"level\":3,\"title\":\"创建数据库\",\"slug\":\"创建数据库\",\"link\":\"#创建数据库\",\"children\":[]},{\"level\":3,\"title\":\"删除数据库\",\"slug\":\"删除数据库\",\"link\":\"#删除数据库\",\"children\":[]}]}],\"git\":{\"updatedTime\":1787545703000,\"contributors\":[{\"name\":\"库玛\",\"username\":\"\",\"email\":\"2791528600@qq.com\",\"commits\":1}],\"changelog\":[{\"hash\":\"c9d39bdfabcbe975366415aa8b08c94583ec4cf2\",\"time\":1787545703000,\"email\":\"2791528600@qq.com\",\"author\":\"库玛\",\"message\":\"修改了连续删除同一个数据库导致的目标重名，且无权删除的问题\"}]},\"filePathRelative\":\"posts/article5.md\",\"excerpt\":\"<h2>MIQL</h2>\\n<div class=\\\"hint-container tip\\\">\\n<p class=\\\"hint-container-title\\\">miql是什么</p>\\n<p>MIQL是用于增删改查数据库的语法,由mq.开头。mq 必须登录用户执行。</p>\\n</div>\\n<h3>创建数据库</h3>\\n<div class=\\\"language-miql line-numbers-mode\\\" data-highlighter=\\\"prismjs\\\" data-ext=\\\"miql\\\"><pre><code><span class=\\\"line\\\">mq.createDB(&lt;name&gt;[, engine])</span>\\n<span class=\\\"line\\\"></span></code></pre>\\n<div class=\\\"line-numbers\\\" aria-hidden=\\\"true\\\" style=\\\"counter-reset:line-number 0\\\"><div class=\\\"line-number\\\"></div></div></div>\"}")
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
