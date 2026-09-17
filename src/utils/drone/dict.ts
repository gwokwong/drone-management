/**
 * 无人机管理平台 - 状态字典
 *
 * 统一管理各业务状态的中文案与标签配色，供列表/表单/大屏复用。
 *
 * @module utils/drone/dict
 */

type TagType = 'success' | 'info' | 'warning' | 'danger' | 'primary'

export interface DictItem {
  label: string
  value: string
  type: TagType
}

/** 无人机存放状态 */
export const STORAGE_STATUS: DictItem[] = [
  { label: '在位', value: 'in_position', type: 'success' },
  { label: '借出', value: 'borrowed', type: 'warning' }
]

/** 无人机设备状态 */
export const DEVICE_STATUS: DictItem[] = [
  { label: '完好', value: 'intact', type: 'success' },
  { label: '损坏待维修', value: 'damaged', type: 'warning' },
  { label: '报废', value: 'scrapped', type: 'danger' }
]

/** 告警类型 */
export const ALARM_TYPE: DictItem[] = [
  { label: '断网异常', value: 'network_disconnect', type: 'danger' },
  { label: '漏水异常', value: 'water_leak', type: 'danger' },
  { label: '门锁异常', value: 'door_alarm', type: 'warning' },
  { label: '设备故障', value: 'device_fault', type: 'warning' },
  { label: '低电量', value: 'low_battery', type: 'warning' },
  { label: '环境异常', value: 'env_abnormal', type: 'info' }
]

/** 告警级别 */
export const ALARM_LEVEL: DictItem[] = [
  { label: '提示', value: 'info', type: 'info' },
  { label: '警告', value: 'warning', type: 'warning' },
  { label: '严重', value: 'critical', type: 'danger' }
]

/** 告警状态 */
export const ALARM_STATUS: DictItem[] = [
  { label: '未处理', value: 'active', type: 'danger' },
  { label: '已处理', value: 'resolved', type: 'success' },
  { label: '已忽略', value: 'ignored', type: 'info' }
]

/** 人员进出类型 */
export const ACCESS_TYPE: DictItem[] = [
  { label: '卡号认证失败', value: 'card_auth_fail', type: 'danger' },
  { label: '门锁打开', value: 'door_open', type: 'success' },
  { label: '门锁关闭', value: 'door_close', type: 'info' },
  { label: '人脸认证通过', value: 'face_auth_pass', type: 'success' }
]

/** 环境设备状态 */
export const ENV_STATUS: DictItem[] = [
  { label: '正常', value: 'normal', type: 'success' },
  { label: '异常', value: 'abnormal', type: 'danger' },
  { label: '离线', value: 'offline', type: 'info' }
]

/** 环境设备类型 */
export const ENV_TYPE: DictItem[] = [
  { label: '温度', value: 'temperature', type: 'warning' },
  { label: '湿度', value: 'humidity', type: 'primary' },
  { label: '烟感', value: 'smoke', type: 'danger' },
  { label: '水浸', value: 'water', type: 'primary' },
  { label: '门禁', value: 'door', type: 'success' },
  { label: '供电', value: 'power', type: 'warning' }
]

/** 借还状态 */
export const BORROW_STATUS: DictItem[] = [
  { label: '借用中', value: 'borrowing', type: 'primary' },
  { label: '已归还', value: 'returned', type: 'success' },
  { label: '已逾期', value: 'overdue', type: 'danger' }
]

/** 视频通道状态 */
export const VIDEO_STATUS: DictItem[] = [
  { label: '在线', value: 'online', type: 'success' },
  { label: '离线', value: 'offline', type: 'info' }
]

/** 开架录像状态 */
export const OPEN_RACK_STATUS: DictItem[] = [
  { label: '录像中', value: 'recording', type: 'danger' },
  { label: '已结束', value: 'finished', type: 'info' }
]

/** 系统日志分类 */
export const LOG_CATEGORY: DictItem[] = [
  { label: '设备', value: 'device', type: 'primary' },
  { label: '用户', value: 'user', type: 'success' },
  { label: '告警', value: 'alarm', type: 'warning' },
  { label: '操作', value: 'operation', type: 'info' }
]

/** 报表类型 */
export const REPORT_TYPE: DictItem[] = [
  { label: '日报', value: 'daily', type: 'primary' },
  { label: '周报', value: 'weekly', type: 'success' },
  { label: '月报', value: 'monthly', type: 'warning' },
  { label: '年报', value: 'yearly', type: 'danger' }
]

/** 角色编码 -> 中文名（与后端 seed 的角色保持一致） */
export const ROLE_LABEL: Record<string, string> = {
  R_SUPER: '超级管理员',
  R_ADMIN: '系统管理员',
  R_OPERATOR: '操作员'
}

/** 取中文文案 */
export function labelOf(list: DictItem[], value: string): string {
  return list.find((i) => i.value === value)?.label ?? value ?? '-'
}

/** 取标签配色 */
export function typeOf(list: DictItem[], value: string): TagType {
  return list.find((i) => i.value === value)?.type ?? 'info'
}
