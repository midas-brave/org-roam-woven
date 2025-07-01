package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func containsDotDot(s string) bool {
	for part := range strings.SplitSeq(s, string(os.PathSeparator)) {
		if part == ".." {
			return true
		}
	}
	return false
}

func isSubPath(basePath string, subPath string) bool {
	if basePath == subPath {
		return true
	}
	rel, err := filepath.Rel(basePath, subPath)
	if err != nil {
		return false
	}
	return !containsDotDot(rel) && rel != "."
}

func validateAndCleanPath(c *gin.Context, requestedPath string) (string, error) {
	baseDir := os.Getenv("ROAM_DIRECTORY")
	if baseDir == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ROAM_DIR environment variable is not set.",
		})
		return "", os.ErrNotExist
	}
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error getting absolute base directory path.",
			"error":   err.Error(),
		})
		return "", err
	}

	cleanPath := filepath.Clean(filepath.Join(absBaseDir, requestedPath))
	absRequestedPath, err := filepath.Abs(cleanPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid path.",
			"error":   err.Error(),
		})
		return "", err
	}
	if !strings.HasPrefix(absRequestedPath, absBaseDir) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Access denied.",
			"path":    requestedPath,
		})
		return "", os.ErrPermission
	}
	if !isSubPath(absBaseDir, absRequestedPath) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Access denied.",
			"path":    requestedPath,
		})
		return "", os.ErrPermission
	}

	return absRequestedPath, nil
}

func listFilePath(c *gin.Context, filePath string) {
	requestedPath := c.Param("filePath")
	files, err := os.ReadDir(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not read directory",
			"error":   err.Error(),
		})
		return
	}

	var fileList []gin.H
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			continue
		}
		fileList = append(fileList, gin.H{
			"name":     file.Name(),
			"is_dir":   file.IsDir(),
			"size":     fileInfo.Size(),
			"mod_time": fileInfo.ModTime().Format("2006-01-02 15:04:05"),
			"mode":     fileInfo.Mode().String(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"path":  requestedPath,
		"files": fileList,
	})
}

func viewFileContent(c *gin.Context, filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not read file content.",
			"error":   err.Error(),
		})
		return
	}

	c.Data(http.StatusOK, "text/plain", content)
}

func viewFilesPath(c *gin.Context) {
	requestedPath := c.Param("filePath")
	cleanAbsPath, err := validateAndCleanPath(c, requestedPath)
	if err != nil {
		return
	}

	fileInfo, err := os.Stat(cleanAbsPath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "File not found",
			"path":    requestedPath,
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error accessing path",
			"error":   err.Error(),
		})
		return
	}
	if fileInfo.IsDir() {
		listFilePath(c, cleanAbsPath)
	} else {
		viewFileContent(c, cleanAbsPath)
	}
}

func createFile(c *gin.Context) {
	requestedPath := c.Param("filePath")
	file, err := c.FormFile("file")
	if file == nil || err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	cleanAbsPath, err := validateAndCleanPath(c, requestedPath)
	if err != nil {
		return
	}
	fileInfo, err := os.Stat(cleanAbsPath)
	if os.IsNotExist(err) {

		createErr := c.SaveUploadedFile(file, cleanAbsPath)
		if createErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Save file error.",
				"path":    requestedPath,
			})
			return
		}
		c.Status(http.StatusCreated)
	} else {
		if fileInfo.IsDir() {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Request path is directory.",
				"path":    requestedPath,
			})
			return
		}
		createErr := c.SaveUploadedFile(file, cleanAbsPath)
		if createErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Save file error.",
				"path":    requestedPath,
			})
			return
		}
		c.Status(http.StatusNoContent)

	}
}

func deleteFile(c *gin.Context) {
	requestedPath := c.Param("filePath")
	cleanAbsPath, err := validateAndCleanPath(c, requestedPath)
	if err != nil {
		return
	}
	fileInfo, err := os.Stat(cleanAbsPath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "File not found",
			"path":    requestedPath,
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error accessing path",
			"error":   err.Error(),
		})
		return
	}
	if fileInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not allow delete directory.",
			"path":    requestedPath,
		})
		return
	} else {
		if err := os.Remove(cleanAbsPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error accessing path",
				"error":   err.Error(),
			})
		}
		return
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to org-roam-woven!",
		})
	})

	{
		unstable := r.Group("api/unstable")
		unstable.GET("/files/*filePath", viewFilesPath)
		unstable.PUT("/files/*filePath", createFile)
		unstable.DELETE("/files/*filePath", deleteFile)
	}

	return r
}

func main() {
	r := setupRouter()

	r.Run(":18080")
}
