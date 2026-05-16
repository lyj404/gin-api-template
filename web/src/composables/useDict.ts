import { ref, onMounted, type Ref } from 'vue'
import { getDictInfo } from '@/api'
import type { DictDetailResponse } from '@/types'

const cache = new Map<string, DictDetailResponse[]>()

export function useDict(dictType: string) {
  const options = ref<{ label: string; value: string }[]>([])
  const loading = ref(false)

  const lookup = (value: string): string => {
    return options.value.find(o => o.value === value)?.label || value
  }

  const load = async () => {
    if (cache.has(dictType)) {
      options.value = cache.get(dictType)!.map(d => ({ label: d.label, value: d.value }))
      return
    }
    loading.value = true
    try {
      const res = await getDictInfo(dictType)
      const items = res.data.data || []
      cache.set(dictType, items)
      options.value = items.map(d => ({ label: d.label, value: d.value }))
    } catch {
      options.value = []
    } finally {
      loading.value = false
    }
  }

  return { options, loading, lookup, load }
}
