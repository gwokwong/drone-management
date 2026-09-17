<!-- 报表管理 - 日/周/月/年统计报表，图形化展示并支持导出 Excel -->
<template>
  <div class="report-page art-full-height">
    <!-- 报表类型 -->
    <ElCard shadow="never" class="filter-card">
      <ElForm :inline="true" @submit.prevent>
        <ElFormItem label="报表类型">
          <ElRadioGroup v-model="reportType" @change="loadReport">
            <ElRadioButton v-for="i in REPORT_TYPE" :key="i.value" :value="i.value">
              {{ i.label }}
            </ElRadioButton>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem>
          <ElButton type="primary" :loading="loading" @click="loadReport">生成报表</ElButton>
          <ElButton :loading="exporting" @click="handleExport">导出 Excel</ElButton>
        </ElFormItem>
        <ElFormItem v-if="report">
          <span class="range-tip">
            统计周期：{{ formatTime(report.start) }} ~ {{ formatTime(report.end) }}
          </span>
        </ElFormItem>
      </ElForm>
    </ElCard>

    <!-- 指标卡片 -->
    <ElRow :gutter="16">
      <ElCol v-for="s in statItems" :key="s.label" :xs="12" :sm="8" :md="8" :lg="4">
        <ElCard shadow="never" class="stat-card">
          <div class="stat-label">{{ s.label }}</div>
          <div class="stat-value" :style="{ color: s.color }">{{ s.value }}</div>
        </ElCard>
      </ElCol>
    </ElRow>

    <!-- 趋势图 -->
    <ElCard shadow="never" class="chart-card">
      <template #header>
        <div class="card-header">周期内借还趋势</div>
      </template>
      <div v-if="trend.length" class="chart">
        <div v-for="t in trend" :key="t.date" class="chart-col">
          <div class="chart-value">{{ t.borrow }}</div>
          <div class="bar-wrap">
            <div class="bar" :style="{ height: barHeight(t.borrow) }"></div>
          </div>
          <div class="chart-label">{{ t.date }}</div>
        </div>
      </div>
      <ElEmpty v-else description="暂无趋势数据" :image-size="60" />
    </ElCard>

    <!-- 趋势明细 -->
    <ElCard shadow="never" class="trend-card">
      <template #header>
        <div class="card-header">趋势明细</div>
      </template>
      <ElTable :data="trend" size="small" max-height="300">
        <ElTableColumn prop="date" label="日期" min-width="120" />
        <ElTableColumn prop="borrow" label="借还次数" min-width="120" />
        <ElTableColumn prop="alarm" label="告警数" min-width="110" />
        <ElTableColumn prop="access" label="进出人次" min-width="120" />
        <template #empty>
          <ElEmpty description="暂无数据" :image-size="60" />
        </template>
      </ElTable>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { fetchReport, exportReport } from '@/api/drone'
  import type { ReportData } from '@/api/drone'
  import { REPORT_TYPE } from '@/utils/drone/dict'
  import { ElMessage } from 'element-plus'

  defineOptions({ name: 'ReportIndex' })

  const reportType = ref('daily')
  const report = ref<ReportData>()
  const loading = ref(false)
  const exporting = ref(false)

  /** 趋势数据 */
  const trend = computed(() => report.value?.trend ?? [])

  /** 指标卡片 */
  const statItems = computed(() => [
    { label: '无人机总数', value: report.value?.droneTotal ?? 0, color: '#1D84FF' },
    { label: '当前借出', value: report.value?.droneBorrowed ?? 0, color: '#F9901F' },
    { label: '周期借还次数', value: report.value?.borrowCount ?? 0, color: '#60C041' },
    { label: '逾期未还', value: report.value?.overdueCount ?? 0, color: '#F56C6C' },
    { label: '周期告警数', value: report.value?.alarmCount ?? 0, color: '#F9901F' },
    { label: '周期进出人次', value: report.value?.accessCount ?? 0, color: '#38C0FC' }
  ])

  /** 柱状图最大高度基准 */
  const maxBorrow = computed(() => Math.max(1, ...trend.value.map((t) => t.borrow)))
  const barHeight = (v: number) => `${Math.round((v / maxBorrow.value) * 100)}%`

  const formatTime = (t?: string) => (t ? t.replace('T', ' ').slice(0, 19) : '-')

  /** 生成报表 */
  const loadReport = async () => {
    loading.value = true
    try {
      report.value = await fetchReport({ type: reportType.value })
    } catch {
      /* 拦截器已提示 */
    } finally {
      loading.value = false
    }
  }

  /** 导出 Excel */
  const handleExport = async () => {
    exporting.value = true
    try {
      const label = REPORT_TYPE.find((i) => i.value === reportType.value)?.label ?? '报表'
      await exportReport(reportType.value, `${label}_${reportType.value}.xlsx`)
      ElMessage.success('报表已导出')
    } catch {
      ElMessage.error('导出失败')
    } finally {
      exporting.value = false
    }
  }

  onMounted(loadReport)
</script>

<style scoped lang="scss">
  .report-page {
    .filter-card {
      margin-bottom: 16px;

      .range-tip {
        font-size: 12px;
        color: var(--art-gray-600, #8a8a8a);
      }
    }

    .stat-card {
      margin-bottom: 16px;

      .stat-label {
        font-size: 13px;
        color: var(--art-gray-600, #8a8a8a);
      }

      .stat-value {
        margin-top: 6px;
        font-size: 24px;
        font-weight: 700;
      }
    }

    .chart-card,
    .trend-card {
      margin-bottom: 16px;
    }

    .card-header {
      font-weight: 600;
    }

    .chart {
      display: flex;
      gap: 12px;
      align-items: flex-end;
      height: 220px;
      padding-top: 10px;
      overflow-x: auto;

      .chart-col {
        display: flex;
        flex: 1;
        flex-direction: column;
        align-items: center;
        min-width: 46px;
        height: 100%;

        .chart-value {
          margin-bottom: 4px;
          font-size: 12px;
          font-weight: 600;
        }

        .bar-wrap {
          display: flex;
          flex: 1;
          align-items: flex-end;
          width: 100%;

          .bar {
            width: 100%;
            min-height: 2px;
            background: linear-gradient(180deg, #5d87ff 0%, #1d84ff 100%);
            border-radius: 4px 4px 0 0;
            transition: height 0.3s;
          }
        }

        .chart-label {
          margin-top: 6px;
          font-size: 11px;
          color: var(--art-gray-600, #8a8a8a);
          white-space: nowrap;
        }
      }
    }
  }
</style>
