export const categoriesMap = JSON.parse("{\"category\":{\"/\":{\"path\":\"/category/\",\"map\":{\"内部实现\":{\"path\":\"/category/%E5%86%85%E9%83%A8%E5%AE%9E%E7%8E%B0/\",\"indexes\":[0,1,2]}}}},\"tag\":{\"/\":{\"path\":\"/tag/\",\"map\":{\"全局锁池\":{\"path\":\"/tag/%E5%85%A8%E5%B1%80%E9%94%81%E6%B1%A0/\",\"indexes\":[2]},\"原子任务\":{\"path\":\"/tag/%E5%8E%9F%E5%AD%90%E4%BB%BB%E5%8A%A1/\",\"indexes\":[1]},\"插件开发\":{\"path\":\"/tag/%E6%8F%92%E4%BB%B6%E5%BC%80%E5%8F%91/\",\"indexes\":[0]}}}}}");

if (import.meta.webpackHot) {
  import.meta.webpackHot.accept();
  if (__VUE_HMR_RUNTIME__.updateBlogCategory)
    __VUE_HMR_RUNTIME__.updateBlogCategory(categoriesMap);
}

if (import.meta.hot)
  import.meta.hot.accept(({ categoriesMap }) => {
    __VUE_HMR_RUNTIME__.updateBlogCategory(categoriesMap);
  });

