/**
 * 无人机管理平台 - 业务接口层
 *
 * 统一对接 Go + Gin 后端（/api 前缀），所有响应遵循 { code, msg, data } 信封。
 * 分页统一使用 current / size 请求，响应为 { records, current, size, total }。
 *
 * @module api/drone
 */

import request from '@/utils/http'
import axios from 'axios'
import { useUserStore } from '@/store/modules/user'

const { VITE_API_URL } = import.meta.env

/** 基础地址（去掉结尾斜杠，避免拼接出 //api 协议相对路径） */
const BASE = (VITE_API_URL || '').replace(/\/$/, '')

/** 拼接完整接口地址 */
const apiUrl = (path: string) => `${BASE}${path}`

/* -------------------------------------------------------------------------- */
/*                                   类型定义                                   */
/* -------------------------------------------------------------------------- */

/** 分页响应 */
export interface PageResult<T> {
  records: T[]
  current: number
  size: number
  total: number
}

/** 分页查询参数 */
export interface PageParams {
  current?: number
  size?: number
  [key: string]: unknown
}

/** 无人机 */
export interface Drone {
  id?: number
  code: string
  name: string
  model: string
  category: string
  storageStatus: 'in_position' | 'borrowed'
  deviceStatus: 'intact' | 'damaged' | 'scrapped'
  battery: number
  positionCode: string
  roomId?: number
  manufacturer: string
  purchaseDate?: string
  imageUrl?: string
  remark?: string
  createdAt?: string
}

/** 无人机统计 */
export interface DroneStats {
  total: number
  inPosition: number
  borrowed: number
  intact: number
  damaged: number
  scrapped: number
}

/** 借还记录 */
export interface BorrowRecord {
  id?: number
  droneId: number
  droneCode: string
  droneName: string
  borrowerId?: number
  borrowerName: string
  borrowTime: string
  expectedReturnTime?: string
  actualReturnTime?: string
  status: 'borrowing' | 'returned' | 'overdue'
  accessories?: string | string[]
  returnCondition?: string
  positionCode: string
  durationMin: number
  remark?: string
}

/** 借出请求（RFID 扫描） */
export interface ScanBorrowParams {
  code: string
  borrowerName: string
  borrowerId?: number
  expectedReturnHrs?: number
  accessories?: string[]
  remark?: string
}

/** 归还请求（RFID 扫描） */
export interface ScanReturnParams {
  code: string
  returnCondition?: 'intact' | 'damaged' | 'scrapped'
  remark?: string
}

/** 储存室 */
export interface StorageRoom {
  id?: number
  code: string
  name: string
  location: string
  capacity: number
  envTemp: number
  envHumidity: number
  status: number
  remark?: string
}

/** 库位视图 */
export interface RackPositionView {
  code: string
  row: number
  col: number
  layer: number
  occupied: boolean
  droneCode: string
  droneId?: number | null
}

/** 密集架视图 */
export interface RackView {
  id: number
  code: string
  name: string
  rows: number
  cols: number
  layers: number
  positions: RackPositionView[]
}

/** 3D 导航图 */
export interface RoomMap {
  room: StorageRoom
  racks: RackView[]
  totalPositions?: number
  usedPositions?: number
}

/** 告警事件 */
export interface AlarmEvent {
  id?: number
  type: string
  level: 'info' | 'warning' | 'critical'
  title: string
  content: string
  source?: string
  deviceId?: number
  roomId?: number
  status: 'active' | 'resolved' | 'ignored'
  resolvedAt?: string
  resolvedBy?: string
  createdAt?: string
}

/** 人员进出记录 */
export interface AccessRecord {
  id?: number
  cardNo: string
  personName: string
  type: 'card_auth_fail' | 'door_open' | 'door_close' | 'face_auth_pass'
  roomId?: number
  deviceId?: number
  imageUrl?: string
  createdAt?: string
}

