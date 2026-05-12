import { defineConfig, presetUno, presetIcons, presetWebFonts } from 'unocss'

export default defineConfig({
  presets: [
    presetUno(),
    presetIcons({
      collections: {
        material: () => import('@iconify-json/material/icons.json').then(i => i.default)
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
  }
})
