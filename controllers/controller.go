package controllers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-lat/models"
)

// FileController is a struct that contains a pointer to the database
type FileController struct {
	DB *gorm.DB
}

// UploadFile is a function that handles the upload of a single file
func (c *FileController) UploadFile(ctx *gin.Context) {
	/*
	  UploadFile function handles the upload of a single file.
	  It gets the file from the form data, saves it to the defined path,
	  generates a unique identifier for the file, saves the file metadata to the database,
	  and returns a success message and the file metadata.
	*/
	// Get the file from the form data
	title := ctx.PostForm("title")
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Define the path where the file will be saved
	filePath := filepath.Join("uploads", file.Filename)
	// Save the file to the defined path
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	// Generate a unique identifier for the file
	uuid := uuid.New().String()
	// Save file metadata to database
	fileMetadata := models.File{
		Title:    title,
		Filename: file.Filename,
		UUID:     uuid,
	}
	if err := c.DB.Create(&fileMetadata).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})
		return
	}
	// Return a success message and the file metadata
	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "Details": fileMetadata})
}
