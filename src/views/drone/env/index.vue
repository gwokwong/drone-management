<!-- 环境设备运行状态 -->
<template>
  <div class="env-page art-full-height">
    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElForm :inline="true" :model="searchForm" @submit.prevent>
            <ElFormItem label="设备类型">
              <ElSelect v-model="searchForm.type" clearable placeholder="全部" style="width: 130px">
                <ElOption v-for="i in ENV_TYPE" :key="i.value" :label="i.label" :value="i.value" />
              </ElSelect>
            </ElFormItem>
            <ElFormItem label="运行状态">
              <ElSelect v-model="searchForm.status" clearable placeholder="全部" style="width: 120px">
                <ElOption v-for="i in ENV_STATUS" :key="i.value" :label="i.label" :value="i.value" />
              </ElSelect>
            </ElFormItem>
            <ElFormItem>
              <ElButton type="primary" @click="handleSearch">查询</ElButton>
              <ElButton @click="handleReset">重置</ElButton>
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
  </div>
</template>

<script setup lang="ts">
  import { useTable } from '@/hooks/core/useTable'
  import { fetchEnvDevices } from '@/api/drone'
  import type { EnvDevice } from '@/api/drone'
  import { ENV_TYPE, ENV_STATUS, labelOf, typeOf } from '@/utils/drone/dict'
  import { ElTag } from 'element-plus'

  defineOptions({ name: 'StorageEnv' })

  const searchForm = ref({ type: '', status: '' })

  /** 时间格式化 */
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
      apiFn: fetchEnvDevices,
      apiParams: { current: 1, size: 20 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'name', label: '设备名称', minWidth: 150 },
        {
          prop: 'type',
          label: '设备类型',
          width: 110,
          formatter: (row) =>
            h(ElTag, { type: typeOf(ENV_TYPE, row.type), size: 'small' }, () =>
              labelOf(ENV_TYPE, row.type)
            )
        },
        {
          prop: 'status',
          label: '运行状态',
          width: 100,
          formatter: (row) =>
            h(ElTag, { type: typeOf(ENV_STATUS, row.status), size: 'small' }, () =>
              labelOf(ENV_STATUS, row.status)
            )
        },
        {
          prop: 'value',
          label: '当前值',
          width: 110,
          formatter: (row) => `${row.value}${row.unit ? ' ' + row.unit : ''}`
        },
        { prop: 'roomId', label: '所属储存室', width: 120 },
        {
          prop: 'lastReportAt',
          label: '最后上报',
          minWidth: 165,
          formatter: (row) => formatTime(row.lastReportAt)
        },
        { prop: 'remark', label: '备注', minWidth: 140, showOverflowTooltip: true }
      ]
    }
  })

  const handleSearch = () => {
    replaceSearchParams({ ...searchForm.value })
    getData()
  }

  const handleReset = () => {
    searchForm.value = { type: '', status: '' }
    resetSearchParams()
    getData()
  }

  // 保持 refreshData 引用，供模板头部刷新按钮使用
  void refreshData
</script>
