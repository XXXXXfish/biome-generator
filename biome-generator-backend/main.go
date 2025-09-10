package main

import (
	"log"
	"net/http"

	"biome-generator-backend/internal/server" 
)

// corsMiddleware 是一个简单的 CORS 中间件，用于设置允许跨域请求的 HTTP 头
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置允许的源 (Origin)。在开发环境中，通常设置为前端的开发服务器地址
		// 生产环境中，你应该将其替换为你的前端生产环境域名
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // Vue 前端的默认开发端口
		// w.Header().Set("Access-Control-Allow-Origin", "*") // 不推荐在生产环境使用，表示允许所有源

		// 设置允许的 HTTP 方法
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		// 设置允许的 HTTP 头
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// 允许发送 Cookies (如果需要的话)
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

		// 对于预检请求 (OPTIONS 方法)，直接返回 200 OK
		// 浏览器会在发送真实的 POST/PUT/DELETE 请求前发送一个 OPTIONS 请求来检查 CORS 策略
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 将请求传递给下一个处理器
		next.ServeHTTP(w, r)
	})
}

func main() {
	// --- 移除前端静态文件服务和 HomeHandler ---
	// fs := http.FileServer(http.Dir("./web"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))
	// http.HandleFunc("/", server.HomeHandler) // 不再需要渲染 index.html

	// 创建一个路由器（Multiplexer）来处理 API 路由
	router := http.NewServeMux()

	// 注册 API 路由
	router.HandleFunc("/generate", server.GenerateHandler)  // 生成世界 API
	router.HandleFunc("/legend", server.BiomeLegendHandler) // 获取生物群系图例 API

	// 将 CORS 中间件应用到路由器上
	// 现在所有的 API 请求都会经过 CORS 处理
	handlerWithCORS := corsMiddleware(router)

	log.Println("Go Backend API Server starting on :8080")
	// 监听并服务 HTTP 请求，使用带有 CORS 的处理器
	if err := http.ListenAndServe(":8080", handlerWithCORS); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