/** 视频通道 */
export interface VideoChannel {
  id?: number
  name: string
  roomId?: number
  rackId?: number
  url: string
  status: 'online' | 'offline'
  remark?: string
}

/** 开架记录 */
export interface OpenRackRecord {
  id?: number
  videoChannelId?: number
  rackId?: number
  rackCode: string
  operator: string
  startTime: string
  endTime?: string
  videoUrl?: string
  status: 'recording' | 'finished'
  remark?: string
}

/** 环境设备 */
export interface EnvDevice {
  id?: number
  name: string
  type: string
  roomId?: number
  status: 'normal' | 'abnormal' | 'offline'
  value: number
  unit: string
  lastReportAt?: string
  remark?: string
}

/** 报表数据 */
export interface ReportData {
  type: string
  start: string
  end: string
  droneTotal: number
  droneBorrowed: number
  borrowCount: number
  overdueCount: number
  alarmCount: number
  accessCount: number
  trend: { date: string; borrow: number; alarm: number; access: number }[]
}

/** 系统日志 */
export interface SystemLog {
  id?: number
  category: string
  deviceType: string
  content: string
  operator: string
  level: string
  ip: string
  createdAt?: string
}

/** 用户 */
export interface UserItem {
  id?: number
  username: string
  password?: string
  nickname: string
  email: string
  phone?: string
  avatar?: string
  status: number
  roleCodes?: string[]
  remark?: string
  createdAt?: string
}

/** 角色 */
export interface RoleItem {
  id?: number
  code: string
  name: string
  permissions?: string[]
  remark?: string
  enabled: number
}

/** 大屏总览 */
export interface DashboardOverview {
  drone: {
    total: number
    inPosition: number
    borrowed: number
    intact: number
    damaged: number
    scrapped: number
    activeBorrow: number
  }
  alarm: { active: number; critical: number }
  accessToday: number
  env: { total: number; abnormal: number; offline: number }
  recentAlarms: AlarmEvent[]
  recentAccess: AccessRecord[]
  roomStats: { roomId: number; name: string; used: number; total: number }[]
}

/* -------------------------------------------------------------------------- */
/*                                  下载工具                                    */
/* -------------------------------------------------------------------------- */

/**
 * 下载文件（Excel 导出等），需携带 Authorization 头
 * @param url 接口地址（/api/...）
 * @param filename 保存的文件名
 */
export async function downloadFile(url: string, filename: string): Promise<void> {
  const { accessToken } = useUserStore()
  const res = await axios.get(apiUrl(url), {
    headers: { Authorization: accessToken },
    responseType: 'blob'
  })
  const blob = new Blob([res.data])
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(link.href)
}

/* -------------------------------------------------------------------------- */
/*                                   大屏总览                                   */
/* -------------------------------------------------------------------------- */

/** 获取大屏数据总览 */
export function fetchDashboardOverview() {
  return request.get<DashboardOverview>({ url: '/api/dashboard/overview' })
}

/* -------------------------------------------------------------------------- */
/*                                  无人机管理                                  */
/* -------------------------------------------------------------------------- */

/** 无人机列表 */
export function fetchDroneList(params: PageParams) {
  return request.get<PageResult<Drone>>({ url: '/api/drones', params })
}

/** 无人机统计 */
export function fetchDroneStats() {
  return request.get<DroneStats>({ url: '/api/drones/stats' })
}

/** 无人机详情 */
export function fetchDroneDetail(id: number) {
  return request.get<Drone>({ url: `/api/drones/${id}` })
}

/** 新增无人机 */
export function createDrone(data: Partial<Drone>) {
  return request.post<Drone>({ url: '/api/drones', params: data })
}

/** 更新无人机 */
export function updateDrone(id: number, data: Partial<Drone>) {
  return request.put<Drone>({ url: `/api/drones/${id}`, params: data })
}

/** 删除无人机 */
export function deleteDrone(id: number) {
  return request.del<null>({ url: `/api/drones/${id}` })
}

