/**
 * 快速入口配置
 * 包含：应用列表、快速链接等配置
 *
 * 全部指向无人机管理平台的真实业务路由（对应 src/router/modules/drone.ts 中的 name）。
 */
import type { FastEnterConfig } from '@/types/config'

const fastEnterConfig: FastEnterConfig = {
  // 显示条件（屏幕宽度）
  minWidth: 1200,
  // 应用列表
  applications: [
    {
      name: '实时总览',
      description: '大屏实时数据与在库统计',
      icon: 'ri:dashboard-3-line',
      iconColor: '#377dff',
      enabled: true,
      order: 1,
      routeName: 'DashboardIndex'
    },
    {
      name: '无人机台账',
      description: 'RFID 标签与设备状态管理',
      icon: 'ri:plane-line',
      iconColor: '#13DEB9',
      enabled: true,
      order: 2,
      routeName: 'DroneList'
    },
    {
      name: '借还管理',
      description: 'RFID 扫码借出与归还',
      icon: 'ri:exchange-line',
      iconColor: '#ffb100',
      enabled: true,
      order: 3,
      routeName: 'DroneBorrow'
    },
    {
      name: '3D 导航图',
      description: '储存室库位三维导航',
      icon: 'ri:archive-2-line',
      iconColor: '#7A7FFF',
      enabled: true,
      order: 4,
      routeName: 'StorageMap'
    },
    {
      name: '告警管理',
      description: '断网 / 漏水等异常告警',
      icon: 'ri:alarm-warning-line',
      iconColor: '#ff3b30',
      enabled: true,
      order: 5,
      routeName: 'MonitorAlarm'
    },
    {
      name: '视频监控',
      description: '实时预览与开架录像',
      icon: 'ri:camera-line',
      iconColor: '#38C0FC',
      enabled: true,
      order: 6,
      routeName: 'MonitorVideo'
    },
    {
      name: '统计报表',
      description: '日报 / 周报 / 月报 / 年报',
      icon: 'ri:file-chart-line',
      iconColor: '#ff6b6b',
      enabled: true,
      order: 7,
      routeName: 'ReportIndex'
    },
    {
      name: '系统日志',
      description: '设备与操作日志查询',
      icon: 'ri:file-search-line',
      iconColor: '#FB7299',
      enabled: true,
      order: 8,
      routeName: 'LogIndex'
    }
  ],
  // 快速链接
  quickLinks: [
    {
      name: '个人中心',
      enabled: true,
      order: 1,
      routeName: 'UserCenter'
    },
    {
      name: '用户管理',
      enabled: true,
      order: 2,
      routeName: 'SystemUser'
    },
    {
      name: '角色管理',
      enabled: true,
      order: 3,
      routeName: 'SystemRole'
    },
    {
      name: '环境设备',
      enabled: true,
      order: 4,
      routeName: 'StorageEnv'
    },
    {
      name: '人员进出',
      enabled: true,
      order: 5,
      routeName: 'MonitorAccess'
    }
  ]
}

export default Object.freeze(fastEnterConfig)
