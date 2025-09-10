package server

import (
	"encoding/json"
	"log"
	"net/http"
	// "path/filepath" // 不再需要，因为不再提供静态文件
	"strconv"
	// "text/template" // 不再需要，因为不再提供静态文件

	// !!! 修改这里：确保导入路径与你的 Go 模块名一致 !!!
	"biome-generator-backend/internal/biome"
)

// HomeHandler 不再需要，因为前端由 Vue 独立提供服务
/*
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("web", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error parsing template %s: %v", tmplPath, err)
		return
	}
	tmpl.Execute(w, nil)
}
*/

// GenerateHandler 处理来自前端的生物群系生成请求
func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	moistureStr := r.FormValue("moisture_spread")
	tempStr := r.FormValue("temperature_spread")
	stabilityStr := r.FormValue("climate_stability")

	moisture, err := strconv.Atoi(moistureStr)
	if err != nil {
		moisture = 50 // 默认值
		log.Printf("Warning: Invalid moisture_spread value '%s', using default %d", moistureStr, moisture)
	}
	temperature, err := strconv.Atoi(tempStr)
	if err != nil {
		temperature = 30 // 默认值
		log.Printf("Warning: Invalid temperature_spread value '%s', using default %d", tempStr, temperature)
	}
	stability, err := strconv.Atoi(stabilityStr)
	if err != nil {
		stability = 100 // 默认值
		log.Printf("Warning: Invalid climate_stability value '%s', using default %d", stabilityStr, stability)
	}

	params := biome.GenerationParameters{
		MoistureSpread:    moisture,
		TemperatureSpread: temperature,
		ClimateStability:  stability,
	}

	log.Printf("Generating world with parameters: %+v", params)

	grid := biome.GenerateWorld(params) // 调用 biome 包的生成函数

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(grid); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding grid to JSON: %v", err)
	}
}

// BiomeLegendHandler 提供生物群系图例信息给前端
func BiomeLegendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(biome.BiomeMap); err != nil { // 使用 biome 包的 BiomeMap
		http.Error(w, "Failed to encode legend", http.StatusInternalServerError)
		log.Printf("Error encoding biome map to JSON: %v", err)
	}
}

