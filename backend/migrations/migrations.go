package migrations

import (
	"face-swap/config"
	"face-swap/models"
)

func Migrate() {
	// 自动迁移模式
	config.DB.AutoMigrate(&models.ImageSwapRecord{})

	// 检查表是否存在
	if !config.DB.Migrator().HasTable(&models.ImageSwapRecord{}) {
		// 创建表
		if err := config.DB.Migrator().CreateTable(&models.ImageSwapRecord{}); err != nil {
			panic("failed to create image_swap_record table")
		}
	}

}
