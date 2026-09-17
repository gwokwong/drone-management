<!-- 数据大屏 - 实时数据总览 -->
<template>
  <div class="dashboard-page art-full-height">
    <!-- 顶部操作条 -->
    <div class="page-toolbar">
      <div class="toolbar-title">实时数据总览</div>
      <div class="toolbar-extra">
        <span class="refresh-tip">最近刷新：{{ lastRefresh || '--' }}</span>
        <ElButton size="small" :loading="loading" @click="loadData">刷新</ElButton>
        <ElSwitch v-model="autoRefresh" size="small" active-text="自动" />
      </div>
    </div>

    <!-- 核心指标卡片 -->
    <ElRow :gutter="16">
      <ElCol v-for="item in statCards" :key="item.label" :xs="12" :sm="8" :md="6" :lg="4">
        <ElCard class="stat-card" shadow="never">
          <div class="stat-label">{{ item.label }}</div>
          <div class="stat-value" :style="{ color: item.color }">{{ item.value }}</div>
        </ElCard>
      </ElCol>
    </ElRow>

    <ElRow :gutter="16" class="section">
      <!-- 储存室库位占用 -->
      <ElCol :xs="24" :lg="10">
        <ElCard shadow="never" class="panel">
          <template #header>
            <div class="panel-header">
              <span>储存室库位占用</span>
              <ElTag size="small" type="info">图形化盘点</ElTag>
            </div>
          </template>
          <div v-if="overview.roomStats.length" class="occ-list">
            <div v-for="r in overview.roomStats" :key="r.roomId" class="occ-item">
              <div class="occ-head">
                <span class="occ-name">{{ r.name }}</span>
                <span class="occ-num">{{ r.used }} / {{ r.total }}</span>
              </div>
              <ElProgress
                :percentage="percent(r.used, r.total)"
                :stroke-width="8"
                :show-text="false"
              />
            </div>
          </div>
          <ElEmpty v-else description="暂无储存室数据" :image-size="60" />
        </ElCard>
      </ElCol>

      <!-- 最新告警 -->
      <ElCol :xs="24" :lg="14">
        <ElCard shadow="never" class="panel">
          <template #header>
            <div class="panel-header">
              <span>最新告警</span>
              <ElTag v-if="overview.alarm.active" size="small" type="danger">
                {{ overview.alarm.active }} 条未处理
              </ElTag>
            </div>
          </template>
          <ElTable :data="overview.recentAlarms" size="small" max-height="260">
            <ElTableColumn prop="title" label="告警标题" min-width="170" show-overflow-tooltip />
            <ElTableColumn prop="type" label="类型" width="110">
              <template #default="{ row }">
                <ElTag size="small" :type="typeOf(ALARM_TYPE, row.type)">
                  {{ labelOf(ALARM_TYPE, row.type) }}
                </ElTag>
              </template>
            </ElTableColumn>
            <ElTableColumn prop="level" label="级别" width="90">
              <template #default="{ row }">
                <ElTag size="small" :type="typeOf(ALARM_LEVEL, row.level)">
                  {{ labelOf(ALARM_LEVEL, row.level) }}
                </ElTag>
              </template>
            </ElTableColumn>
            <ElEmpty slot="empty" description="暂无告警" />
          </ElTable>
        </ElCard>
      </ElCol>
    </ElRow>

    <ElRow :gutter="16" class="section">
      <!-- 最近人员进出 -->
      <ElCol :span="24">
        <ElCard shadow="never" class="panel">
          <template #header>
            <div class="panel-header">
              <span>最近人员进出</span>
              <ElTag size="small" type="success">今日 {{ overview.accessToday }} 人次</ElTag>
            </div>
          </template>
          <ElTable :data="overview.recentAccess" size="small" max-height="240">
            <ElTableColumn prop="personName" label="人员" width="120" />
            <ElTableColumn prop="cardNo" label="卡号" width="140" />
            <ElTableColumn prop="type" label="事件类型" width="140">
              <template #default="{ row }">
                <ElTag size="small" :type="typeOf(ACCESS_TYPE, row.type)">
                  {{ labelOf(ACCESS_TYPE, row.type) }}
                </ElTag>
              </template>
            </ElTableColumn>
            <ElTableColumn prop="createdAt" label="时间" min-width="180">
              <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
            </ElTableColumn>
            <template #empty>
              <ElEmpty description="暂无进出记录" :image-size="60" />
            </template>
          </ElTable>
        </ElCard>
      </ElCol>
    </ElRow>
  </div>
