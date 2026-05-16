<template>
  <div class="page-padding">
    <div class="toolbar-row mb-3">
      <n-h2 class="!my-0">组织管理</n-h2>
      <n-button type="primary" @click="openModal(null)">新增组织</n-button>
    </div>

    <n-card>
      <n-tree
        :data="treeData"
        :loading="loading"
        key-field="id"
        label-field="label"
        children-field="children"
        default-expand-all
        selectable
        :render-label="renderNodeLabel"
      />
    </n-card>

    <n-modal v-model:show="showModal" preset="card" :title="editingId ? '编辑组织' : '新增组织'" :style="{ width: '90vw', maxWidth: '500px' }">
      <n-form :model="form" label-placement="left" label-width="80">
        <n-form-item label="组织名称">
          <n-input v-model:value="form.name" placeholder="请输入组织名称" />
        </n-form-item>
        <n-form-item label="上级组织">
          <n-select v-model:value="form.parent_id" :options="orgOptions" placeholder="请选择上级组织" clearable />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" :loading="saving" class="ml-8" @click="handleSave">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { NButton, NTree, NCard, NModal, NForm, NFormItem, NInput, NSelect, NH2, NSpace, useMessage } from 'naive-ui'
import type { TreeOption, SelectOption } from 'naive-ui'
import { getOrgTree, createOrgUnit, updateOrgUnit, deleteOrgUnit } from '@/api'
import type { OrgUnitResponse } from '@/types'

const message = useMessage()
const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const editingId = ref<string | null>(null)

const form = reactive({ name: '', parent_id: null as string | null })
const allOrgs = ref<OrgUnitResponse[]>([])

const buildTree = (nodes: OrgUnitResponse[], parentId: string | null = null): TreeOption[] => {
  return nodes
    .filter(n => n.parent_id === parentId)
    .map(n => ({
      id: n.id,
      label: n.name,
      key: n.id,
      children: buildTree(nodes, n.id)
    }))
}

const treeData = computed(() => buildTree(allOrgs.value))

const orgOptions = computed<SelectOption[]>(() => [{ label: '顶级组织', value: '0' }, ...allOrgs.value.map((o: OrgUnitResponse) => ({ label: o.name, value: o.id }))])

const renderNodeLabel = (info: { option: TreeOption }) => {
  const opt = info.option
  return h('div', { style: 'display: flex; align-items: center; gap: 8px; width: 100%; padding: 4px 0' }, [
    h('span', { style: 'flex: 1' }, opt.label as string),
    h(NButton, { size: 'tiny', quaternary: true, onClick: (e: Event) => { e.stopPropagation(); openModal(opt) } }, { default: () => '新增子级' }),
    h(NButton, { size: 'tiny', quaternary: true, onClick: (e: Event) => { e.stopPropagation(); openModal(opt, true) } }, { default: () => '编辑' }),
    h(NButton, { size: 'tiny', quaternary: true, type: 'error', onClick: (e: Event) => { e.stopPropagation(); handleDelete(opt) } }, { default: () => '删除' })
  ])
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getOrgTree()
    allOrgs.value = res.data.data || []
  } catch {} finally {
    loading.value = false
  }
}

const openModal = (node: TreeOption | null, isEdit?: boolean) => {
  if (node && isEdit) {
    const org = allOrgs.value.find(o => o.id === node.id)
    editingId.value = node.id as string
    form.name = (org?.name || node.label) as string
    form.parent_id = org?.parent_id || null
  } else if (node) {
    editingId.value = null
    form.name = ''
    form.parent_id = node.id as string
  } else {
    editingId.value = null
    form.name = ''
    form.parent_id = null
  }
  showModal.value = true
}

const handleSave = async () => {
  saving.value = true
  try {
    const payload: any = { name: form.name }
    if (form.parent_id) payload.parent_id = form.parent_id
    if (editingId.value) {
      await updateOrgUnit(editingId.value, payload)
    } else {
      await createOrgUnit(payload)
    }
    message.success(editingId.value ? '更新成功' : '创建成功')
    showModal.value = false
    fetchData()
  } catch (e: any) {
    message.error(e?.response?.data?.message || '操作失败')
  } finally {
    saving.value = false
  }
}

const handleDelete = (opt: TreeOption) => {
  deleteOrgUnit(opt.id as string).then(() => { message.success('删除成功'); fetchData() }).catch((e: any) => message.error(e?.response?.data?.message || '删除失败'))
}

onMounted(fetchData)
</script>