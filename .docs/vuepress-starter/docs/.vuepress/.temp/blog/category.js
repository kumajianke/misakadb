export const categoriesMap = JSON.parse("{\"category\":{\"/\":{\"path\":\"/category/\",\"map\":{\"内部实现\":{\"path\":\"/category/%E5%86%85%E9%83%A8%E5%AE%9E%E7%8E%B0/\",\"indexes\":[0,1]},\"插件\":{\"path\":\"/category/%E6%8F%92%E4%BB%B6/\",\"indexes\":[2]},\"其他\":{\"path\":\"/category/%E5%85%B6%E4%BB%96/\",\"indexes\":[3]}}}},\"tag\":{\"/\":{\"path\":\"/tag/\",\"map\":{\"全局锁池\":{\"path\":\"/tag/%E5%85%A8%E5%B1%80%E9%94%81%E6%B1%A0/\",\"indexes\":[1]},\"原子任务\":{\"path\":\"/tag/%E5%8E%9F%E5%AD%90%E4%BB%BB%E5%8A%A1/\",\"indexes\":[0]},\"插件开发\":{\"path\":\"/tag/%E6%8F%92%E4%BB%B6%E5%BC%80%E5%8F%91/\",\"indexes\":[2]},\"更新\":{\"path\":\"/tag/%E6%9B%B4%E6%96%B0/\",\"indexes\":[3]}}}}}");

if (import.meta.webpackHot) {
  import.meta.webpackHot.accept();
  if (__VUE_HMR_RUNTIME__.updateBlogCategory)
    __VUE_HMR_RUNTIME__.updateBlogCategory(categoriesMap);
}

if (import.meta.hot)
  import.meta.hot.accept(({ categoriesMap }) => {
    __VUE_HMR_RUNTIME__.updateBlogCategory(categoriesMap);
  });

