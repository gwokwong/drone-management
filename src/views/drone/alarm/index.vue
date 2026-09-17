<!-- 告警管理 - 断网 / 漏水 / 门锁 / 设备等异常告警 -->
<template>
  <div class="alarm-page art-full-height">
    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElForm :inline="true" :model="searchForm" @submit.prevent>
            <ElFormItem label="告警类型">
              <ElSelect v-model="searchForm.type" clearable placeholder="全部" style="width: 140px">
                <ElOption v-for="i in ALARM_TYPE" :key="i.value" :label="i.label" :value="i.value" />
              </ElSelect>
            </ElFormItem>
            <ElFormItem label="级别">
              <ElSelect v-model="searchForm.level" clearable placeholder="全部" style="width: 110px">
                <ElOption v-for="i in ALARM_LEVEL" :key="i.value" :label="i.label" :value="i.value" />
              </ElSelect>
            </ElFormItem>
            <ElFormItem label="状态">
              <ElSelect v-model="searchForm.status" clearable placeholder="全部" style="width: 110px">
                <ElOption v-for="i in ALARM_STATUS" :key="i.value" :label="i.label" :value="i.value" />
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
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useTable } from '@/hooks/core/useTable'
  import { fetchAlarmList, resolveAlarm, ignoreAlarm, deleteAlarm } from '@/api/drone'
  import type { AlarmEvent } from '@/api/drone'
  import { ALARM_TYPE, ALARM_LEVEL, ALARM_STATUS, labelOf, typeOf } from '@/utils/drone/dict'
  import { ElTag, ElMessage, ElMessageBox, ElButton } from 'element-plus'

  defineOptions({ name: 'MonitorAlarm' })

  const searchForm = ref({ type: '', level: '', status: '' })

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
      apiFn: fetchAlarmList,
      apiParams: { current: 1, size: 20 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'title', label: '告警标题', minWidth: 170, showOverflowTooltip: true },
        {
          prop: 'type',
          label: '告警类型',
          width: 120,
          formatter: (row) =>
            h(ElTag, { type: typeOf(ALARM_TYPE, row.type), size: 'small' }, () =>
              labelOf(ALARM_TYPE, row.type)
            )
        },
        {
          prop: 'level',
          label: '级别',
          width: 90,
          formatter: (row) =>
            h(ElTag, { type: typeOf(ALARM_LEVEL, row.level), size: 'small' }, () =>
              labelOf(ALARM_LEVEL, row.level)
            )
        },
        {
          prop: 'status',
          label: '状态',
          width: 100,
          formatter: (row) =>
            h(ElTag, { type: typeOf(ALARM_STATUS, row.status), size: 'small' }, () =>
              labelOf(ALARM_STATUS, row.status)
            )
        },
        { prop: 'source', label: '来源', width: 110 },
        { prop: 'content', label: '告警内容', minWidth: 180, showOverflowTooltip: true },
        {
          prop: 'createdAt',
          label: '发生时间',
          width: 165,
          formatter: (row) => formatTime(row.createdAt)
        },
        {
          prop: 'operation',
          label: '操作',
          width: 190,
          fixed: 'right',
          formatter: (row) =>
            h('div', { class: 'flex items-center' }, [
              h(
                ElButton,
                {
                  link: true,
                  type: 'success',
                  size: 'small',
                  disabled: row.status !== 'active',
                  onClick: () => handleResolve(row)
                },
                () => '处理'
              ),
              h(
                ElButton,
                {
                  link: true,
                  type: 'info',
                  size: 'small',
                  disabled: row.status !== 'active',
                  onClick: () => handleIgnore(row)
                },
                () => '忽略'
              ),
              h(ArtButtonTable, { type: 'delete', onClick: () => handleDelete(row) })
            ])
        }
      ]
    }
  })

  /** 处理告警 */
  const handleResolve = (row: AlarmEvent) => {
    ElMessageBox.confirm(`确认已处理告警「${row.title}」？`, '告警处理', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      await resolveAlarm(row.id!)
      ElMessage.success('告警已处理')
      refreshData()
    })
  }

  /** 忽略告警 */
  const handleIgnore = (row: AlarmEvent) => {
    ElMessageBox.confirm(`确认忽略告警「${row.title}」？`, '告警忽略', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    }).then(async () => {
      await ignoreAlarm(row.id!)
      ElMessage.success('告警已忽略')
      refreshData()
    })
  }

  /** 删除告警 */
  const handleDelete = (row: AlarmEvent) => {
    ElMessageBox.confirm(`确定要删除告警「${row.title}」吗？`, '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'error'
    }).then(async () => {
      await deleteAlarm(row.id!)
      ElMessage.success('删除成功')
      refreshData()
    })
  }

  const handleSearch = () => {
    replaceSearchParams({ ...searchForm.value })
    getData()
  }

  const handleReset = () => {
    searchForm.value = { type: '', level: '', status: '' }
    resetSearchParams()
    getData()
  }
</script>
