<script setup>
import { ref, onMounted } from 'vue';

// 响应式数据
const moistureSpread = ref(50);
const temperatureSpread = ref(30);
const climateStability = ref(100);
const grid = ref([]); // 存储 10x10 的生物群系网格数据
const biomeLegend = ref({}); // 存储生物群系图例数据
const isLoading = ref(false); // 控制加载状态
const error = ref(null); // 存储错误信息

// 后端 API 基础 URL
const API_BASE_URL = 'http://localhost:8080'; // 确保与 Go 后端运行的端口一致

// --- 获取并渲染图例 ---
const fetchAndRenderLegend = async () => {
  error.value = null; // 清除之前的错误
  try {
    const response = await fetch(`${API_BASE_URL}/legend`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    biomeLegend.value = await response.json();
  } catch (err) {
    console.error('Error fetching legend:', err);
    error.value = '加载图例失败。';
  }
};

// --- 生成并渲染世界网格 ---
const generateAndRenderWorld = async () => {
  isLoading.value = true; // 设置加载状态为 true
  error.value = null; // 清除之前的错误
  try {
    const formData = new URLSearchParams();
    formData.append('moisture_spread', moistureSpread.value);
    formData.append('temperature_spread', temperatureSpread.value);
    formData.append('climate_stability', climateStability.value);

    const response = await fetch(`${API_BASE_URL}/generate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: formData.toString(),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    grid.value = await response.json(); // 更新网格数据
    console.log('Backend response (updated grid.value):', grid.value);
  } catch (err) {
    console.error('Error generating world:', err);
    error.value = '生成世界失败，请检查服务器是否运行。';
  } finally {
    isLoading.value = false; // 无论成功或失败，都结束加载状态
  }
};

// --- 生命周期钩子：组件挂载后执行 ---
onMounted(() => {
  fetchAndRenderLegend(); // 页面加载时首先获取并渲染图例
  generateAndRenderWorld(); // 页面加载时自动生成一次世界并显示
});
</script>

<template>
  <div id="app-container">
    <h1>生物群系生成器</h1>

    <div class="controls">
      <label for="moisture_spread">西侧影响强度 (Moisture Spread %):</label>
      <input type="number" id="moisture_spread" v-model.number="moistureSpread" min="0" max="200">

      <label for="temperature_spread">北侧影响强度 (Temperature Spread %):</label>
      <input type="number" id="temperature_spread" v-model.number="temperatureSpread" min="0" max="200">

      <label for="climate_stability">气候稳定区影响强度 (Climate Stability %):</label>
      <input type="number" id="climate_stability" v-model.number="climateStability" min="0" max="200">

      <button @click="generateAndRenderWorld" :disabled="isLoading">
        {{ isLoading ? '生成中...' : '重新生成世界' }}
      </button>
    </div>

    <!-- 错误信息显示 -->
    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <!-- 生物群系图例区域 -->
    <div class="legend" id="biomeLegend">
      <h3>生物群系图例:</h3>
      <div v-for="(biome, typeKey) in biomeLegend" :key="typeKey" class="legend-item">
        <div class="legend-color-box" :style="{ backgroundColor: biome.color }"></div>
        <span>{{ biome.Name }}</span>
      </div>
    </div>

    <!-- 10x10 网格容器 -->
    <div class="grid-container" id="gridContainer">
      <!-- 使用 v-for 遍历 grid 数据来渲染网格单元 -->
      <!-- grid[x][y] 对应后端返回的结构，所以需要两层循环 -->
      <template v-for="(col, x) in grid" :key="x">
        <div v-for="(block, y) in col" :key="y"
             class="grid-cell"
             :style="{ backgroundColor: block.info.color }"
             :data-tooltip="`${block.info.name} (${block.x},${block.y})`">
          ({{ block.x }},{{ block.y }})
        </div>
      </template>
      <div v-if="isLoading && !error" class="loading-overlay">
          <div class="spinner"></div>
          <p>生成中...</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 整个应用的容器 */
html, body { height: 100%; }

#app-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 6px;
  background-color: #eef2f7; /* 淡蓝色背景 */
  color: #333;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  min-height: 100vh; /* 确保背景覆盖整个视口 */
  box-sizing: border-box;
  padding-bottom: 8px; /* 底部留白 */
}

h1 {
  color: #2c3e50; /* 深蓝色标题 */
  margin-bottom: 10px;
  font-size: 1.4rem;
}

.controls {
  background-color: #ffffff;
  padding: 12px;
  border-radius: 10px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08); /* 更柔和的阴影 */
  margin-bottom: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
  justify-content: center;
  width: 100%;
  max-width: 600px;
}

.controls label {
  font-weight: bold;
  color: #555;
  flex-shrink: 0; /* 防止标签收缩 */
  font-size: 0.95em;
}

.controls input[type="number"] {
  padding: 8px;
  border: 1px solid #cce0ff; /* 浅蓝色边框 */
  border-radius: 5px;
  width: 80px;
  text-align: center;
  font-size: 0.95em;
  transition: border-color 0.2s ease-in-out;
}

.controls input[type="number"]:focus {
  border-color: #007bff; /* 聚焦时深蓝色 */
  outline: none;
}

.controls button {
  padding: 10px 18px;
  background-color: #007bff; /* 蓝色按钮 */
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 0.95em;
  transition: background-color 0.2s ease-in-out, transform 0.1s ease-in-out;
}

.controls button:hover:not(:disabled) {
  background-color: #0056b3; /* 鼠标悬停时深蓝色 */
  transform: translateY(-2px); /* 悬停时轻微上浮 */
}

.controls button:disabled {
  background-color: #cccccc; /* 禁用状态的按钮颜色 */
  cursor: not-allowed;
  transform: none;
}

.error-message {
  color: #dc3545; /* 红色 */
  background-color: #f8d7da; /* 浅红色背景 */
  border: 1px solid #f5c6cb;
  border-radius: 5px;
  padding: 10px 20px;
  margin-bottom: 20px;
  font-weight: bold;
}

.legend {
  background-color: #ffffff;
  padding: 10px;
  border-radius: 10px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
  margin-bottom: 12px;
  text-align: center;
  width: 100%;
  max-width: 520px;
}

.legend h3 {
  margin-top: 0;
  color: #2c3e50;
  font-size: 1rem;
}

.legend-item {
  display: inline-flex;
  align-items: center;
  margin: 6px 12px;
  font-size: 0.9em;
  color: #444;
}

.legend-color-box {
  width: 18px;
  height: 18px;
  border: 1px solid #e0e0e0;
  margin-right: 10px;
  border-radius: 4px;
  box-shadow: inset 0 1px 3px rgba(0,0,0,0.05); /* 内部阴影 */
}

.grid-container {
  display: grid;
  grid-template-columns: repeat(10, 40px); /* 10列，每列40px */
  grid-template-rows: repeat(10, 40px);    /* 10行，每行40px */
  border: 1px solid #a0a0a0; /* 网格整体边框 */
  box-shadow: 0 6px 18px rgba(0,0,0,0.15); /* 显著的网格阴影 */
  background-color: #eee; /* 网格背景 */
  border-radius: 8px; /* 网格圆角 */
  overflow: visible; /* 允许悬浮提示越界显示 */
  position: relative; /* 用于加载覆盖层定位 */
}

.grid-cell {
  width: 40px;
  height: 40px;
  box-sizing: border-box; /* 包含 padding 和 border */
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px; /* 坐标文字大小 */
  color: rgba(0, 0, 0, 0.6); /* 坐标文字颜色 */
  position: relative;
  border: 1px solid rgba(220, 220, 220, 0.7); /* 区块间细线，略透明 */
  cursor: help; /* 鼠标悬停时显示帮助光标 */
  transition: border-color 0.1s ease-in-out;
}

.grid-cell:hover {
  border: 2px solid #3498db; /* 悬停时更明显的边框颜色 */
  z-index: 1; /* 确保 hover 时边框显示在最上层 */
  box-shadow: 0 0 8px rgba(52, 152, 219, 0.5); /* 悬停时蓝色光晕 */
}

/* 悬停提示样式 */
.grid-cell::after {
  content: attr(data-tooltip); /* 使用 data-tooltip 属性作为提示内容 */
  position: absolute;
  bottom: calc(100% + 5px); /* 向上显示，距离单元格 5px */
  left: 50%;
  transform: translateX(-50%);
  background-color: rgba(44, 62, 80, 0.9); /* 深色半透明背景 */
  color: white;
  padding: 8px 12px;
  border-radius: 6px;
  white-space: nowrap;
  opacity: 0;
  pointer-events: none; /* 允许鼠标穿透，不影响点击下层元素 */
  transition: opacity 0.2s ease-in-out, transform 0.2s ease-in-out;
  font-size: 13px;
  z-index: 100; /* 确保在最顶层 */
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
}

.grid-cell:hover::after {
  opacity: 1;
  transform: translateX(-50%) translateY(-5px); /* 悬停时轻微向上移动 */
}

/* 加载覆盖层 */
.loading-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(255, 255, 255, 0.8);
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    font-size: 1.2em;
    color: #007bff;
    z-index: 50;
}

.spinner {
    border: 4px solid rgba(0, 123, 255, 0.3);
    border-top: 4px solid #007bff;
    border-radius: 50%;
    width: 30px;
    height: 30px;
    animation: spin 1s linear infinite;
    margin-bottom: 10px;
}

@keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
}

/* 根据视口高度进一步压缩，避免出现滚动条 */
@media (max-height: 800px) {
  .grid-container {
    grid-template-columns: repeat(10, 36px);
    grid-template-rows: repeat(10, 36px);
  }
  .grid-cell {
    width: 36px;
    height: 36px;
    font-size: 9.5px;
  }
  .controls { padding: 10px; gap: 8px; }
  .legend { padding: 8px; }
}

@media (max-height: 700px) {
  .grid-container {
    grid-template-columns: repeat(10, 34px);
    grid-template-rows: repeat(10, 34px);
  }
  .grid-cell {
    width: 34px;
    height: 34px;
    font-size: 9px;
  }
  h1 { margin-bottom: 8px; font-size: 1.2rem; }
}

@media (max-height: 620px) {
  .grid-container {
    grid-template-columns: repeat(10, 30px);
    grid-template-rows: repeat(10, 30px);
  }
  .grid-cell {
    width: 30px;
    height: 30px;
    font-size: 8.5px;
  }
  .controls { padding: 8px; }
}
</style>
