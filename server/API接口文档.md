# 无人机管理平台 · 后端接口技术文档

> 面向第三方二次开发与系统维护，描述 Go + Gin 后端全部 HTTP 接口的调用方式、参数与返回结构。

- **后端技术栈**：Go 1.22 + Gin + GORM
- **默认数据库**：SQLite（本地文件，零依赖，纯 Go 实现，无需 cgo）
- **可切换数据库**：MySQL（仅改配置，无需改代码）
- **服务默认端口**：`8080`
- **接口根路径**：`/api`

---

## 1. 通用约定

### 1.1 响应信封

所有业务接口统一返回如下结构：

```json
{
  "code": 200,
  "msg": "success",
  "data": {}
}
```

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | int | 业务状态码，`200` 表示成功 |
| `msg` | string | 提示信息，失败时为错误原因 |
| `data` | any | 业务数据，失败时通常为 `null` |

**状态码约定**

| code | 含义 |
| --- | --- |
| 200 | 成功 |
| 401 | 未认证（令牌缺失/无效/过期） |
| 403 | 无权限（角色不匹配） |
| 500 | 业务失败（如参数错误、记录不存在） |

### 1.2 分页约定

- **请求参数**：`current`（页码，从 1 开始，默认 1）、`size`（每页条数，默认 10）
- **响应结构**：

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "records": [],
    "current": 1,
    "size": 20,
    "total": 100
  }
}
```

### 1.3 认证方式

除登录接口外，所有接口需在请求头携带 JWT：

```
Authorization: <token>
```

也兼容 `Authorization: Bearer <token>` 形式。令牌有效期默认 **72 小时**。

### 1.4 健康检查

```
GET /health
```

返回 `{"status":"ok"}`，无需认证。

---

## 2. 认证与用户

### 2.1 登录

```
POST /api/auth/login
```

请求体：

```json
{ "userName": "admin", "password": "admin123" }
```

返回：

```json
{ "code": 200, "msg": "success", "data": { "token": "xxx", "refreshToken": "xxx" } }
```

### 2.2 退出登录

```
POST /api/auth/logout
```

### 2.3 获取当前用户信息

```
GET /api/user/info
```

返回：

```json
{
  "code": 200,
  "data": {
    "userId": 1,
    "userName": "admin",
    "email": "admin@drone.com",
    "avatar": "",
    "nickname": "系统管理员",
    "roles": ["R_SUPER"],
    "buttons": []
  }
}
```

### 2.4 用户管理

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/users` | 用户列表，支持 `userName`、`nickName`、`status` 模糊/精确过滤 |
| GET | `/api/users/:id` | 用户详情 |
| POST | `/api/users` | 新增用户 |
| PUT | `/api/users/:id` | 更新用户 |
| DELETE | `/api/users/:id` | 删除用户 |
| POST | `/api/users/:id/reset-password` | 重置密码 |

新增用户请求体：

```json
{
  "username": "operator",
  "password": "123456",
  "nickname": "操作员",
  "email": "op@drone.com",
  "phone": "13800000000",
  "status": 1,
  "roleCodes": ["R_OPERATOR"],
  "remark": ""
}
```

> 说明：`password` 留空时默认初始化为 `123456`；`roleCodes` 为角色编码数组。

更新用户请求体（`nickname`/`email`/`phone`/`avatar`/`status`/`roleCodes`/`remark`）：

```json
{ "nickname": "操作员", "status": 1, "roleCodes": ["R_OPERATOR"] }
```

重置密码：

```json
{ "password": "newPassword" }
```