</template>

<script setup lang="ts">
  import { fetchDashboardOverview } from '@/api/drone'
  import type { DashboardOverview } from '@/api/drone'
  import { ALARM_TYPE, ALARM_LEVEL, ACCESS_TYPE, labelOf, typeOf } from '@/utils/drone/dict'

  defineOptions({ name: 'DashboardIndex' })

  const loading = ref(false)
  const autoRefresh = ref(true)
  const lastRefresh = ref('')
  let timer: ReturnType<typeof setInterval> | null = null

  /** 大屏数据 */
  const overview = ref<DashboardOverview>({
    drone: {
      total: 0,
      inPosition: 0,
      borrowed: 0,
      intact: 0,
      damaged: 0,
      scrapped: 0,
      activeBorrow: 0
    },
    alarm: { active: 0, critical: 0 },
    accessToday: 0,
    env: { total: 0, abnormal: 0, offline: 0 },
    recentAlarms: [],
    recentAccess: [],
    roomStats: []
  })

  /** 指标卡片 */
  const statCards = computed(() => {
    const d = overview.value.drone
    return [
      { label: '无人机总数', value: d.total, color: '#1D84FF' },
      { label: '在位', value: d.inPosition, color: '#60C041' },
      { label: '借出', value: d.borrowed, color: '#F9901F' },
      { label: '设备完好', value: d.intact, color: '#60C041' },
      { label: '损坏待修', value: d.damaged, color: '#F9901F' },
      { label: '已报废', value: d.scrapped, color: '#F56C6C' },
      { label: '未处理告警', value: overview.value.alarm.active, color: '#F56C6C' },
      { label: '环境异常', value: overview.value.env.abnormal, color: '#F9901F' }
    ]
  })

  /** 占用百分比 */
  const percent = (used: number, total: number) => (total > 0 ? Math.round((used / total) * 100) : 0)

  /** 时间格式化 */
  const formatTime = (t?: string) => (t ? t.replace('T', ' ').slice(0, 19) : '-')

  /** 加载数据 */
  const loadData = async () => {
    loading.value = true
    try {
      overview.value = await fetchDashboardOverview()
      lastRefresh.value = new Date().toLocaleTimeString('zh-CN', { hour12: false })
    } catch {
      // 错误已由 http 拦截器统一提示
    } finally {
      loading.value = false
    }
  }

  /** 启停自动刷新（30 秒） */
  watch(autoRefresh, (on) => {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
    if (on) timer = setInterval(loadData, 30000)
  })

  onMounted(() => {
    loadData()
    if (autoRefresh.value) timer = setInterval(loadData, 30000)
  })

  onUnmounted(() => {
    if (timer) clearInterval(timer)
  })
</script>

<style scoped lang="scss">
  .dashboard-page {
    padding: 16px;

    .page-toolbar {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 14px;

      .toolbar-title {
        font-size: 16px;
        font-weight: 600;
      }

      .toolbar-extra {
        display: flex;
        gap: 10px;
        align-items: center;

        .refresh-tip {
          font-size: 12px;
          color: var(--art-gray-600, #8a8a8a);
        }
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
        font-size: 26px;
        font-weight: 700;
      }
    }

    .section {
      .panel {
        margin-bottom: 16px;
      }

      .panel-header {
        display: flex;
        gap: 8px;
        align-items: center;
        font-weight: 600;
      }

      .occ-list {
        .occ-item {
          margin-bottom: 12px;

          .occ-head {
            display: flex;
            justify-content: space-between;
            margin-bottom: 4px;
            font-size: 13px;
          }

          .occ-name {
            font-weight: 500;
          }

          .occ-num {
            color: var(--art-gray-600, #8a8a8a);
          }
        }
      }
    }
  }
</style>
