# 生物群系生成器前端（Vue 3 + Vite）

本项目是“生物群系生成系统”的前端界面，提供实时可视化、参数调节与重新生成能力。需要与后端 Go API 一起运行。

## 功能概览

- 实时展示 10×10 生物群系网格（森林、沙漠、海洋、山地、平原）
- 调节气候参数：`moisture_spread`、`temperature_spread`、`climate_stability`
- 一键“重新生成世界”，立即更新可视化
- 网格悬浮提示：显示生物群系名称与坐标
- 生物群系图例：显示 5 种类型及对应颜色
- 自适应高度样式：在常见屏幕高度下无需滚动

## 运行环境

- Node.js ≥ 18
- npm ≥ 9
- 后端 Go API（默认监听 `http://localhost:8080`）

## 后端启动

在项目根目录（包含 `biome-generator-backend`）按你的运行方式启动后端，确保日志显示 `Go Backend API Server starting on :8080`：

```sh
# 如果已构建可执行文件（示例）
# Windows PowerShell
./biome-generator-backend/biome-generator-backend.exe

# 或直接用 go 运行（需在 backend 目录）
# cd biome-generator-backend
# go run ./...
```

注意：后端启用了 CORS，允许来自 `http://localhost:5173` 的前端请求。

## 前端启动

在 `biome-generator-ui` 目录中：

```sh
npm install
npm run dev
```

启动成功后，访问终端提示的地址（通常是 `http://localhost:5173`）。

## 配置

- 前端默认请求后端 `http://localhost:8080`。如需修改，请在 `src/App.vue` 中调整 `API_BASE_URL` 常量。
- 如需修改前端开发端口或别名，见 `vite.config.js`。

## 使用说明

1. 打开网页后，前端会自动：
   - 调用 `/legend` 获取图例
   - 调用 `/generate` 获取 10×10 网格
2. 在顶部控制区调整参数：
   - 西侧影响强度（moisture_spread）
   - 北侧影响强度（temperature_spread）
   - 气候稳定区影响强度（climate_stability）
3. 点击“重新生成世界”按钮，网格会基于新参数重新计算并渲染。
4. 将鼠标悬停在任意格子上，可查看对应生物群系与 `(x,y)` 坐标。

参数建议范围：`0–200`。超过该范围也能工作，但可能导致概率偏置过强。

## 生物群系与颜色

- 森林（Forest）：#228B22
- 沙漠（Desert）：#F0E68C
- 海洋（Ocean）：#4682B4
- 山地（Mountain）：#A9A9A9
- 平原（Plains）：#90EE90

## 常见问题（Troubleshooting）

- 网格或图例不显示颜色 / 悬浮提示文字为 `undefined`：
  - 确认前端使用的字段名为小写 `name`、`color`（与后端 JSON 标签一致）。
- 前端无法请求后端：
  - 确认后端在 `:8080` 运行，且终端无错误日志。
  - 确认浏览器地址为 `http://localhost:5173`（与后端 CORS 允许的来源一致）。
  - 若更换了前端端口，请同步更新后端的 CORS 设置。
- 悬浮标签在边缘被截断：
  - 已将网格容器的 `overflow` 设为 `visible`；若仍有问题，可考虑用全局 `position: fixed` 的 tooltip 实现。
- 页面需要上下滚动才能完整展示：
  - 已内置自适应压缩（媒体查询）。若仍需要更紧凑，可在 `src/App.vue` 中进一步调整网格单元尺寸或控件间距。

## 代码结构（前端）

- `src/App.vue`：主界面与样式，负责参数输入、调用 API、渲染图例与网格
- `src/main.js`：Vue 应用入口
- `vite.config.js`：Vite 与插件配置

## 生产构建

```sh
npm run build
```

生成的静态文件位于 `dist/` 目录。你可以将其部署到任意静态资源服务器，注意配置好与后端 API 的跨域与地址。

