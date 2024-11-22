package handlers

import (
	"fmt"
	"math"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yantology/go-csv-to-postgres-coroutine/config"
	"github.com/yantology/go-csv-to-postgres-coroutine/services"
)

func HandleFileUpload(db *config.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
			return
		}

		// Validate file extension
		if filepath.Ext(file.Filename) != ".csv" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only CSV files are allowed"})
			return
		}

		// Process the uploaded file
		err = services.ProcessCSVFile(db, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File processed successfully"})
		duration := time.Since(start)
		fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")
	}
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
