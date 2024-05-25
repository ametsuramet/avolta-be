package handler

import (
	"avolta/config"
	"avolta/util"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
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

	flipped := c.PostForm("flipped")

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

	if flipped == "1" {

		srcFile, err := os.Open(destination)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer srcFile.Close()

		// Detect the file type
		ext := filepath.Ext(file.Filename)
		var img image.Image

		switch ext {

		case ".png":
			img, err = png.Decode(srcFile)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode PNG image"})
				return
			}
		default:
			img, err = jpeg.Decode(srcFile)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode JPEG image"})
				return
			}
		}

		flippedImg := flipImageHorizontally(img)

		newDest := filepath.Join(assetsFolder, "flipped_"+filename)

		// Save the flipped image
		outFile, err := os.Create(newDest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create output file"})
			return
		}
		defer outFile.Close()

		switch ext {
		case ".png":
			err = png.Encode(outFile, flippedImg)
		default:
			err = jpeg.Encode(outFile, flippedImg, nil)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode image"})
			return
		}
		destination = newDest
	}

	util.ResponseSuccess(c, "file uploaded", gin.H{
		"path":     destination,
		"filename": file.Filename,
		"url":      fmt.Sprintf("%s/%s", config.App.Server.BaseURL, destination),
	}, nil)
}

func flipImageHorizontally(img image.Image) image.Image {
	bounds := img.Bounds()
	flipped := image.NewRGBA(bounds)
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			flipped.Set(bounds.Dx()-x-1, y, img.At(x, y))
		}
	}
	return flipped
}