/* -------------------------------------------------------------------------- */
/*                              借还管理（RFID）                                */
/* -------------------------------------------------------------------------- */

/** 扫描借出 */
export function scanBorrow(data: ScanBorrowParams) {
  return request.post<BorrowRecord>({ url: '/api/borrow/scan', params: data })
}

/** 扫描归还 */
export function scanReturn(data: ScanReturnParams) {
  return request.post<BorrowRecord>({ url: '/api/return/scan', params: data })
}

/** 借还记录列表 */
export function fetchBorrowRecords(params: PageParams) {
  return request.get<PageResult<BorrowRecord>>({ url: '/api/borrow/records', params })
}

/* -------------------------------------------------------------------------- */
/*                              储存室 / 密集架                                 */
/* -------------------------------------------------------------------------- */

/** 储存室列表 */
export function fetchRoomList(params: PageParams) {
  return request.get<PageResult<StorageRoom>>({ url: '/api/storage-rooms', params })
}

/** 储存室 3D 导航图 */
export function fetchRoomMap(id: number) {
  return request.get<RoomMap>({ url: `/api/storage-rooms/${id}/map` })
}

/** 新增储存室 */
export function createRoom(data: Partial<StorageRoom>) {
  return request.post<StorageRoom>({ url: '/api/storage-rooms', params: data })
}

/** 更新储存室 */
export function updateRoom(id: number, data: Partial<StorageRoom>) {
  return request.put<StorageRoom>({ url: `/api/storage-rooms/${id}`, params: data })
}

/** 删除储存室 */
export function deleteRoom(id: number) {
  return request.del<null>({ url: `/api/storage-rooms/${id}` })
}

/** 密集架列表 */
export function fetchRackList(params: PageParams) {
  return request.get<PageResult<RackView>>({ url: '/api/racks', params })
}

/** 新增密集架（自动生成库位） */
export function createRack(data: Partial<RackView>) {
  return request.post<RackView>({ url: '/api/racks', params: data })
}

/* -------------------------------------------------------------------------- */
/*                                  环境设备                                    */
/* -------------------------------------------------------------------------- */

/** 环境设备列表 */
export function fetchEnvDevices(params: PageParams) {
  return request.get<PageResult<EnvDevice>>({ url: '/api/env-devices', params })
}

/** 环境设备汇总 */
export function fetchEnvSummary() {
  return request.get<Record<string, number>>({ url: '/api/env-devices/summary' })
}

/* -------------------------------------------------------------------------- */
/*                                   告警管理                                   */
/* -------------------------------------------------------------------------- */

/** 告警列表 */
export function fetchAlarmList(params: PageParams) {
  return request.get<PageResult<AlarmEvent>>({ url: '/api/alarms', params })
}

/** 告警统计 */
export function fetchAlarmStats() {
  return request.get<Record<string, number>>({ url: '/api/alarms/stats' })
}

/** 处理告警 */
export function resolveAlarm(id: number) {
  return request.post<null>({ url: `/api/alarms/${id}/resolve` })
}

/** 忽略告警 */
export function ignoreAlarm(id: number) {
  return request.post<null>({ url: `/api/alarms/${id}/ignore` })
}

/** 删除告警 */
export function deleteAlarm(id: number) {
  return request.del<null>({ url: `/api/alarms/${id}` })
}

/* -------------------------------------------------------------------------- */
/*                                  人员进出                                    */
/* -------------------------------------------------------------------------- */

/** 进出记录列表 */
export function fetchAccessList(params: PageParams) {
  return request.get<PageResult<AccessRecord>>({ url: '/api/access', params })
}

/** 进出统计 */
export function fetchAccessStats() {
  return request.get<Record<string, number>>({ url: '/api/access/stats' })
}

/** 新增进出记录（第三方对接） */
export function createAccess(data: Partial<AccessRecord>) {
  return request.post<AccessRecord>({ url: '/api/access', params: data })
}

