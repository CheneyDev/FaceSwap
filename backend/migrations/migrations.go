package migrations

import (
	"face-swap/database"
	"face-swap/models"
)

func Migrate() {
	// 自动迁移模式
	database.DB.AutoMigrate(&models.ImageSwapRecord{})

	// 检查表是否存在
	if !database.DB.Migrator().HasTable(&models.ImageSwapRecord{}) {
		// 创建表
		if err := database.DB.Migrator().CreateTable(&models.ImageSwapRecord{}); err != nil {
			panic("failed to create image_swap_record table")
		}
	}

}
