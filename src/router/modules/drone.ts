import { AppRouteRecord } from '@/types/router'

/**
 * 无人机管理平台 - 业务路由（前端模式菜单）
 *
 * component 使用 '/xxx' 形式，由 ComponentLoader 动态解析到 src/views 下的组件。
 * meta.title 直接使用中文，formatMenuTitle 会原样返回（非 menus. 开头的键）。
 */

/** 数据大屏 */
const dashboardRoute: AppRouteRecord = {
  path: '/dashboard',
  name: 'Dashboard',
  component: '/index/index',
  meta: {
    title: '数据大屏',
    icon: 'ri:dashboard-3-line',
    keepAlive: true
  },
  children: [
    {
      path: 'index',
      name: 'DashboardIndex',
      component: '/drone/dashboard',
      meta: {
        title: '实时总览',
        icon: 'ri:dashboard-3-line',
        keepAlive: true
      }
    }
  ]
}

/** 无人机管理（RFID 借还 + 台账） */
const droneManageRoute: AppRouteRecord = {
  path: '/drone',
  name: 'Drone',
  component: '/index/index',
  meta: {
    title: '无人机管理',
    icon: 'ri:plane-line'
  },
  children: [
    {
      path: 'list',
      name: 'DroneList',
      component: '/drone/drone',
      meta: {
        title: '无人机台账',
        icon: 'ri:list-check-2',
        keepAlive: true
      }
    },
    {
      path: 'borrow',
      name: 'DroneBorrow',
      component: '/drone/borrow',
      meta: {
        title: '借还管理',
        icon: 'ri:exchange-line',
        keepAlive: true
      }
    }
  ]
}

/** 储存管理（3D 导航图 + 环境设备） */
const storageRoute: AppRouteRecord = {
  path: '/storage',
  name: 'Storage',
  component: '/index/index',
  meta: {
    title: '储存管理',
    icon: 'ri:building-2-line'
  },
  children: [
    {
      path: 'map',
      name: 'StorageMap',
      component: '/drone/storage',
      meta: {
        title: '3D 导航图',
        icon: 'ri:archive-2-line',
        keepAlive: true
      }
    },
    {
      path: 'env',
      name: 'StorageEnv',
      component: '/drone/env',
      meta: {
        title: '环境设备',
        icon: 'ri:temp-cold-line',
        keepAlive: true
      }
    }
  ]
}

/** 监控中心（告警 / 人员进出 / 视频） */
const monitorRoute: AppRouteRecord = {
  path: '/monitor',
  name: 'Monitor',
  component: '/index/index',
  meta: {
    title: '监控中心',
    icon: 'ri:shield-line'
  },
  children: [
    {
      path: 'alarm',
      name: 'MonitorAlarm',
      component: '/drone/alarm',
      meta: {
        title: '告警管理',
        icon: 'ri:alarm-warning-line',
        keepAlive: true
      }
    },
    {
      path: 'access',
      name: 'MonitorAccess',
      component: '/drone/access',
      meta: {
        title: '人员进出',
        icon: 'ri:door-open-line',
        keepAlive: true
      }
    },
    {
      path: 'video',
      name: 'MonitorVideo',
      component: '/drone/video',
      meta: {
        title: '视频监控',
        icon: 'ri:camera-line',
        keepAlive: true
      }
    }
  ]
}

/** 报表管理 */
const reportRoute: AppRouteRecord = {
  path: '/report',
  name: 'Report',
  component: '/index/index',
  meta: {
    title: '报表管理',
    icon: 'ri:file-chart-line'
  },
  children: [
    {
      path: 'index',
      name: 'ReportIndex',
      component: '/drone/report',
      meta: {
        title: '统计报表',
        icon: 'ri:file-chart-line',
        keepAlive: true
      }
    }
  ]
}

/** 系统日志 */
const logRoute: AppRouteRecord = {
  path: '/log',
  name: 'Log',
  component: '/index/index',
  meta: {
    title: '系统日志',
    icon: 'ri:file-list-3-line'
  },
  children: [
    {
      path: 'index',
      name: 'LogIndex',
      component: '/drone/log',
      meta: {
        title: '日志查询',
        icon: 'ri:file-search-line',
        keepAlive: true
      }
    }
  ]
}

/** 系统管理（用户 / 角色，按角色鉴权） */
const systemRoute: AppRouteRecord = {
  path: '/system',
  name: 'System',
  component: '/index/index',
  meta: {
    // 注意：此处不设 roles —— 父级 roles 会过滤整个子树，
    // 否则操作员将看不到「个人中心」（顶部用户菜单固定跳转 /system/user-center）。
    // 权限收敛到下方子菜单各自控制。
    title: '系统管理',
    icon: 'ri:settings-3-line'
  },
  children: [
    {
      path: 'user-center',
      name: 'UserCenter',
      component: '/drone/system/user-center',
      meta: {
        title: '个人中心',
        icon: 'ri:user-settings-line',
        keepAlive: true
      }
    },
    {
      path: 'user',
      name: 'SystemUser',
      component: '/drone/system/user',
      meta: {
        title: '用户管理',
        icon: 'ri:user-3-line',
        keepAlive: true,
        roles: ['R_SUPER', 'R_ADMIN']
      }
    },
    {
      path: 'role',
      name: 'SystemRole',
      component: '/drone/system/role',
      meta: {
        title: '角色管理',
        icon: 'ri:shield-user-line',
        keepAlive: true,
        roles: ['R_SUPER']
      }
    }
  ]
}

export const droneRoutes: AppRouteRecord[] = [
  dashboardRoute,
  droneManageRoute,
  storageRoute,
  monitorRoute,
  reportRoute,
  logRoute,
  systemRoute
]
