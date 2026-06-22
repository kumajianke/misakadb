import { blogPlugin } from '@vuepress/plugin-blog'
import { defaultTheme } from '@vuepress/theme-default'
import { defineUserConfig } from 'vuepress'
import { viteBundler } from '@vuepress/bundler-vite'
import { copyFile, cp, mkdir, readdir, stat } from 'node:fs/promises'
import path from 'node:path'

const IMAGE_EXTENSIONS = new Set([
  '.png',
  '.jpg',
  '.jpeg',
  '.gif',
  '.webp',
  '.svg',
  '.ico',
  '.bmp',
  '.avif',
])

const isImageFile = (fileName) => IMAGE_EXTENSIONS.has(path.extname(fileName).toLowerCase())

const copyRootImages = async (sourceDir, destDir) => {
  const entries = await readdir(sourceDir, { withFileTypes: true })

  await Promise.all(
    entries
      .filter((entry) => entry.isFile() && isImageFile(entry.name))
      .map((entry) =>
        copyFile(path.join(sourceDir, entry.name), path.join(destDir, entry.name)),
      ),
  )
}

const copyImgsDirectory = async (sourceDir, destDir) => {
  const imgsSourceDir = path.join(sourceDir, 'imgs')
  const imgsDestDir = path.join(destDir, 'imgs')

  try {
    const imgsStat = await stat(imgsSourceDir)
    if (!imgsStat.isDirectory()) return
  } catch {
    return
  }

  await mkdir(imgsDestDir, { recursive: true })
  await cp(imgsSourceDir, imgsDestDir, { recursive: true, force: true })
}

export default defineUserConfig({
  lang: 'en-US',

  title: 'Wiki of MisakaDB',
  description: 'misaka DB数据库的wiki',

  theme: defaultTheme({
    logo: '/logos.png',
    logoDark: '/logos-dark.png',

    navbar: [
      '/',
      {
        text: '文档',
        link: '/article/',
      },
      {
        text: '分类',
        link: '/category/',
      },
      {
        text: '模块',
        link: '/tag/',
      },

    ],
  }),

  plugins: [
    {
      name: 'copy-doc-images',
      onGenerated: async (app) => {
        const sourceDir = app.dir.source()
        const destDir = app.dir.dest()

        await copyRootImages(sourceDir, destDir)
        await copyImgsDirectory(sourceDir, destDir)
      },
    },
    blogPlugin({
      // Only files under posts are articles
      filter: ({ filePathRelative }) =>
        filePathRelative ? filePathRelative.startsWith('posts/') : false,

      // Getting article info
      getInfo: ({ frontmatter, title, data }) => ({
        title,
        author: frontmatter.author || '',
        date: frontmatter.date || null,
        category: frontmatter.category || [],
        tag: frontmatter.tag || [],
        excerpt:
          // Support manually set excerpt through frontmatter
          typeof frontmatter.excerpt === 'string'
            ? frontmatter.excerpt
            : data?.excerpt || '',
      }),

      // Generate excerpt for all pages excerpt those users choose to disable
      excerptFilter: ({ frontmatter }) =>
        !frontmatter.home &&
        frontmatter.excerpt !== false &&
        typeof frontmatter.excerpt !== 'string',

      category: [
        {
          key: 'category',
          getter: (page) => page.frontmatter.category || [],
          layout: 'Category',
          itemLayout: 'Category',
          frontmatter: () => ({
            title: 'Categories',
            sidebar: false,
          }),
          itemFrontmatter: (name) => ({
            title: `Category ${name}`,
            sidebar: false,
          }),
        },
        {
          key: 'tag',
          getter: (page) => page.frontmatter.tag || [],
          layout: 'Tag',
          itemLayout: 'Tag',
          frontmatter: () => ({
            title: 'Tags',
            sidebar: false,
          }),
          itemFrontmatter: (name) => ({
            title: `Tag ${name}`,
            sidebar: false,
          }),
        },
      ],

      type: [
        {
          key: 'article',
          // Remove archive articles
          filter: (page) => !page.frontmatter.archive,
          layout: 'Article',
          frontmatter: () => ({
            title: 'Articles',
            sidebar: false,
          }),
          // Sort pages with time and sticky
          sorter: (pageA, pageB) => {
            if (pageA.frontmatter.sticky && pageB.frontmatter.sticky)
              return pageB.frontmatter.sticky - pageA.frontmatter.sticky

            if (pageA.frontmatter.sticky && !pageB.frontmatter.sticky) return -1

            if (!pageA.frontmatter.sticky && pageB.frontmatter.sticky) return 1

            if (!pageB.frontmatter.date) return 1
            if (!pageA.frontmatter.date) return -1

            return (
              new Date(pageB.frontmatter.date).getTime() -
              new Date(pageA.frontmatter.date).getTime()
            )
          },
        },
        {
          key: 'timeline',
          // Only article with date should be added to timeline
          filter: (page) => page.frontmatter.date instanceof Date,
          // Sort pages with time
          sorter: (pageA, pageB) =>
            new Date(pageB.frontmatter.date).getTime() -
            new Date(pageA.frontmatter.date).getTime(),
          layout: 'Timeline',
          frontmatter: () => ({
            title: 'Timeline',
            sidebar: false,
          }),
        },
      ],
      hotReload: true,
    }),
  ],

  bundler: viteBundler(),
})
