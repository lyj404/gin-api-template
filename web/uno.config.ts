import { defineConfig, presetUno, presetIcons, presetWebFonts } from 'unocss'

// 菜单与业务侧用到的图标（来源于数据库的动态字段，UnoCSS 静态扫描扫不到，需要 safelist）
const menuIcons = [
  'i-material-symbols:dashboard-outline',
  'i-material-symbols:group-outline',
  'i-material-symbols:manage-accounts-outline',
  'i-material-symbols:list-alt-outline',
  'i-material-symbols:account-tree-outline',
  'i-material-symbols:shield-outline',
  'i-material-symbols:receipt-long-outline',
  'i-material-symbols:circle-outline'
]

export default defineConfig({
  presets: [
    presetUno(),
    presetIcons({
      scale: 1.2,
      collections: {
        'material-symbols': () => import('@iconify-json/material-symbols/icons.json').then(i => i.default)
      }
    }),
    presetWebFonts({
      fonts: {
        sans: 'Noto Sans SC:400,500,700'
      }
    })
  ],
  shortcuts: {
    'flex-center': 'flex justify-center items-center'
  },
  safelist: menuIcons
})
