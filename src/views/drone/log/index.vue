<!-- 系统日志 - 按设备分类 / 时间 / 关键字查询，支持导出 Excel -->
<template>
  <div class="log-page art-full-height">
    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElForm :inline="true" :model="searchForm" @submit.prevent>
            <ElFormItem label="日志分类">
              <ElSelect
                v-model="searchForm.category"
                clearable
                placeholder="全部"
                style="width: 120px"
              >
                <ElOption v-for="i in LOG_CATEGORY" :key="i.value" :label="i.label" :value="i.value" />
              </ElSelect>
            </ElFormItem>
            <ElFormItem label="设备类型">
              <ElInput v-model="searchForm.deviceType" clearable style="width: 120px" />
            </ElFormItem>
            <ElFormItem label="级别">
              <ElSelect v-model="searchForm.level" clearable placeholder="全部" style="width: 110px">
                <ElOption label="信息" value="info" />
                <ElOption label="警告" value="warning" />
                <ElOption label="错误" value="error" />
              </ElSelect>
            </ElFormItem>
            <ElFormItem label="关键字">
              <ElInput
                v-model="searchForm.keyword"
                placeholder="模糊查询内容/操作人"
                clearable
                style="width: 180px"
              />
            </ElFormItem>
            <ElFormItem label="时间范围">
              <ElDatePicker
                v-model="dateRange"
                type="datetimerange"
                start-placeholder="开始时间"
                end-placeholder="结束时间"
                value-format="YYYY-MM-DDTHH:mm:ssZ"
                style="width: 340px"
              />
            </ElFormItem>
            <ElFormItem>
              <ElButton type="primary" @click="handleSearch">查询</ElButton>
              <ElButton @click="handleReset">重置</ElButton>
              <ElButton :loading="exporting" @click="handleExport">导出 Excel</ElButton>
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
  import { fetchLogList, exportLogs } from '@/api/drone'
  import type { SystemLog } from '@/api/drone'
  import { LOG_CATEGORY, labelOf, typeOf } from '@/utils/drone/dict'
  import { ElTag, ElMessage } from 'element-plus'

  defineOptions({ name: 'LogIndex' })

  const searchForm = ref({
    category: '',
    deviceType: '',
    level: '',
    keyword: ''
  })
  const dateRange = ref<string[]>([])
  const exporting = ref(false)

  const formatTime = (t?: string) => (t ? t.replace('T', ' ').slice(0, 19) : '-')

  /** 构建查询参数（含时间范围） */
  const buildParams = () => ({
    ...searchForm.value,
    start: dateRange.value?.[0] ?? '',
    end: dateRange.value?.[1] ?? ''
  })

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
      apiFn: fetchLogList,
      apiParams: { current: 1, size: 20 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        {
          prop: 'category',
          label: '日志分类',
          width: 110,
          formatter: (row) =>
            h(ElTag, { type: typeOf(LOG_CATEGORY, row.category), size: 'small' }, () =>
              labelOf(LOG_CATEGORY, row.category)
            )
        },
        { prop: 'deviceType', label: '设备类型', width: 120 },
        { prop: 'content', label: '日志内容', minWidth: 260, showOverflowTooltip: true },
        { prop: 'operator', label: '操作人', width: 120 },
        {
          prop: 'level',
          label: '级别',
          width: 90,
          formatter: (row) =>
            h(
              ElTag,
              {
                size: 'small',
                type: row.level === 'error' ? 'danger' : row.level === 'warning' ? 'warning' : 'info'
              },
              () => (row.level === 'error' ? '错误' : row.level === 'warning' ? '警告' : '信息')
            )
        },
        { prop: 'ip', label: 'IP 地址', width: 130 },
        {
          prop: 'createdAt',
          label: '产生时间',
          width: 170,
          formatter: (row) => formatTime(row.createdAt)
        }
      ]
    }
  })

  const handleSearch = () => {
    replaceSearchParams(buildParams())
    getData()
  }

  const handleReset = () => {
    searchForm.value = { category: '', deviceType: '', level: '', keyword: '' }
    dateRange.value = []
    resetSearchParams()
    getData()
  }

  /** 导出 Excel */
  const handleExport = async () => {
    exporting.value = true
    try {
      await exportLogs(buildParams())
      ElMessage.success('日志已导出')
    } catch {
      ElMessage.error('导出失败')
    } finally {
      exporting.value = false
    }
  }

  void refreshData
</script>
