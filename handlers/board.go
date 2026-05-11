package handlers

import (
	"net/http"
	"strconv"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func currentUserID(c *gin.Context) uint {
	user, _ := c.Get("user")
	m := user.(gin.H)
	return m["id"].(uint)
}

func ListBoards(c *gin.Context) {
	var boards []models.Board
	database.DB.Where("user_id = ?", currentUserID(c)).Find(&boards)
	c.JSON(http.StatusOK, boards)
}

func CreateBoard(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required,max=100"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	board := models.Board{UserID: currentUserID(c), Name: input.Name}
	if err := database.DB.Create(&board).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create board"})
		return
	}
	c.JSON(http.StatusCreated, board)
}

func GetBoard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}
	var board models.Board
	if err := database.DB.First(&board, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Board not found"})
		return
	}
	if board.UserID != currentUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	database.DB.
		Preload("Labels").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Order("position asc")
		}).
		Preload("Cards.Labels").
		First(&board, id)
	c.JSON(http.StatusOK, board)
}

func UpdateBoard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}
	var board models.Board
	if err := database.DB.First(&board, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Board not found"})
		return
	}
	if board.UserID != currentUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	var input struct {
		Name string `json:"name" binding:"required,max=100"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Model(&board).Update("name", input.Name)
	c.JSON(http.StatusOK, board)
}

func DeleteBoard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}
	var board models.Board
	if err := database.DB.First(&board, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Board not found"})
		return
	}
	if board.UserID != currentUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	database.DB.Delete(&board)
	c.JSON(http.StatusOK, gin.H{"message": "Board deleted"})
}
