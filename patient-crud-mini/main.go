package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Patient struct {
	ID              uint   `json:"id" gorm:"primaryKey"`
	Name            string `json:"name" gorm:"size:100;not null"`
	Phone           string `json:"phone" gorm:"size:30;not null"`
	AdmissionNumber string `json:"admission_number" gorm:"size:50;not null"`
}

type PatientInput struct {
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	AdmissionNumber string `json:"admission_number"`
}

func readInput(c *gin.Context) (PatientInput, bool) {
	var input PatientInput
	// Context 保存本次请求和响应；JSON Binding 把请求体写入 input。
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体必须是有效 JSON"})
		return input, false
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Phone = strings.TrimSpace(input.Phone)
	input.AdmissionNumber = strings.TrimSpace(input.AdmissionNumber)
	if input.Name == "" || input.Phone == "" || input.AdmissionNumber == "" || len([]rune(input.Name)) > 100 || len([]rune(input.Phone)) > 30 || len([]rune(input.AdmissionNumber)) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "姓名、手机号和住院号必填，且长度不能超过字段限制"})
		return input, false
	}
	return input, true
}

func readID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID 必须是正整数"})
		return 0, false
	}
	return id, true
}

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("请设置 MYSQL_DSN")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	if err := db.AutoMigrate(&Patient{}); err != nil {
		log.Fatalf("创建表失败: %v", err)
	}

	r := gin.Default()
	r.POST("/patients", func(c *gin.Context) {
		input, ok := readInput(c)
		if !ok {
			return
		}
		patient := Patient{Name: input.Name, Phone: input.Phone, AdmissionNumber: input.AdmissionNumber}
		// 指针让 GORM 填入数据库生成的 ID。
		if err := db.Create(&patient).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "新增失败"})
			return
		}
		c.JSON(http.StatusCreated, patient)
	})
	r.GET("/patients", func(c *gin.Context) {
		patients := []Patient{}
		// Find 查询全部行，并通过切片指针写入结果。
		if err := db.Find(&patients).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
			return
		}
		c.JSON(http.StatusOK, patients)
	})
	r.PUT("/patients/:id", func(c *gin.Context) {
		id, ok := readID(c)
		if !ok {
			return
		}
		input, ok := readInput(c)
		if !ok {
			return
		}
		var patient Patient
		// First 按主键查询；区分不存在和其他数据库错误。
		err := db.First(&patient, id).Error
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "患者不存在"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
			return
		}
		patient.Name, patient.Phone, patient.AdmissionNumber = input.Name, input.Phone, input.AdmissionNumber
		if err := db.Save(&patient).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "修改失败"})
			return
		}
		c.JSON(http.StatusOK, patient)
	})
	r.DELETE("/patients/:id", func(c *gin.Context) {
		id, ok := readID(c)
		if !ok {
			return
		}
		result := db.Delete(&Patient{}, id)
		// Error 表示执行失败；RowsAffected 为 0 表示没有删除任何行。
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "患者不存在"})
			return
		}
		c.Status(http.StatusNoContent)
	})
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