### 2.5 角色管理

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/roles` | 角色分页列表，支持 `name`、`code` 过滤 |
| GET | `/api/roles/all` | 全部角色（不分页） |
| GET | `/api/roles/:id` | 角色详情 |
| POST | `/api/roles` | 新增角色 |
| PUT | `/api/roles/:id` | 更新角色 |
| DELETE | `/api/roles/:id` | 删除角色 |

角色请求体：

```json
{
  "code": "R_OPERATOR",
  "name": "操作员",
  "permissions": ["drone:view", "borrow:manage"],
  "enabled": 1,
  "remark": ""
}
```

**内置角色**：`R_SUPER`（超级管理员）、`R_ADMIN`（管理员）、`R_OPERATOR`（操作员）。

---

## 3. 无人机管理

### 3.1 无人机列表

```
GET /api/drones?keyword=&storageStatus=&deviceStatus=&roomId=&category=
```

| 参数 | 说明 |
| --- | --- |
| `keyword` | 模糊匹配名称 / 编号 / 型号 |
| `storageStatus` | `in_position` 在位 / `borrowed` 借出 |
| `deviceStatus` | `intact` 完好 / `damaged` 损坏 / `scrapped` 报废 |
| `roomId` | 所属储存室 ID |
| `category` | 分类 |

返回 `records` 中无人机对象：

```json
{
  "ID": 1,
  "code": "RFID-001",
  "name": "测绘无人机 M300",
  "model": "M300 RTK",
  "category": "多旋翼",
  "storageStatus": "in_position",
  "deviceStatus": "intact",
  "battery": 96,
  "positionCode": "A-01-1",
  "roomId": 1,
  "manufacturer": "DJI",
  "imageUrl": "",
  "remark": ""
}
```

### 3.2 无人机统计

```
GET /api/drones/stats
```

返回：

```json
{
  "total": 8,
  "inPosition": 6,
  "borrowed": 2,
  "intact": 7,
  "damaged": 1,
  "scrapped": 0,
  "byCategory": [{ "category": "多旋翼", "count": 5 }]
}
```

### 3.3 无人机增删改查

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/drones/:id` | 详情 |
| POST | `/api/drones` | 新增 |
| PUT | `/api/drones/:id` | 更新 |
| DELETE | `/api/drones/:id` | 删除 |

---

## 4. 借还管理（RFID / 电子标签）

### 4.1 扫描借出

```
POST /api/borrow/scan
```

请求体：

```json
{
  "code": "RFID-001",
  "borrowerName": "张三",
  "borrowerId": 1,
  "expectedReturnHrs": 24,
  "accessories": ["遥控器", "电池", "充电器"],
  "remark": ""
}
```

**处理逻辑**：

1. 校验无人机存在且当前为「在位」，否则返回错误；
2. 自动生成借还记录：借用人、借出时间、预计归还时间、随行配件；
3. 无人机状态置为 `borrowed`；
4. 释放原库位（库位 `drone_id` 置空、`status` 置 0）；
5. 写入系统操作日志。

返回借还记录对象。

### 4.2 扫描归还

```
POST /api/return/scan
```

请求体：

```json
{ "code": "RFID-001", "returnCondition": "intact", "remark": "" }
```

`returnCondition`：`intact` / `damaged` / `scrapped`，留空默认 `intact`。

**处理逻辑**：

1. 找到该无人机最近一条 `borrowing` 记录；
2. 记录实际归还时间、计算使用时长（`durationMin`）；
3. 无人机回到原位并重新占用库位；
4. 按归还状态更新设备状态；
5. 记录开门日志，便于与门禁 / 密集架联动。

### 4.3 借还记录

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/borrow/records` | 列表，支持 `droneCode`、`borrowerName`、`status`、`start`、`end` |
| GET | `/api/borrow/records/:id` | 详情 |

`start` / `end` 使用 **RFC3339** 时间格式，例如 `2026-09-16T00:00:00Z`。

记录对象：

```json
{
  "droneCode": "RFID-001",
  "droneName": "测绘无人机 M300",
  "borrowerName": "张三",
  "borrowTime": "2026-09-16T10:00:00+08:00",
  "expectedReturnTime": "2026-09-17T10:00:00+08:00",
  "actualReturnTime": null,
  "status": "borrowing",
  "accessories": "[\"遥控器\",\"电池\"]",
  "positionCode": "A-01-1",
  "durationMin": 125
}
```

> `status` 会在查询时实时计算：超过预计归还时间自动置为 `overdue`，并实时刷新 `durationMin`。

---

## 5. 储存室与密集架

### 5.1 储存室

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/storage-rooms` | 列表，支持 `name` |
| GET | `/api/storage-rooms/:id` | 详情 |
| GET | `/api/storage-rooms/:id/map` | **3D 导航图数据** |
| POST | `/api/storage-rooms` | 新增 |
| PUT | `/api/storage-rooms/:id` | 更新 |
| DELETE | `/api/storage-rooms/:id` | 删除 |

### 5.2 3D 导航图

```
GET /api/storage-rooms/:id/map
```

返回密集架及全部库位，用于图形化盘点与导航：

