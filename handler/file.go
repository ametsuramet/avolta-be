package handler

import (
	"avolta/config"
	"avolta/util"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"path"

	"github.com/gin-gonic/gin"
)

func FileUploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assetsFolder := "./assets/images/"
	if _, err := os.Stat(assetsFolder); os.IsNotExist(err) {
		os.Mkdir(assetsFolder, os.ModePerm)
	}

	// Generate a unique filename
	filename := fmt.Sprintf("%v%s", util.GetCurrentTimestamp(), path.Ext(file.Filename))

	// Save the uploaded file to the assets folder
	destination := filepath.Join(assetsFolder, filename)
	if err := c.SaveUploadedFile(file, destination); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	util.ResponseSuccess(c, "file uploaded", gin.H{
		"path":     destination,
		"filename": file.Filename,
		"url":      fmt.Sprintf("%s/%s", config.App.Server.BaseURL, destination),
	}, nil)
}
