<!-- 视频监控 - 摄像头集中管理 + 密集架开架录像 -->
<template>
  <div class="video-page art-full-height">
    <!-- 摄像头管理 -->
    <ElCard shadow="never" class="panel">
      <template #header>
        <div class="card-header">
          <span>摄像头集中管理</span>
          <ElButton size="small" :loading="vLoading" @click="vRefresh">刷新</ElButton>
        </div>
      </template>
      <ArtTable
        :loading="vLoading"
        :data="vData"
        :columns="vColumns"
        :pagination="vPagination"
        @pagination:size-change="vHandleSizeChange"
        @pagination:current-change="vHandleCurrentChange"
      />
    </ElCard>

    <!-- 开架录像记录 -->
    <ElCard shadow="never" class="panel">
      <template #header>
        <div class="card-header">
          <span>密集架开架录像记录</span>
          <ElButton size="small" :loading="rLoading" @click="rRefresh">刷新</ElButton>
        </div>
      </template>
      <ArtTable
        :loading="rLoading"
        :data="rData"
        :columns="rColumns"
        :pagination="rPagination"
        @pagination:size-change="rHandleSizeChange"
        @pagination:current-change="rHandleCurrentChange"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { useTable } from '@/hooks/core/useTable'
  import { fetchVideoList, fetchOpenRackList, openRack, stopOpenRack } from '@/api/drone'
  import type { VideoChannel, OpenRackRecord } from '@/api/drone'
  import { VIDEO_STATUS, OPEN_RACK_STATUS, labelOf, typeOf } from '@/utils/drone/dict'
  import { ElTag, ElMessage, ElButton } from 'element-plus'

  defineOptions({ name: 'MonitorVideo' })

  const formatTime = (t?: string) => (t ? t.replace('T', ' ').slice(0, 19) : '-')

  /** 开架录像 */
  const handleOpenRack = async (row: VideoChannel) => {
    await openRack(row.id!)
    ElMessage.success(`已启动「${row.name}」开架录像`)
    rRefresh()
  }

  /** 停止录像 */
  const handleStop = async (row: OpenRackRecord) => {
    await stopOpenRack(row.id!)
    ElMessage.success('录像已结束')
    rRefresh()
  }

  /** 摄像头列表 */
  const {
    columns: vColumns,
    data: vData,
    loading: vLoading,
    pagination: vPagination,
    handleSizeChange: vHandleSizeChange,
    handleCurrentChange: vHandleCurrentChange,
    refreshData: vRefresh
  } = useTable({
    core: {
      apiFn: fetchVideoList,
      apiParams: { current: 1, size: 10 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'name', label: '摄像头名称', minWidth: 160 },
        { prop: 'roomId', label: '所属储存室', width: 120 },
        { prop: 'rackId', label: '关联密集架', width: 120 },
        {
          prop: 'status',
          label: '状态',
          width: 100,
          formatter: (row) =>
            h(ElTag, { type: typeOf(VIDEO_STATUS, row.status), size: 'small' }, () =>
              labelOf(VIDEO_STATUS, row.status)
            )
        },
        { prop: 'url', label: '视频流地址', minWidth: 220, showOverflowTooltip: true },
        {
          prop: 'operation',
          label: '操作',
          width: 130,
          fixed: 'right',
          formatter: (row) =>
            h(
              ElButton,
              { link: true, type: 'primary', size: 'small', onClick: () => handleOpenRack(row) },
              () => '开架录像'
            )
        }
      ]
    }
  })

  /** 开架录像记录 */
  const {
    columns: rColumns,
    data: rData,
    loading: rLoading,
    pagination: rPagination,
    handleSizeChange: rHandleSizeChange,
    handleCurrentChange: rHandleCurrentChange,
    refreshData: rRefresh
  } = useTable({
    core: {
      apiFn: fetchOpenRackList,
      apiParams: { current: 1, size: 10 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'rackCode', label: '密集架编号', width: 140 },
        { prop: 'operator', label: '操作人', width: 120 },
        {
          prop: 'startTime',
          label: '开始时间',
          width: 170,
          formatter: (row) => formatTime(row.startTime)
        },
        {
          prop: 'endTime',
          label: '结束时间',
          width: 170,
          formatter: (row) => formatTime(row.endTime)
        },
        {
          prop: 'status',
          label: '状态',
          width: 100,
          formatter: (row) =>
            h(ElTag, { type: typeOf(OPEN_RACK_STATUS, row.status), size: 'small' }, () =>
              labelOf(OPEN_RACK_STATUS, row.status)
            )
        },
        { prop: 'videoUrl', label: '回放地址', minWidth: 200, showOverflowTooltip: true },
        {
          prop: 'operation',
          label: '操作',
          width: 100,
          fixed: 'right',
          formatter: (row) =>
            row.status === 'recording'
              ? h(
                  ElButton,
                  { link: true, type: 'danger', size: 'small', onClick: () => handleStop(row) },
                  () => '结束录像'
                )
              : '-'
        }
      ]
    }
  })
</script>

<style scoped lang="scss">
  .video-page {
    .panel {
      margin-bottom: 16px;
    }

    .card-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      font-weight: 600;
    }
  }
</style>
