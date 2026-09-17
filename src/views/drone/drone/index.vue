<!-- 无人机台账 - 无序存放 / 有序管理 -->
<template>
  <div class="drone-page art-full-height">
    <!-- 实时统计 -->
    <ElCard shadow="never" class="stats-card">
      <div class="stats-row">
        <div v-for="s in statItems" :key="s.label" class="stat-item">
          <span class="stat-label">{{ s.label }}</span>
          <span class="stat-value" :style="{ color: s.color }">{{ s.value }}</span>
        </div>
      </div>
    </ElCard>

    <!-- 搜索栏 -->
    <ElCard shadow="never" class="search-card">
      <ElForm :inline="true" :model="searchForm" @submit.prevent>
        <ElFormItem label="关键字">
          <ElInput
            v-model="searchForm.keyword"
            placeholder="名称 / 编号 / 型号"
            clearable
            style="width: 180px"
          />
        </ElFormItem>
        <ElFormItem label="存放状态">
          <ElSelect v-model="searchForm.storageStatus" clearable placeholder="全部" style="width: 130px">
            <ElOption v-for="i in STORAGE_STATUS" :key="i.value" :label="i.label" :value="i.value" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="设备状态">
          <ElSelect v-model="searchForm.deviceStatus" clearable placeholder="全部" style="width: 140px">
            <ElOption v-for="i in DEVICE_STATUS" :key="i.value" :label="i.label" :value="i.value" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem>
          <ElButton type="primary" @click="handleSearch">查询</ElButton>
          <ElButton @click="handleReset">重置</ElButton>
        </ElFormItem>
      </ElForm>
    </ElCard>

    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElSpace wrap>
            <ElButton @click="showDialog('add')" v-ripple>新增无人机</ElButton>
          </ElSpace>
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
      :title="dialogType === 'add' ? '新增无人机' : '编辑无人机'"
      width="700"
    >
      <ElForm ref="formRef" :model="form" :rules="rules" label-width="100px">
        <ElRow :gutter="16">
          <ElCol :span="12">
            <ElFormItem label="RFID编码" prop="code">
              <ElInput v-model="form.code" placeholder="电子标签编码" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="名称" prop="name">
              <ElInput v-model="form.name" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="型号">
              <ElInput v-model="form.model" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="分类">
              <ElInput v-model="form.category" placeholder="如 多旋翼 / 固定翼" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="存放状态">
              <ElSelect v-model="form.storageStatus" style="width: 100%">
                <ElOption v-for="i in STORAGE_STATUS" :key="i.value" :label="i.label" :value="i.value" />
              </ElSelect>
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="设备状态">
              <ElSelect v-model="form.deviceStatus" style="width: 100%">
                <ElOption v-for="i in DEVICE_STATUS" :key="i.value" :label="i.label" :value="i.value" />
              </ElSelect>
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="电量(%)">
              <ElInputNumber v-model="form.battery" :min="0" :max="100" style="width: 100%" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="库位编号">
              <ElInput v-model="form.positionCode" placeholder="如 A-01-1" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="所属储存室ID">
              <ElInputNumber v-model="form.roomId" :min="0" style="width: 100%" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="12">
            <ElFormItem label="制造商">
              <ElInput v-model="form.manufacturer" />
            </ElFormItem>
          </ElCol>
          <ElCol :span="24">
            <ElFormItem label="备注">
              <ElInput v-model="form.remark" type="textarea" :rows="2" />
            </ElFormItem>
          </ElCol>
        </ElRow>
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
  import { fetchDroneList, fetchDroneStats, createDrone, updateDrone, deleteDrone } from '@/api/drone'
  import type { Drone, DroneStats } from '@/api/drone'
  import { STORAGE_STATUS, DEVICE_STATUS, labelOf, typeOf } from '@/utils/drone/dict'
  import { ElTag, ElMessage, ElMessageBox } from 'element-plus'
  import { DialogType } from '@/types'

  defineOptions({ name: 'DroneList' })

  const searchForm = ref({ keyword: '', storageStatus: '', deviceStatus: '' })

  /** 实时统计 */
  const stats = ref<DroneStats>({
    total: 0,
    inPosition: 0,
    borrowed: 0,
    intact: 0,
    damaged: 0,
    scrapped: 0
  })
  const statItems = computed(() => [
    { label: '总数', value: stats.value.total, color: '#1D84FF' },
    { label: '在位', value: stats.value.inPosition, color: '#60C041' },
    { label: '借出', value: stats.value.borrowed, color: '#F9901F' },
    { label: '完好', value: stats.value.intact, color: '#60C041' },
    { label: '损坏待修', value: stats.value.damaged, color: '#F9901F' },
    { label: '已报废', value: stats.value.scrapped, color: '#F56C6C' }
  ])

  const dialogType = ref<DialogType>('add')
  const dialogVisible = ref(false)
  const submitting = ref(false)
  const formRef = ref()
  const form = ref<Partial<Drone>>({})

  const rules = {
    code: [{ required: true, message: '请输入 RFID 编码', trigger: 'blur' }],
    name: [{ required: true, message: '请输入无人机名称', trigger: 'blur' }]
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
      apiFn: fetchDroneList,
      apiParams: { current: 1, size: 20 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'code', label: 'RFID编码', width: 140 },
        { prop: 'name', label: '名称', minWidth: 140 },
        { prop: 'model', label: '型号', width: 120 },
        { prop: 'category', label: '分类', width: 110 },
        {
          prop: 'storageStatus',
          label: '存放状态',
          width: 100,
          formatter: (row) =>
            h(
              ElTag,
              { type: typeOf(STORAGE_STATUS, row.storageStatus), size: 'small' },
              () => labelOf(STORAGE_STATUS, row.storageStatus)
            )
        },
        {
          prop: 'deviceStatus',
          label: '设备状态',
          width: 120,
          formatter: (row) =>
            h(
              ElTag,
              { type: typeOf(DEVICE_STATUS, row.deviceStatus), size: 'small' },
              () => labelOf(DEVICE_STATUS, row.deviceStatus)
            )
        },
        { prop: 'battery', label: '电量', width: 90, formatter: (row) => `${row.battery}%` },
        { prop: 'positionCode', label: '库位', width: 110 },
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
  const showDialog = (type: DialogType, row?: Drone) => {
    dialogType.value = type
    form.value = row
      ? { ...row }
      : { storageStatus: 'in_position', deviceStatus: 'intact', battery: 100 }
    nextTick(() => {
      dialogVisible.value = true
    })
  }

  /** 提交新增 / 编辑 */
  const handleSubmit = async () => {
    await formRef.value?.validate(async (valid: boolean) => {
      if (!valid) return
      submitting.value = true
      try {
        if (dialogType.value === 'add') {
          await createDrone(form.value)
          ElMessage.success('新增成功')
        } else {
          await updateDrone(form.value.id!, form.value)
          ElMessage.success('更新成功')
        }
        dialogVisible.value = false
        refreshData()
        loadStats()
      } finally {
        submitting.value = false
      }
    })
  }

  /** 删除 */
  const handleDelete = (row: Drone) => {
    ElMessageBox.confirm(`确定要删除无人机「${row.name}」吗？`, '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'error'
    }).then(async () => {
      await deleteDrone(row.id!)
      ElMessage.success('删除成功')
      refreshData()
      loadStats()
    })
  }

  /** 查询 */
  const handleSearch = () => {
    replaceSearchParams({ ...searchForm.value })
    getData()
  }

  /** 重置 */
  const handleReset = () => {
    searchForm.value = { keyword: '', storageStatus: '', deviceStatus: '' }
    resetSearchParams()
    getData()
  }

  /** 加载统计 */
  const loadStats = async () => {
    try {
      stats.value = await fetchDroneStats()
    } catch {
      /* 拦截器已提示 */
    }
  }

  onMounted(loadStats)
</script>

<style scoped lang="scss">
  .drone-page {
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

    .search-card {
      margin-bottom: 16px;
    }
  }
</style>