```json
{
  "room": { "ID": 1, "code": "R001", "name": "一号储存室", "location": "A 栋 1 层", "envTemp": 22.5, "envHumidity": 45 },
  "racks": [
    {
      "id": 1,
      "code": "A",
      "name": "A 区密集架",
      "rows": 4,
      "cols": 6,
      "layers": 3,
      "positions": [
        { "code": "A-01-1", "row": 1, "col": 1, "layer": 1, "occupied": true, "droneCode": "RFID-001", "droneId": 1 }
      ]
    }
  ],
  "totalPositions": 72,
  "usedPositions": 6
}
```

### 5.3 密集架

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/racks` | 列表，支持 `roomId` |
| POST | `/api/racks` | 新增（**自动生成全部库位**） |
| PUT | `/api/racks/:id` | 更新 |
| DELETE | `/api/racks/:id` | 删除 |

新增密集架请求体：

```json
{ "roomId": 1, "code": "A", "name": "A 区密集架", "rows": 4, "cols": 6, "layers": 3 }
```

系统按 `列-行-层` 自动生成库位编码（如 `A-01-1`）。

### 5.4 环境设备

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/env-devices` | 列表，支持 `roomId`、`type`、`status` |
| GET | `/api/env-devices/summary` | 汇总统计 |
| POST | `/api/env-devices` | 新增 |
| PUT | `/api/env-devices/:id` | 更新 |
| DELETE | `/api/env-devices/:id` | 删除 |

设备类型：`temperature` 温度、`humidity` 湿度、`smoke` 烟感、`water` 水浸、`door` 门禁、`power` 供电。
设备状态：`normal` / `abnormal` / `offline`。

---

## 6. 告警管理

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/alarms` | 列表，支持 `type`、`level`、`status`、`roomId`、`source`、`start`、`end` |
| GET | `/api/alarms/stats` | 告警统计 |
| POST | `/api/alarms` | 手动创建告警 |
| POST | `/api/alarms/:id/resolve` | 处理告警 |
| POST | `/api/alarms/:id/ignore` | 忽略告警 |
| DELETE | `/api/alarms/:id` | 删除告警 |

**告警类型**

| 值 | 含义 |
| --- | --- |
| `network_disconnect` | 断网异常 |
| `water_leak` | 漏水异常 |
| `door_alarm` | 门锁异常 |
| `device_fault` | 设备故障 |
| `low_battery` | 低电量 |
| `env_abnormal` | 环境异常 |

**告警级别**：`info` / `warning` / `critical`
**告警状态**：`active` 未处理 / `resolved` 已处理 / `ignored` 已忽略
**告警来源**：`system` 系统 / `third-party` 第三方 / `device` 设备

---

## 7. 人员进出

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/access` | 列表，支持 `type`、`roomId`、`personName`、`cardNo`、`start`、`end` |
| GET | `/api/access/stats` | 进出统计 |
| POST | `/api/access` | 新增进出记录 |

**事件类型**

| 值 | 含义 |
| --- | --- |
| `card_auth_fail` | 卡号认证失败（超次报警） |
| `door_open` | 门锁打开 |
| `door_close` | 门锁关闭 |
| `face_auth_pass` | 人脸认证通过 |

统计接口返回：

```json
{
  "total": 120,
  "today": 8,
  "cardFail": 2,
  "doorOpen": 40,
  "doorClose": 40,
  "facePass": 38
}
```

---

