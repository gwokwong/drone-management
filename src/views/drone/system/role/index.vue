<!-- 角色管理 - 不同角色分配不同权限 -->
<template>
  <div class="role-page art-full-height">
    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElForm :inline="true" :model="searchForm" @submit.prevent>
            <ElFormItem label="角色名称">
              <ElInput v-model="searchForm.name" clearable style="width: 140px" />
            </ElFormItem>
            <ElFormItem label="角色编码">
              <ElInput v-model="searchForm.code" clearable style="width: 140px" />
            </ElFormItem>
            <ElFormItem>
              <ElButton type="primary" @click="handleSearch">查询</ElButton>
              <ElButton @click="handleReset">重置</ElButton>
              <ElButton @click="showDialog('add')">新增角色</ElButton>
            </ElFormItem>
          </ElForm>
        </template>
      </ArtTableHeader>

      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      />
    </ElCard>

    <!-- 新增 / 编辑 -->
    <ElDialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增角色' : '编辑角色'"
      width="620"
    >
      <ElForm ref="formRef" :model="form" :rules="rules" label-width="90px">
        <ElFormItem label="角色编码" prop="code">
          <ElInput v-model="form.code" :disabled="dialogType === 'edit'" placeholder="如 R_OPERATOR" />
        </ElFormItem>
        <ElFormItem label="角色名称" prop="name">
          <ElInput v-model="form.name" placeholder="如 操作员" />
        </ElFormItem>
        <ElFormItem label="权限分配">
          <ElSelect
            v-model="form.permissions"
            multiple
            collapse-tags
            collapse-tags-tooltip
            placeholder="请选择权限"
            style="width: 100%"
          >
            <ElOption
              v-for="p in PERMISSION_OPTIONS"
              :key="p.value"
              :label="p.label"
              :value="p.value"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="启用状态">
          <ElSwitch
            v-model="form.enabled"
            :active-value="1"
            :inactive-value="0"
            active-text="启用"
            inactive-text="停用"
          />
        </ElFormItem>
        <ElFormItem label="备注">
          <ElInput v-model="form.remark" type="textarea" :rows="2" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitting" @click="handleSubmit">确定</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useTable } from '@/hooks/core/useTable'
  import { fetchRoleList, createRole, updateRole, deleteRole } from '@/api/drone'
  import type { RoleItem } from '@/api/drone'
  import { ElTag, ElMessage, ElMessageBox } from 'element-plus'
  import { DialogType } from '@/types'

  defineOptions({ name: 'SystemRole' })

  /** 可选权限清单 */
  const PERMISSION_OPTIONS = [
    { label: '数据大屏查看', value: 'dashboard:view' },
    { label: '无人机查看', value: 'drone:view' },
    { label: '无人机管理', value: 'drone:edit' },
    { label: '借还管理', value: 'borrow:manage' },
    { label: '储存室管理', value: 'storage:manage' },
    { label: '环境设备查看', value: 'env:view' },
    { label: '告警管理', value: 'alarm:manage' },
    { label: '人员进出查看', value: 'access:view' },
    { label: '视频监控管理', value: 'video:manage' },
    { label: '报表导出', value: 'report:export' },
    { label: '日志查看', value: 'log:view' },
    { label: '系统管理', value: 'system:manage' }
  ]

  /** 权限编码 -> 名称 */
  const permissionNameMap = computed(() => {
    const m: Record<string, string> = {}
    PERMISSION_OPTIONS.forEach((p) => {
      m[p.value] = p.label
    })
    return m
  })

  const searchForm = ref<{ name?: string; code?: string }>({})

  const dialogType = ref<DialogType>('add')
  const dialogVisible = ref(false)
  const submitting = ref(false)
  const formRef = ref()
  const form = ref<Partial<RoleItem>>({})

  const rules = {
    code: [{ required: true, message: '请输入角色编码', trigger: 'blur' }],
    name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }]
  }

  const {
    columns,
    columnChecks,
    data,
    loading,
    pagination,
    getData,
    replaceSearchParams,
    resetSearchParams,
    handleSizeChange,
    handleCurrentChange,
    refreshData
  } = useTable({
    core: {
      apiFn: fetchRoleList,
      apiParams: { current: 1, size: 20 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'code', label: '角色编码', width: 150 },
        { prop: 'name', label: '角色名称', width: 140 },
        {
          prop: 'permissions',
          label: '权限',
          minWidth: 260,
          formatter: (row) => {
            const perms = row.permissions ?? []
            if (!perms.length) return '-'
            return h(
              'div',
              { class: 'flex flex-wrap gap-1' },
              perms.map((p: string) =>
                h(ElTag, { size: 'small', type: 'info' }, () => permissionNameMap.value[p] || p)
              )
            )
          }
        },
        {
          prop: 'enabled',
          label: '状态',
          width: 90,
          formatter: (row) =>
            h(ElTag, { size: 'small', type: row.enabled === 1 ? 'success' : 'danger' }, () =>
              row.enabled === 1 ? '启用' : '停用'
            )
        },
        { prop: 'remark', label: '备注', minWidth: 160, showOverflowTooltip: true },
        {
          prop: 'operation',
          label: '操作',
          width: 120,
          fixed: 'right',
          formatter: (row) =>
            h('div', [
              h(ArtButtonTable, { type: 'edit', onClick: () => showDialog('edit', row) }),
              h(ArtButtonTable, { type: 'delete', onClick: () => handleDelete(row) })
            ])
        }
      ]
    }
  })

  /** 打开弹窗 */
  const showDialog = (type: DialogType, row?: RoleItem) => {
    dialogType.value = type
    form.value = row
      ? { ...row, permissions: row.permissions ?? [] }
      : { enabled: 1, permissions: [] }
    nextTick(() => {
      dialogVisible.value = true
    })
  }

  /** 提交 */
  const handleSubmit = async () => {
    await formRef.value?.validate(async (valid: boolean) => {
      if (!valid) return
      submitting.value = true
      try {
        if (dialogType.value === 'add') {
          await createRole(form.value)
          ElMessage.success('新增成功')
        } else {
          await updateRole(form.value.id!, form.value)
          ElMessage.success('更新成功')
        }
        dialogVisible.value = false
        refreshData()
      } finally {
        submitting.value = false
      }
    })
  }

  /** 删除 */
  const handleDelete = (row: RoleItem) => {
    ElMessageBox.confirm(`确定要删除角色「${row.name}」吗？`, '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'error'
    }).then(async () => {
      await deleteRole(row.id!)
      ElMessage.success('删除成功')
      refreshData()
    })
  }

  const handleSearch = () => {
    replaceSearchParams({ ...searchForm.value })
    getData()
  }

  const handleReset = () => {
    searchForm.value = {}
    resetSearchParams()
    getData()
  }
</script>
