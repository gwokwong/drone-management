<!-- 储存室 3D 导航图 - 库位图形化盘点 -->
<template>
  <div class="storage-page art-full-height">
    <!-- 储存室选择 -->
    <ElCard shadow="never" class="filter-card">
      <ElForm :inline="true" @submit.prevent>
        <ElFormItem label="储存室">
          <ElSelect
            v-model="roomId"
            placeholder="请选择储存室"
            style="width: 260px"
            @change="loadMap"
          >
            <ElOption v-for="r in rooms" :key="r.id" :label="`${r.name}（${r.code}）`" :value="r.id" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem>
          <ElButton :loading="loading" @click="loadMap">刷新</ElButton>
        </ElFormItem>
        <ElFormItem v-if="map">
          <ElTag type="info">库位占用 {{ usedCount }} / {{ totalCount }}（{{ occupancyRate }}%）</ElTag>
        </ElFormItem>
      </ElForm>
    </ElCard>

    <!-- 储存室环境信息 -->
    <ElCard v-if="map" shadow="never" class="room-card">
      <ElDescriptions :column="4" border size="small">
        <ElDescriptionsItem label="储存室">{{ map.room.name }}</ElDescriptionsItem>
        <ElDescriptionsItem label="位置">{{ map.room.location || '-' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="温度">{{ map.room.envTemp }} ℃</ElDescriptionsItem>
        <ElDescriptionsItem label="湿度">{{ map.room.envHumidity }} %</ElDescriptionsItem>
      </ElDescriptions>
    </ElCard>

    <!-- 密集架库位图 -->
    <div v-if="map" class="rack-list">
      <ElCard v-for="rack in map.racks" :key="rack.id" shadow="never" class="rack-card">
        <template #header>
          <div class="rack-header">
            <span class="rack-name">{{ rack.name }}（{{ rack.code }}）</span>
            <span class="rack-meta">{{ rack.rows }} 行 × {{ rack.cols }} 列 × {{ rack.layers }} 层</span>
          </div>
        </template>
        <div class="pos-grid">
          <div
            v-for="p in rack.positions"
            :key="p.code"
            class="pos-cell"
            :class="p.occupied ? 'is-occupied' : 'is-empty'"
            @click="openPosition(p)"
          >
            <div class="pos-code">{{ p.code }}</div>
            <div class="pos-drone">{{ p.occupied ? p.droneCode : '空置' }}</div>
          </div>
        </div>
      </ElCard>
    </div>

    <ElCard v-else shadow="never">
      <ElEmpty description="请选择储存室以查看 3D 导航图" />
    </ElCard>

    <!-- 库位详情 -->
    <ElDialog v-model="posVisible" title="库位详情" width="440">
      <ElDescriptions v-if="currentPos" :column="1" border size="small">
        <ElDescriptionsItem label="库位编号">{{ currentPos.code }}</ElDescriptionsItem>
        <ElDescriptionsItem label="行 / 列 / 层">
          {{ currentPos.row }} / {{ currentPos.col }} / {{ currentPos.layer }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="状态">
          <ElTag :type="currentPos.occupied ? 'success' : 'info'" size="small">
            {{ currentPos.occupied ? '已占用' : '空置' }}
          </ElTag>
        </ElDescriptionsItem>
        <ElDescriptionsItem label="无人机标签">{{ currentPos.droneCode || '-' }}</ElDescriptionsItem>
      </ElDescriptions>
      <template #footer>
        <ElButton @click="posVisible = false">关闭</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { fetchRoomList, fetchRoomMap } from '@/api/drone'
  import type { StorageRoom, RoomMap, RackPositionView } from '@/api/drone'

  defineOptions({ name: 'StorageMap' })

  const rooms = ref<StorageRoom[]>([])
  const roomId = ref<number>()
  const map = ref<RoomMap>()
  const loading = ref(false)

  const posVisible = ref(false)
  const currentPos = ref<RackPositionView>()

  /** 库位总数 */
  const totalCount = computed(
    () => map.value?.racks.reduce((sum, r) => sum + r.positions.length, 0) ?? 0
  )
  /** 已占用库位数 */
  const usedCount = computed(
    () =>
      map.value?.racks.reduce((sum, r) => sum + r.positions.filter((p) => p.occupied).length, 0) ?? 0
  )
  /** 占用率 */
  const occupancyRate = computed(() =>
    totalCount.value ? Math.round((usedCount.value / totalCount.value) * 100) : 0
  )

  /** 加载储存室列表 */
  const loadRooms = async () => {
    try {
      const res = await fetchRoomList({ current: 1, size: 100 })
      rooms.value = res.records ?? []
      if (rooms.value.length && !roomId.value) {
        roomId.value = rooms.value[0].id
        loadMap()
      }
    } catch {
      /* 拦截器已提示 */
    }
  }

  /** 加载 3D 导航图 */
  const loadMap = async () => {
    if (!roomId.value) return
    loading.value = true
    try {
      map.value = await fetchRoomMap(roomId.value)
    } catch {
      /* 拦截器已提示 */
    } finally {
      loading.value = false
    }
  }

  /** 查看库位详情 */
  const openPosition = (p: RackPositionView) => {
    currentPos.value = p
    posVisible.value = true
  }

  onMounted(loadRooms)
</script>

<style scoped lang="scss">
  .storage-page {
    .filter-card,
    .room-card {
      margin-bottom: 16px;
    }

    .rack-list {
      .rack-card {
        margin-bottom: 16px;

        .rack-header {
          display: flex;
          align-items: center;
          justify-content: space-between;

          .rack-name {
            font-weight: 600;
          }

          .rack-meta {
            font-size: 12px;
            color: var(--art-gray-600, #8a8a8a);
          }
        }

        .pos-grid {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(96px, 1fr));
          gap: 10px;

          .pos-cell {
            padding: 8px 6px;
            text-align: center;
            cursor: pointer;
            border: 1px solid var(--art-border-color, #e4e7ed);
            border-radius: 6px;
            transition: all 0.2s;

            &:hover {
              transform: translateY(-2px);
              box-shadow: 0 4px 10px rgb(0 0 0 / 8%);
            }

            .pos-code {
              font-size: 12px;
              font-weight: 600;
            }

            .pos-drone {
              margin-top: 2px;
              font-size: 11px;
              opacity: 0.75;
            }

            &.is-occupied {
              color: #fff;
              background-color: #60c041;
              border-color: #60c041;
            }

            &.is-empty {
              color: var(--art-gray-700, #5a5a5a);
              background-color: var(--art-bg-color, #f5f7fa);
            }
          }
        }
      }
    }
  }
</style>