## 8. 视频监控与开架录像

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/videos` | 摄像头列表，支持 `roomId` |
| POST | `/api/videos` | 新增摄像头 |
| PUT | `/api/videos/:id` | 更新摄像头 |
| DELETE | `/api/videos/:id` | 删除摄像头 |
| POST | `/api/videos/:id/open-rack` | **开架并开始录像** |
| POST | `/api/open-racks/:id/stop` | 结束录像 |
| GET | `/api/open-racks` | 开架录像记录，支持 `rackId` |

摄像头对象：

```json
{ "name": "A 区摄像头", "roomId": 1, "rackId": 1, "url": "rtsp://192.168.1.10:554/stream", "status": "online" }
```

---

## 9. 报表与日志

### 9.1 生成报表

```
GET /api/reports/generate?type=daily
```

`type`：`daily` 日报 / `weekly` 周报 / `monthly` 月报 / `yearly` 年报。

返回：

```json
{
  "type": "daily",
  "start": "2026-09-16T00:00:00+08:00",
  "end": "2026-09-16T17:00:00+08:00",
  "droneTotal": 8,
  "droneBorrowed": 2,
  "borrowCount": 5,
  "overdueCount": 1,
  "alarmCount": 3,
  "accessCount": 12,
  "trend": [{ "date": "09-16", "borrow": 5, "alarm": 3, "access": 12 }]
}
```

### 9.2 导出报表 Excel

```
GET /api/reports/export?type=daily
```

返回 `.xlsx` 文件流（需携带 `Authorization` 头）。

### 9.3 系统日志

```
GET /api/logs?category=&deviceType=&level=&keyword=&start=&end=
```

- `category`：`device` 设备 / `user` 用户 / `alarm` 告警 / `operation` 操作
- `keyword`：模糊匹配日志内容与操作人
- `start` / `end`：RFC3339 时间范围

### 9.4 导出日志 Excel

```
GET /api/logs/export?category=&start=&end=...
```

---

## 10. 大屏数据总览

```
GET /api/dashboard/overview
```

一次返回大屏所需的全部聚合数据：

```json
{
  "drone": { "total": 8, "inPosition": 6, "borrowed": 2, "intact": 7, "damaged": 1, "scrapped": 0, "activeBorrow": 2 },
  "alarm": { "active": 3, "critical": 1 },
  "accessToday": 8,
  "env": { "total": 12, "abnormal": 1, "offline": 0 },
  "recentAlarms": [],
  "recentAccess": [],
  "roomStats": [{ "roomId": 1, "name": "一号储存室", "used": 6, "total": 72 }]
}
```

---

## 11. 第三方对接接口（预留）

用于门禁、环控、RFID 读写器等第三方系统主动推送数据。

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/integration/alarm/push` | 推送告警事件 |
| POST | `/api/integration/access/push` | 推送人员进出记录 |
| POST | `/api/integration/env/push` | 推送环境设备数据 |
| POST | `/api/integration/rfid/borrow` | RFID 设备触发借出 |
| POST | `/api/integration/rfid/return` | RFID 设备触发归还 |

推送告警示例：

```json
{
  "type": "water_leak",
  "level": "critical",
  "title": "一号储存室检测到漏水",
  "content": "水浸传感器 WS-01 触发",
  "deviceId": 3,
  "roomId": 1
}
```

推送环境数据示例：

```json
{ "name": "温湿度传感器", "type": "temperature", "roomId": 1, "status": "normal", "value": 22.5, "unit": "℃" }
```

---

## 12. 数据库切换（SQLite ↔ MySQL）

配置文件：`server/config/config.yaml`

```yaml
database:
  driver: sqlite            # sqlite | mysql
  sqlite:
    path: data/drone.db
  mysql:
    host: 127.0.0.1
    port: 3306
    user: root
    password: "123456"
    dbname: drone
```

**切换到 MySQL 仅需两步**：

1. 将 `driver` 改为 `mysql`，并填写 MySQL 连接信息；
2. 重启服务，GORM 会自动执行表结构迁移与初始化数据。

无需修改任何业务代码。

> 说明：SQLite 使用纯 Go 驱动（`github.com/glebarez/sqlite`），不需要 cgo，交叉编译与部署无额外依赖。

---

## 13. 初始化账号

服务首次启动会自动建表并写入**基础数据**（角色 + 管理员账号）：

| 账号 | 密码 | 角色 | 写入条件 |
| --- | --- | --- | --- |
| `admin` | `admin123` | 超级管理员（R_SUPER） | 始终写入 |
| `operator` | `operator123` | 操作员（R_OPERATOR） | 仅 `app.seedDemoData: true` 时写入 |

### 13.1 演示数据开关

配置项 `app.seedDemoData`（见 `config/config.yaml`）控制是否写入**演示业务数据**：
无人机、储存室/密集架/库位、环境设备、摄像头、告警事件、人员进出记录等。

| 取值 | 行为 |
| --- | --- |
| `false`（默认） | 仅初始化角色与管理员账号，系统为**不含任何测试数据的干净空系统**，交付即用 |
| `true` | 写入演示业务数据，便于演示与联调 |

也可用环境变量临时覆盖，无需改动配置文件：

```bash
SEED_DEMO_DATA=true ./drone-server.exe
```

> 已存在的数据库不会被覆盖：各张表仅在为空时才写入种子数据。
> 如需彻底重置，删除 SQLite 文件（`data/drone.db`）后重启服务即可。

---

## 14. 启动方式

```bash
# 开发运行
go run ./main.go

# 编译
go build -o drone-server ./...

# 指定配置文件
CONFIG_PATH=./config/config.yaml ./drone-server
```

服务启动后监听 `:8080`，前端通过 `VITE_API_PROXY_URL` 代理 `/api` 到该地址。
