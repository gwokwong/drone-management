# DroneManagement 项目长期约定

## 项目定位
无人机管理平台。前端基于 art-design-pro（Vue 3 + Vite + TS + Element Plus），后端为自建 Go + Gin 服务，默认 SQLite、可切换 MySQL。

## 后端约定
- 技术栈：Go 1.22 + Gin + GORM；纯 Go SQLite 驱动（`glebarez/sqlite`，无需 cgo）。
- 模块名 `drone-server`，入口 `server/main.go`，默认端口 **8080**，接口根路径 `/api`。
- 数据库切换只改 `server/config/config.yaml` 的 `database.driver`（`sqlite` / `mysql`），无需改代码；启动时自动迁移与灌种子数据。
- 响应信封 `{ code, msg, data }`；分页请求 `current` / `size`，响应 `{ records, current, size, total }`。
- 认证：JWT，请求头 `Authorization: <token>`（也兼容 `Bearer <token>`），默认有效期 72 小时。
- 初始化账号：`admin / admin123`（R_SUPER，始终存在）；`operator / operator123`（R_OPERATOR，仅 `app.seedDemoData: true` 时写入）。
- **演示数据开关** `app.seedDemoData`（默认 **false**）：关闭时系统为不含任何测试数据的干净空系统，只初始化角色与管理员账号；置 true（或环境变量 `SEED_DEMO_DATA=true`）才写入无人机/储存室/告警/进出等演示数据。交付环境保持关闭。
- 配置与 SQLite 相对路径以**可执行文件所在目录**为基准解析，数据固定落在 `server/data/drone.db`，避免从任意工作目录启动时漂移。
- 接口文档：`server/API接口文档.md`。

## 前端约定
- 菜单为**前端模式**（`.env` 中 `VITE_ACCESS_MODE = frontend`），菜单来源 `src/router/modules/index.ts` → `drone.ts`，由 `ComponentLoader` 把 `component: '/xxx'` 解析到 `src/views` 下的组件。
- 开发代理：`.env.development` 的 `VITE_API_PROXY_URL = http://localhost:8080`，Vite 将 `/api` 转发到后端。
- 业务页面统一在 `src/views/drone/**`，接口统一在 `src/api/drone.ts`，状态字典在 `src/utils/drone/dict.ts`。
- 表格统一用 `useTable` + `ArtTable` + `ArtTableHeader` 模式。

## 易踩的坑
- `useTable` 的泛型是 **API 函数类型**（`useTable<TApiFn>`），**不是记录类型**。不要写 `useTable<Drone>(...)`，应省略泛型让 `apiFn` 推断。
- `models.User` 的 `Password` 与 `RoleCodes` 标记为 `json:"-"`：绑定/序列化时会丢失，需单独解析或用视图对象返回。
- 本环境 Bash 的 `cd` / `ls` / `head` 等 coreutils 不可用，Go 命令用 `go -C <dir>`；`pnpm exec` 解析不到 vite / vue-tsc，需直接调用 `node_modules/.bin/vite.CMD`。
- Element Plus 组件由 unplugin 全局自动导入，模板中可直接用 `El*`；但 `ElMessage` / `ElMessageBox` 等函数需在 `<script setup>` 中显式 import。
- 菜单父级 `meta.roles` 会过滤**整个子树**（`MenuProcessor.filterMenuByRoles`）：若某个子路由要对所有角色开放（如「个人中心」），父级就不要设 `roles`，把权限下沉到各个子菜单。
- 删除文件时要同时排查**间接引用**：桶文件（如 `utils/sys/index.ts` 的 `export *`）、`App.vue` 的全局调用、`components.d.ts` 的自动导入条目。本环境 `vue-tsc` 跑不起来，靠 `vite build` 才能暴露这类断裂。
- Python `urllib` 会走系统代理，访问 localhost 服务会返回 502；需 `urllib.request.build_opener(ProxyHandler({}))` 绕过。
