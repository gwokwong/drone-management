<!-- RFID 借还管理 - 扫描借出 / 扫描归还 / 借还记录 -->
<template>
  <div class="borrow-page art-full-height">
    <!-- RFID 扫描借还 -->
    <ElCard shadow="never" class="scan-card">
      <template #header>
        <div class="card-header">RFID 电子标签借还</div>
      </template>
      <div class="scan-body">
        <ElInput
          v-model="scanCode"
          size="large"
          placeholder="请扫描或输入无人机 RFID 标签编码"
          clearable
          class="scan-input"
        />
        <div class="scan-actions">
          <ElButton type="primary" :disabled="!scanCode" @click="openBorrow">扫描借出</ElButton>
          <ElButton type="success" :disabled="!scanCode" @click="openReturn">扫描归还</ElButton>
        </div>
      </div>
      <ElAlert
        type="info"
        :closable="false"
        show-icon
        class="scan-tip"
        title="借出后自动记录借用人、借出时间、预计归还时间与全套配件；归还后自动打开对应区域、更新库存并记录实际归还时间与设备状态。"
      />
    </ElCard>

    <!-- 借还记录 -->
    <ElCard class="art-table-card">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElForm :inline="true" :model="searchForm" @submit.prevent>
            <ElFormItem label="标签编码">
              <ElInput v-model="searchForm.droneCode" clearable style="width: 150px" />
            </ElFormItem>
            <ElFormItem label="借用人">
              <ElInput v-model="searchForm.borrowerName" clearable style="width: 130px" />
            </ElFormItem>
            <ElFormItem label="状态">
              <ElSelect v-model="searchForm.status" clearable placeholder="全部" style="width: 120px">
                <ElOption v-for="i in BORROW_STATUS" :key="i.value" :label="i.label" :value="i.value" />
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

    <!-- 借出弹窗 -->
    <ElDialog v-model="borrowVisible" title="扫描借出" width="540">
      <ElForm :model="borrowForm" label-width="120px">
        <ElFormItem label="标签编码">
          <ElInput v-model="borrowForm.code" disabled />
        </ElFormItem>
        <ElFormItem label="借用人">
          <ElInput v-model="borrowForm.borrowerName" placeholder="请输入借用人姓名" />
        </ElFormItem>
        <ElFormItem label="预计归还(小时)">
          <ElInputNumber v-model="borrowForm.expectedReturnHrs" :min="1" />
        </ElFormItem>
        <ElFormItem label="配套配件">
          <ElSelect
            v-model="borrowForm.accessories"
            multiple
            filterable
            allow-create
            placeholder="选择或输入随行配件"
            style="width: 100%"
          >
            <ElOption v-for="a in ACCESSORY_OPTIONS" :key="a" :label="a" :value="a" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="备注">
          <ElInput v-model="borrowForm.remark" type="textarea" :rows="2" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="borrowVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitting" @click="submitBorrow">确认借出</ElButton>
      </template>
    </ElDialog>

    <!-- 归还弹窗 -->
    <ElDialog v-model="returnVisible" title="扫描归还" width="540">
      <ElForm :model="returnForm" label-width="120px">
        <ElFormItem label="标签编码">
          <ElInput v-model="returnForm.code" disabled />
        </ElFormItem>
        <ElFormItem label="归还设备状态">
          <ElSelect v-model="returnForm.returnCondition" style="width: 100%">
            <ElOption v-for="i in DEVICE_STATUS" :key="i.value" :label="i.label" :value="i.value" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="备注">
          <ElInput v-model="returnForm.remark" type="textarea" :rows="2" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="returnVisible = false">取消</ElButton>
        <ElButton type="success" :loading="submitting" @click="submitReturn">确认归还</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { useTable } from '@/hooks/core/useTable'
  import { fetchBorrowRecords, scanBorrow, scanReturn } from '@/api/drone'
  import type { BorrowRecord } from '@/api/drone'
  import { BORROW_STATUS, DEVICE_STATUS, labelOf, typeOf } from '@/utils/drone/dict'
  import { ElTag, ElMessage } from 'element-plus'

  defineOptions({ name: 'DroneBorrow' })

  /** 常用配件 */
  const ACCESSORY_OPTIONS = ['遥控器', '电池', '充电器', '螺旋桨', '存储卡', '云台', '数据线', '运输箱']

  const scanCode = ref('')
  const submitting = ref(false)
  const borrowVisible = ref(false)
  const returnVisible = ref(false)

  const searchForm = ref({ droneCode: '', borrowerName: '', status: '' })

  const borrowForm = ref({
    code: '',
    borrowerName: '',
    expectedReturnHrs: 24,
    accessories: [] as string[],
    remark: ''
  })
  const returnForm = ref({ code: '', returnCondition: 'intact', remark: '' })

  /** 打开借出弹窗 */
  const openBorrow = () => {
    borrowForm.value = {
      code: scanCode.value,
      borrowerName: '',
      expectedReturnHrs: 24,
      accessories: [],
      remark: ''
    }
    borrowVisible.value = true
  }

  /** 打开归还弹窗 */
  const openReturn = () => {
    returnForm.value = { code: scanCode.value, returnCondition: 'intact', remark: '' }
    returnVisible.value = true
  }

  /** 提交借出 */
  const submitBorrow = async () => {
    if (!borrowForm.value.borrowerName) {
      ElMessage.warning('请输入借用人')
      return
    }
    submitting.value = true
    try {
      await scanBorrow({ ...borrowForm.value })
      ElMessage.success('借出成功，已自动记录借用信息')
      borrowVisible.value = false
      scanCode.value = ''
      refreshData()
    } finally {
      submitting.value = false
    }
  }

  /** 提交归还 */
  const submitReturn = async () => {
    submitting.value = true
    try {
      await scanReturn({ ...returnForm.value })
      ElMessage.success('归还成功，库存与库位已更新')
      returnVisible.value = false
      scanCode.value = ''
      refreshData()
    } finally {
      submitting.value = false
    }
  }

  /** 时间格式化 */
  const formatTime = (t?: string) => (t ? t.replace('T', ' ').slice(0, 19) : '-')

  /** 时长格式化（分钟） */
  const formatDuration = (min?: number) => {
    if (min === undefined || min === null) return '-'
    if (min < 60) return `${min} 分钟`
    return `${Math.floor(min / 60)} 小时 ${min % 60} 分`
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
      apiFn: fetchBorrowRecords,
      apiParams: { current: 1, size: 20 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'droneCode', label: '标签编码', width: 130 },
        { prop: 'droneName', label: '无人机', minWidth: 130 },
        { prop: 'borrowerName', label: '借用人', width: 110 },
        {
          prop: 'borrowTime',
          label: '借出时间',
          width: 165,
          formatter: (row) => formatTime(row.borrowTime)
        },
        {
          prop: 'expectedReturnTime',
          label: '预计归还',
          width: 165,
          formatter: (row) => formatTime(row.expectedReturnTime)
        },
        {
          prop: 'actualReturnTime',
          label: '实际归还',
          width: 165,
          formatter: (row) => formatTime(row.actualReturnTime)
        },
        {
          prop: 'status',
          label: '状态',
          width: 100,
          formatter: (row) =>
            h(
              ElTag,
              { type: typeOf(BORROW_STATUS, row.status), size: 'small' },
              () => labelOf(BORROW_STATUS, row.status)
            )
        },
        {
          prop: 'durationMin',
          label: '使用时长',
          width: 120,
          formatter: (row) => formatDuration(row.durationMin)
        },
        { prop: 'positionCode', label: '库位', width: 110 }
      ]
    }
  })

  /** 查询 */
  const handleSearch = () => {
    replaceSearchParams({ ...searchForm.value })
    getData()
  }

  /** 重置 */
  const handleReset = () => {
    searchForm.value = { droneCode: '', borrowerName: '', status: '' }
    resetSearchParams()
    getData()
  }
</script>

<style scoped lang="scss">
  .borrow-page {
    .scan-card {
      margin-bottom: 16px;

      .card-header {
        font-weight: 600;
      }

      .scan-body {
        display: flex;
        gap: 12px;
        align-items: center;

        .scan-input {
          max-width: 420px;
        }

        .scan-actions {
          display: flex;
          gap: 8px;
        }
      }

      .scan-tip {
        margin-top: 12px;
      }
    }
  }
</style>
