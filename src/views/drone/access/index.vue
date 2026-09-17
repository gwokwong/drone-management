<!-- 人员进出 - 刷卡/人脸认证与门锁事件统计 -->
<template>
  <div class="access-page art-full-height">
    <!-- 统计卡片 -->
    <ElCard shadow="never" class="stats-card">
      <div class="stats-row">
        <div v-for="s in statItems" :key="s.label" class="stat-item">
          <span class="stat-label">{{ s.label }}</span>
          <span class="stat-value" :style="{ color: s.color }">{{ s.value }}</span>
        </div>
      </div>
    </ElCard>

    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElForm :inline="true" :model="searchForm" @submit.prevent>
            <ElFormItem label="事件类型">
              <ElSelect v-model="searchForm.type" clearable placeholder="全部" style="width: 150px">
                <ElOption v-for="i in ACCESS_TYPE" :key="i.value" :label="i.label" :value="i.value" />
              </ElSelect>
            </ElFormItem>
            <ElFormItem label="人员">
              <ElInput v-model="searchForm.personName" clearable style="width: 130px" />
            </ElFormItem>
            <ElFormItem label="卡号">
              <ElInput v-model="searchForm.cardNo" clearable style="width: 140px" />
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
  import { fetchAccessList, fetchAccessStats } from '@/api/drone'
  import type { AccessRecord } from '@/api/drone'
  import { ACCESS_TYPE, labelOf, typeOf } from '@/utils/drone/dict'
  import { ElTag } from 'element-plus'

  defineOptions({ name: 'MonitorAccess' })

  const searchForm = ref({ type: '', personName: '', cardNo: '' })

  /** 进出统计 */
  const stats = ref<Record<string, number>>({})
  const statItems = computed(() => [
    { label: '累计进出', value: stats.value.total ?? 0, color: '#1D84FF' },
    { label: '今日进出', value: stats.value.today ?? 0, color: '#60C041' },
    { label: '卡号认证失败', value: stats.value.cardFail ?? 0, color: '#F56C6C' },
    { label: '门锁打开', value: stats.value.doorOpen ?? 0, color: '#60C041' },
    { label: '门锁关闭', value: stats.value.doorClose ?? 0, color: '#909399' },
    { label: '人脸认证通过', value: stats.value.facePass ?? 0, color: '#38C0FC' }
  ])

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
      apiFn: fetchAccessList,
      apiParams: { current: 1, size: 20 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'personName', label: '人员姓名', width: 130 },
        { prop: 'cardNo', label: '卡号', width: 150 },
        {
          prop: 'type',
          label: '事件类型',
          width: 150,
          formatter: (row) =>
            h(ElTag, { type: typeOf(ACCESS_TYPE, row.type), size: 'small' }, () =>
              labelOf(ACCESS_TYPE, row.type)
            )
        },
        { prop: 'roomId', label: '所属储存室', width: 120 },
        { prop: 'deviceId', label: '设备ID', width: 100 },
        {
          prop: 'createdAt',
          label: '发生时间',
          minWidth: 170,
          formatter: (row) => formatTime(row.createdAt)
        }
      ]
    }
  })

  /** 加载统计 */
  const loadStats = async () => {
    try {
      stats.value = await fetchAccessStats()
    } catch {
      /* 拦截器已提示 */
    }
  }

  const handleSearch = () => {
    replaceSearchParams({ ...searchForm.value })
    getData()
    loadStats()
  }

  const handleReset = () => {
    searchForm.value = { type: '', personName: '', cardNo: '' }
    resetSearchParams()
    getData()
  }

  onMounted(loadStats)
</script>

<style scoped lang="scss">
  .access-page {
    .stats-card {
      margin-bottom: 16px;

      .stats-row {
        display: flex;
        flex-wrap: wrap;
        gap: 28px;

        .stat-item {
          display: flex;
          flex-direction: column;
          gap: 4px;

          .stat-label {
            font-size: 13px;
            color: var(--art-gray-600, #8a8a8a);
          }

          .stat-value {
            font-size: 22px;
            font-weight: 700;
          }
        }
      }
    }
  }
</style>
