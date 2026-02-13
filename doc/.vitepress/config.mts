import { defineConfig } from 'vitepress'

export default defineConfig({
  lang: 'en-US',
  title: 'Alias Manager',
  description: 'Fast, safe CLI for managing shell aliases.',
  lastUpdated: true,
  cleanUrls: true,
  themeConfig: {
    nav: [
      { text: 'Getting Started', link: '/getting-started' },
      { text: 'Usage', link: '/usage' },
      { text: 'Internals', link: '/internals' },
      { text: 'Development', link: '/development' },
      { text: 'Commands', link: '/commands' }
    ],
    sidebar: [
      {
        text: 'Guide',
        items: [
          { text: 'Getting Started', link: '/getting-started' },
          { text: 'Usage', link: '/usage' },
          { text: 'Internals & Safety', link: '/internals' }
        ]
      },
      {
        text: 'Reference',
        items: [
          { text: 'Command Reference', link: '/commands' },
          { text: 'Development', link: '/development' }
        ]
      }
    ],
    socialLinks: [
      { icon: 'github', link: 'https://github.com/navio/am' }
    ],
    outline: [2, 3]
  }
})
