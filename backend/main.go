package main

import (
	"face-swap/config"
	"face-swap/controllers"
	"face-swap/migrations"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// 加载环境变量
	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatal("Error loading app.env")
	}

	// 建立数据库连接
	config.ConnectDB()

	// 运行数据库迁移
	migrations.Migrate()

	// 初始化 Gin
	r := gin.Default()

	// 注册路由
	r.POST("/api/image_swap_records", controllers.CreateImageSwapRecord)
	r.GET("/api/image_swap_records", controllers.GetImageSwapRecords)
	r.GET("/api/image_swap_records/:id", controllers.GetImageSwapRecord)
	r.PUT("/api/image_swap_records/:id", controllers.UpdateImageSwapRecord)
	r.DELETE("/api/image_swap_records/:id", controllers.DeleteImageSwapRecord)

	r.POST("/api/swap", controllers.SwapImage)

	// 启动服务器
	r.Run()

}
