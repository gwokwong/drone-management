<!-- 用户管理 -->
<template>
  <div class="user-page art-full-height">
    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElForm :inline="true" :model="searchForm" @submit.prevent>
            <ElFormItem label="用户名">
              <ElInput v-model="searchForm.userName" clearable style="width: 140px" />
            </ElFormItem>
            <ElFormItem label="昵称">
              <ElInput v-model="searchForm.nickName" clearable style="width: 140px" />
            </ElFormItem>
            <ElFormItem label="状态">
              <ElSelect v-model="searchForm.status" clearable placeholder="全部" style="width: 110px">
                <ElOption label="启用" :value="1" />
                <ElOption label="禁用" :value="0" />
              </ElSelect>
            </ElFormItem>
            <ElFormItem>
              <ElButton type="primary" @click="handleSearch">查询</ElButton>
              <ElButton @click="handleReset">重置</ElButton>
              <ElButton @click="showDialog('add')">新增用户</ElButton>
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
      :title="dialogType === 'add' ? '新增用户' : '编辑用户'"
      width="600"
    >
      <ElForm ref="formRef" :model="form" :rules="rules" label-width="90px">
        <ElFormItem label="用户名" prop="username">
          <ElInput v-model="form.username" :disabled="dialogType === 'edit'" />
        </ElFormItem>
        <ElFormItem v-if="dialogType === 'add'" label="初始密码">
          <ElInput v-model="form.password" type="password" show-password placeholder="留空默认 123456" />
        </ElFormItem>
        <ElFormItem label="昵称">
          <ElInput v-model="form.nickname" />
        </ElFormItem>
        <ElFormItem label="邮箱">
          <ElInput v-model="form.email" />
        </ElFormItem>
        <ElFormItem label="手机号">
          <ElInput v-model="form.phone" />
        </ElFormItem>
        <ElFormItem label="角色">
          <ElSelect v-model="form.roleCodes" multiple placeholder="请选择角色" style="width: 100%">
            <ElOption v-for="r in allRoles" :key="r.code" :label="r.name" :value="r.code" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="状态">
          <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用" />
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
  import {
    fetchUserList,
    fetchAllRoles,
    createUser,
    updateUser,
    deleteUser,
    resetPassword
  } from '@/api/drone'
  import type { UserItem, RoleItem } from '@/api/drone'
  import { ElTag, ElMessage, ElMessageBox, ElButton } from 'element-plus'
  import { DialogType } from '@/types'

  defineOptions({ name: 'SystemUser' })

  const searchForm = ref<{ userName?: string; nickName?: string; status?: number }>({})

  const dialogType = ref<DialogType>('add')
  const dialogVisible = ref(false)
  const submitting = ref(false)
  const formRef = ref()
  const form = ref<Partial<UserItem>>({})
  const allRoles = ref<RoleItem[]>([])

  /** 角色编码 -> 名称 */
  const roleNameMap = computed(() => {
    const m: Record<string, string> = {}
    allRoles.value.forEach((r) => {
      m[r.code] = r.name
    })
    return m
  })

  const rules = {
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }]
  }

  const formatTime = (t?: string) => (t ? t.replace('T', ' ').slice(0, 19) : '-')

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
      apiFn: fetchUserList,
      apiParams: { current: 1, size: 20 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'username', label: '用户名', width: 130 },
        { prop: 'nickname', label: '昵称', width: 130 },
        { prop: 'email', label: '邮箱', minWidth: 180 },
        { prop: 'phone', label: '手机号', width: 140 },
        {
          prop: 'roleCodes',
          label: '角色',
          minWidth: 160,
          formatter: (row) => {
            const codes = row.roleCodes ?? []
            if (!codes.length) return '-'
            return h(
              'div',
              { class: 'flex flex-wrap gap-1' },
              codes.map((c: string) =>
                h(ElTag, { size: 'small', type: 'primary' }, () => roleNameMap.value[c] || c)
              )
            )
          }
        },
        {
          prop: 'status',
          label: '状态',
          width: 90,
          formatter: (row) =>
            h(ElTag, { size: 'small', type: row.status === 1 ? 'success' : 'danger' }, () =>
              row.status === 1 ? '启用' : '禁用'
            )
        },
        {
          prop: 'createdAt',
          label: '创建时间',
          width: 170,
          formatter: (row) => formatTime(row.createdAt)
        },
        {
          prop: 'operation',
          label: '操作',
          width: 200,
          fixed: 'right',
          formatter: (row) =>
            h('div', { class: 'flex items-center' }, [
              h(ArtButtonTable, { type: 'edit', onClick: () => showDialog('edit', row) }),
              h(
                ElButton,
                { link: true, type: 'warning', size: 'small', onClick: () => handleResetPwd(row) },
                () => '重置密码'
              ),
              h(ArtButtonTable, { type: 'delete', onClick: () => handleDelete(row) })
            ])
        }
      ]
    }
  })

  /** 打开弹窗 */
  const showDialog = (type: DialogType, row?: UserItem) => {
    dialogType.value = type
    form.value = row
      ? { ...row, password: undefined, roleCodes: row.roleCodes ?? [] }
      : { status: 1, roleCodes: [] }
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
          await createUser(form.value)
          ElMessage.success('新增成功')
        } else {
          await updateUser(form.value.id!, form.value)
          ElMessage.success('更新成功')
        }
        dialogVisible.value = false
        refreshData()
      } finally {
        submitting.value = false
      }
    })
  }

  /** 重置密码 */
  const handleResetPwd = (row: UserItem) => {
    ElMessageBox.prompt(`请输入「${row.username}」的新密码`, '重置密码', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPlaceholder: '请输入新密码'
    }).then(async ({ value }) => {
      if (!value) {
        ElMessage.warning('密码不能为空')
        return
      }
      await resetPassword(row.id!, value)
      ElMessage.success('密码已重置')
    })
  }

  /** 删除 */
  const handleDelete = (row: UserItem) => {
    ElMessageBox.confirm(`确定要删除用户「${row.username}」吗？`, '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'error'
    }).then(async () => {
      await deleteUser(row.id!)
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

  /** 加载全部角色 */
  const loadRoles = async () => {
    try {
      allRoles.value = await fetchAllRoles()
    } catch {
      /* 拦截器已提示 */
    }
  }

  onMounted(loadRoles)
</script>
