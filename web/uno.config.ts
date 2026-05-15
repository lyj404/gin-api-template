import { defineConfig, presetUno, presetIcons, presetWebFonts } from 'unocss'

// 菜单与业务侧用到的图标（来源于数据库的动态字段，UnoCSS 静态扫描扫不到，需要 safelist）
const menuIcons = [
  'i-material-symbols:dashboard-outline',
  'i-material-symbols:group-outline',
  'i-material-symbols:person-outline',
  'i-material-symbols:manage-accounts-outline',
  'i-material-symbols:menu-outline',
  'i-material-symbols:corporate-fare-outline',
  'i-material-symbols:security-outline',
  'i-material-symbols:history-outline',
  'i-material-symbols:settings-outline',
  'i-material-symbols:home-outline',
  'i-material-symbols:notifications-outline',
  'i-material-symbols:mail-outline',
  'i-material-symbols:search',
  'i-material-symbols:edit-outline',
  'i-material-symbols:delete-outline',
  'i-material-symbols:add-circle-outline',
  'i-material-symbols:refresh',
  'i-material-symbols:download',
  'i-material-symbols:upload',
  'i-material-symbols:print',
  'i-material-symbols:share',
  'i-material-symbols:star-outline',
  'i-material-symbols:favorite-outline',
  'i-material-symbols:lock-outline',
  'i-material-symbols:visibility-outline',
  'i-material-symbols:map',
  'i-material-symbols:bar-chart',
  'i-material-symbols:pie-chart',
  'i-material-symbols:table',
  'i-material-symbols:calendar-month',
  'i-material-symbols:description',
  'i-material-symbols:folder-outline',
  'i-material-symbols:file-upload',
  'i-material-symbols:article-outline',
  'i-material-symbols:build-outline',
  'i-material-symbols:help-outline',
  'i-material-symbols:info-outline',
  'i-material-symbols:warning-outline',
  'i-material-symbols:check-circle-outline',
  'i-material-symbols:cloud-outline',
  'i-material-symbols:link',
  'i-material-symbols:tag',
  'i-material-symbols:label-outline',
  'i-material-symbols:category-outline',
  'i-material-symbols:layers',
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
