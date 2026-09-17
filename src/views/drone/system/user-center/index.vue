<template>
  <div class="user-center-page art-full-height">
    <el-row :gutter="16">
      <el-col :xs="24" :md="8">
        <el-card shadow="never" class="user-card">
          <div class="user-card-header">
            <el-avatar :size="72" :src="userInfo.avatar || defaultAvatar" />
            <div class="user-card-meta">
              <h3>{{ userInfo.nickname || userInfo.userName || '-' }}</h3>
              <p>{{ userInfo.userName }}</p>
              <div class="user-card-roles">
                <el-tag v-for="r in roles" :key="r" type="primary" size="small" effect="light">
                  {{ roleName(r) }}
                </el-tag>
                <span v-if="!roles.length" class="text-placeholder">未分配角色</span>
              </div>
            </div>
          </div>
          <el-divider />
          <el-descriptions :column="1" size="small">
            <el-descriptions-item label="用户 ID">{{ userInfo.userId ?? '-' }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ userInfo.email || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="16">
        <el-card shadow="never">
          <template #header>
            <span>修改密码</span>
          </template>
          <el-form
            ref="formRef"
            :model="form"
            :rules="rules"
            label-width="100px"
            style="max-width: 460px"
          >
            <el-form-item label="原密码" prop="oldPassword">
              <el-input v-model="form.oldPassword" type="password" show-password />
            </el-form-item>
            <el-form-item label="新密码" prop="newPassword">
              <el-input v-model="form.newPassword" type="password" show-password />
            </el-form-item>
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input v-model="form.confirmPassword" type="password" show-password />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="submitting" @click="submit">保存</el-button>
              <el-button @click="reset">重置</el-button>
            </el-form-item>
          </el-form>
          <el-alert
            type="info"
            :closable="false"
            show-icon
            title="修改密码需后端提供对应接口，当前页面已按接口规范预留，接入后即可生效。"
          />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
  import { computed, reactive, ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import type { FormInstance, FormRules } from 'element-plus'
  import { useUserStore } from '@/store/modules/user'
  import { ROLE_LABEL } from '@/utils/drone/dict'

  defineOptions({ name: 'DroneUserCenter' })

  const userStore = useUserStore()
  const userInfo = computed(() => userStore.getUserInfo)
  const roles = computed(() => userInfo.value.roles ?? [])

  const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

  const roleName = (code: string) => ROLE_LABEL[code] ?? code

  const formRef = ref<FormInstance>()
  const submitting = ref(false)
  const form = reactive({
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  })

  const rules: FormRules = {
    oldPassword: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
    newPassword: [
      { required: true, message: '请输入新密码', trigger: 'blur' },
      { min: 6, message: '密码长度不少于 6 位', trigger: 'blur' }
    ],
    confirmPassword: [
      { required: true, message: '请再次输入新密码', trigger: 'blur' },
      {
        validator: (_r, v, cb) => (v === form.newPassword ? cb() : cb(new Error('两次输入不一致'))),
        trigger: 'blur'
      }
    ]
  }

  const submit = async () => {
    if (!formRef.value) return
    await formRef.value.validate(async (ok) => {
      if (!ok) return
      submitting.value = true
      try {
        // 后端暂未提供修改密码接口，此处保留调用位；接口上线后替换为真实请求即可。
        ElMessage.info('后端暂未开放修改密码接口，请先在系统管理中由管理员重置')
      } finally {
        submitting.value = false
      }
    })
  }

  const reset = () => {
    form.oldPassword = ''
    form.newPassword = ''
    form.confirmPassword = ''
    formRef.value?.clearValidate()
  }
</script>

<style scoped lang="scss">
  .user-center-page {
    .user-card-header {
      display: flex;
      gap: 16px;
      align-items: center;
    }

    .user-card-meta {
      h3 {
        margin: 0 0 4px;
        font-size: 18px;
        font-weight: 600;
      }

      p {
        margin: 0 0 8px;
        font-size: 13px;
        color: var(--art-text-gray-600);
      }
    }

    .user-card-roles {
      display: flex;
      flex-wrap: wrap;
      gap: 6px;
    }

    .text-placeholder {
      font-size: 12px;
      color: var(--art-text-gray-500);
    }
  }
</style>