/* -------------------------------------------------------------------------- */
/*                                  视频监控                                    */
/* -------------------------------------------------------------------------- */

/** 视频通道列表 */
export function fetchVideoList(params: PageParams) {
  return request.get<PageResult<VideoChannel>>({ url: '/api/videos', params })
}

/** 新增视频通道 */
export function createVideo(data: Partial<VideoChannel>) {
  return request.post<VideoChannel>({ url: '/api/videos', params: data })
}

/** 更新视频通道 */
export function updateVideo(id: number, data: Partial<VideoChannel>) {
  return request.put<VideoChannel>({ url: `/api/videos/${id}`, params: data })
}

/** 删除视频通道 */
export function deleteVideo(id: number) {
  return request.del<null>({ url: `/api/videos/${id}` })
}

/** 开架（开始录像） */
export function openRack(id: number) {
  return request.post<OpenRackRecord>({ url: `/api/videos/${id}/open-rack` })
}

/** 停止录像 */
export function stopOpenRack(id: number) {
  return request.post<null>({ url: `/api/open-racks/${id}/stop` })
}

/** 开架记录列表 */
export function fetchOpenRackList(params: PageParams) {
  return request.get<PageResult<OpenRackRecord>>({ url: '/api/open-racks', params })
}

/* -------------------------------------------------------------------------- */
/*                                   报表管理                                   */
/* -------------------------------------------------------------------------- */

/** 生成报表数据 */
export function fetchReport(params: { type: string }) {
  return request.get<ReportData>({ url: '/api/reports/generate', params })
}

/** 导出报表 Excel */
export function exportReport(type: string, filename?: string) {
  return downloadFile(`/api/reports/export?type=${type}`, filename || `报表_${type}.xlsx`)
}

/* -------------------------------------------------------------------------- */
/*                                  系统日志                                    */
/* -------------------------------------------------------------------------- */

/** 日志列表 */
export function fetchLogList(params: PageParams) {
  return request.get<PageResult<SystemLog>>({ url: '/api/logs', params })
}

/** 导出日志 Excel */
export function exportLogs(params: PageParams, filename?: string) {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([k, v]) => {
    if (v !== undefined && v !== null && v !== '') query.append(k, String(v))
  })
  return downloadFile(`/api/logs/export?${query.toString()}`, filename || '系统日志.xlsx')
}

/* -------------------------------------------------------------------------- */
/*                                  系统管理                                    */
/* -------------------------------------------------------------------------- */

/** 用户列表 */
export function fetchUserList(params: PageParams) {
  return request.get<PageResult<UserItem>>({ url: '/api/users', params })
}

/** 新增用户 */
export function createUser(data: Partial<UserItem>) {
  return request.post<UserItem>({ url: '/api/users', params: data })
}

/** 更新用户 */
export function updateUser(id: number, data: Partial<UserItem>) {
  return request.put<UserItem>({ url: `/api/users/${id}`, params: data })
}

/** 删除用户 */
export function deleteUser(id: number) {
  return request.del<null>({ url: `/api/users/${id}` })
}

/** 重置密码 */
export function resetPassword(id: number, password: string) {
  return request.post<null>({ url: `/api/users/${id}/reset-password`, params: { password } })
}

/** 角色列表 */
export function fetchRoleList(params: PageParams) {
  return request.get<PageResult<RoleItem>>({ url: '/api/roles', params })
}

/** 全部角色（不分页） */
export function fetchAllRoles() {
  return request.get<RoleItem[]>({ url: '/api/roles/all' })
}

/** 新增角色 */
export function createRole(data: Partial<RoleItem>) {
  return request.post<RoleItem>({ url: '/api/roles', params: data })
}

/** 更新角色 */
export function updateRole(id: number, data: Partial<RoleItem>) {
  return request.put<RoleItem>({ url: `/api/roles/${id}`, params: data })
}

/** 删除角色 */
export function deleteRole(id: number) {
  return request.del<null>({ url: `/api/roles/${id}` })
}
