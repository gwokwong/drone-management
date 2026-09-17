import { AppRouteRecord } from '@/types/router'
import { droneRoutes } from './drone'

/**
 * 导出所有模块化路由
 *
 * 无人机管理平台业务菜单：
 * 数据大屏 / 无人机管理 / 储存管理 / 监控中心 / 报表管理 / 系统日志 / 系统管理
 */
export const routeModules: AppRouteRecord[] = [...droneRoutes]
