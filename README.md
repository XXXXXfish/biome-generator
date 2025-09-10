# 生物群系生成系统（Go + Vue 3）

本项目实现一个 10×10 的生物群系生成器，包含实时可视化与参数调节。后端使用 Go 提供 `/generate` 与 `/legend` 接口，前端使用 Vue 3 + Vite 渲染网格、图例与交互。

## 项目结构

- `biome-generator-backend/`：Go 后端（API 服务、生成算法）
- `biome-generator-ui/`：前端（Vue 3 + Vite 可视化界面）

## 功能说明

- 生成 10×10 网格，包含 5 种生物群系：森林、沙漠、海洋、山地、平原
- 可调节参数：
  - `moisture_spread`（西侧影响强度）
  - `temperature_spread`（北侧影响强度）
  - `climate_stability`（稳定区额外影响）
- 一键重新生成世界，实时更新可视化
- 悬浮提示显示每格生物群系与坐标
- 图例展示 5 种生物群系颜色与名称

## 核心算法（简述）

- 起点 `(0,0)` 随机生成任意生物群系
- 其他格依次生成（行优先），其分数来自：
  - 各生物群系基础分 20（对应自然出现率 20%）
  - 若有西侧格子，则将 `moisture_spread` 加到西侧相同生物群系分数上
  - 若有北侧格子，则将 `temperature_spread` 加到北侧相同生物群系分数上
  - 若西北相同生物群系，额外加上 `climate_stability`
- 将所有分数归一化为概率，按累积分布随机抽取最终生物群系

## 重难点与解决方案

- 概率归一化与累积选择：
  - 先累加各生物群系分数求总分，再计算 `score/total` 概率，并构造累积数组用于一次抽样
- 邻接影响与稳定区叠加：
  - 生成顺序保证 `(x-1,y)` 与 `(x,y-1)` 已存在；若同类则叠加稳定区权重
- 前后端数据字段一致性：
  - 后端 JSON 标签为小写 `name`、`color`；前端也应使用小写访问，避免出现 `undefined`
- CORS 与端口：
  - 后端允许 `http://localhost:5173` 来源，前端默认请求 `http://localhost:8080`
- 视觉与可用性：
  - 自适应高度的样式压缩，避免常见屏幕需要滚动
  - `overflow: visible` 处理，使边缘悬浮提示不被容器裁剪

## 快速开始

前提：Node.js ≥ 18，npm ≥ 9，Go ≥ 1.20（建议）。

1) 启动后端（项目根目录）：

```powershell
# 已有可执行文件（示例，Windows）
./biome-generator-backend/biome-generator-backend.exe

# 或源码方式
# cd biome-generator-backend
# go run ./...
```

看到日志 `Go Backend API Server starting on :8080` 即后端就绪。

2) 启动前端：

```powershell
cd biome-generator-ui
npm install
npm run dev
```

浏览器访问终端提示的地址（通常 `http://localhost:5173`）。

## 配置与端点

- 前端后端地址：在 `biome-generator-ui/src/App.vue` 中修改 `API_BASE_URL`
- 后端端点：
  - `GET /legend` → 返回生物群系图例（键为类型枚举值字符串）
  - `POST /generate` → 生成网格，表单字段：
    - `moisture_spread`（int）
    - `temperature_spread`（int）
    - `climate_stability`（int）

响应数据（示例，单格）：

```json
{
  "x": 2,
  "y": 3,
  "type": 1,
  "info": { "name": "沙漠", "color": "#F0E68C" }
}
```

## 常见问题（FAQ）

- 图例/颜色不显示或悬浮提示为 `undefined`：
  - 前端模板字段请使用小写 `name`、`color` 对应后端 JSON 标签
- 前端请求失败（CORS 或连接错误）：
  - 确认后端在 `:8080` 运行，且 CORS 允许 `http://localhost:5173`
  - 如更换前端端口，需同步更新后端 CORS
- 悬浮提示在边缘被裁剪：
  - 已设置容器 `overflow: visible`；如需更复杂的避让，可改用全局 tooltip
- 页面需要滚动：
  - 已内置媒体查询按视口高度收缩尺寸；可在 `App.vue` 进一步微调格子尺寸/间距

## 目录导航

- 前端详细说明：`biome-generator-ui/README.md`
- 后端代码：`biome-generator-backend/internal/biome/*.go`（类型、算法）、`internal/server/handler.go`（HTTP 处理）

## 生产构建（前端）

```bash
cd biome-generator-ui
npm run build
```

将 `dist/` 部署到静态服务器，确保正确的 API 地址与 CORS 设置。
