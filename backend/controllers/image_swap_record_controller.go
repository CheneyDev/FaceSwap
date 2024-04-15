package controllers

import (
	"face-swap/config"
	"face-swap/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateImageSwapRecord 创建记录
func CreateImageSwapRecord(c *gin.Context) {
	var record models.ImageSwapRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		fmt.Println("Error binding JSON:", err) // 打印错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//打印record
	fmt.Println("Received record:", record)
	if err := config.DB.Create(&record).Error; err != nil {
		fmt.Println("Error creating record:", err) // 打印错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
		return
	}
	c.JSON(http.StatusOK, record)
}

// GetImageSwapRecords 获取所有记录
func GetImageSwapRecords(c *gin.Context) {
	var records []models.ImageSwapRecord
	config.DB.Find(&records)
	c.JSON(200, records)
}

// GetImageSwapRecord 获取单个记录
func GetImageSwapRecord(c *gin.Context) {
	var record models.ImageSwapRecord
	if err := config.DB.Where("id = ?", c.Param("id")).First(&record).Error; err != nil {
		c.JSON(404, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(200, record)
}

// UpdateImageSwapRecord 更新记录
func UpdateImageSwapRecord(c *gin.Context) {
	var record models.ImageSwapRecord
	if err := config.DB.Where("id = ?", c.Param("id")).First(&record).Error; err != nil {
		c.JSON(404, gin.H{"error": "Record not found"})
		return
	}
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&record)
	c.JSON(200, record)
}

// DeleteImageSwapRecord 删除记录
func DeleteImageSwapRecord(c *gin.Context) {
	var record models.ImageSwapRecord
	if err := config.DB.Where("id = ?", c.Param("id")).First(&record).Error; err != nil {
		c.JSON(404, gin.H{"error": "Record not found"})
		return
	}
	config.DB.Delete(&record)
	c.JSON(200, gin.H{"message": "Record deleted"})
}
