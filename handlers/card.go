package handlers

import (
	"net/http"
	"strconv"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

var validColumns = map[string]bool{
	"backlog":    true,
	"todo":       true,
	"inprogress": true,
	"done":       true,
}

func boardOwnerCheck(c *gin.Context, boardID uint) (*models.Board, bool) {
	var board models.Board
	if err := database.DB.First(&board, boardID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Board not found"})
		return nil, false
	}
	if board.UserID != currentUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return nil, false
	}
	return &board, true
}

func cardOwnerCheck(c *gin.Context, cardID uint) (*models.Card, bool) {
	var card models.Card
	if err := database.DB.First(&card, cardID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return nil, false
	}
	if _, ok := boardOwnerCheck(c, card.BoardID); !ok {
		return nil, false
	}
	return &card, true
}

func nextPosition(boardID uint, column string) int {
	var count int64
	database.DB.Model(&models.Card{}).
		Where("board_id = ? AND column = ?", boardID, column).
		Count(&count)
	return int(count)
}

func compactPositions(boardID uint, column string) {
	var cards []models.Card
	database.DB.Where("board_id = ? AND column = ?", boardID, column).
		Order("position asc").Find(&cards)
	for i, card := range cards {
		if card.Position != i {
			database.DB.Model(&card).Update("position", i)
		}
	}
}

func CreateCard(c *gin.Context) {
	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}
	if _, ok := boardOwnerCheck(c, uint(boardID)); !ok {
		return
	}
	var input struct {
		Title       string  `json:"title" binding:"required,max=200"`
		Column      string  `json:"column" binding:"required"`
		Description string  `json:"description"`
		Priority    string  `json:"priority"`
		DueDate     *string `json:"dueDate"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !validColumns[input.Column] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid column"})
		return
	}
	card := models.Card{
		BoardID:     uint(boardID),
		Column:      input.Column,
		Position:    nextPosition(uint(boardID), input.Column),
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
	}
	if err := database.DB.Create(&card).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create card"})
		return
	}
	c.JSON(http.StatusCreated, card)
}

func UpdateCard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}
	card, ok := cardOwnerCheck(c, uint(id))
	if !ok {
		return
	}
	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Priority    *string `json:"priority"`
		DueDate     *string `json:"dueDate"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]any{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Priority != nil {
		updates["priority"] = *input.Priority
	}
	if input.DueDate != nil {
		updates["due_date"] = *input.DueDate
	}
	database.DB.Model(card).Updates(updates)
	database.DB.Preload("Labels").Preload("SubCards").First(card, card.ID)
	c.JSON(http.StatusOK, card)
}

func MoveCard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}
	card, ok := cardOwnerCheck(c, uint(id))
	if !ok {
		return
	}
	var input struct {
		Column   string `json:"column" binding:"required"`
		Position int    `json:"position"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !validColumns[input.Column] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid column"})
		return
	}

	oldColumn := card.Column

	// Shift cards in destination column to make room
	database.DB.Model(&models.Card{}).
		Where("board_id = ? AND column = ? AND position >= ? AND id != ?",
			card.BoardID, input.Column, input.Position, card.ID).
		UpdateColumn("position", clause.Expr{SQL: "position + 1"})

	database.DB.Model(card).Updates(map[string]any{
		"column":   input.Column,
		"position": input.Position,
	})

	// Compact old column if card moved to a different column
	if oldColumn != input.Column {
		compactPositions(card.BoardID, oldColumn)
	}

	c.JSON(http.StatusOK, card)
}

func DeleteCard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}
	card, ok := cardOwnerCheck(c, uint(id))
	if !ok {
		return
	}
	col := card.Column
	boardID := card.BoardID
	database.DB.Delete(card)
	compactPositions(boardID, col)
	c.JSON(http.StatusOK, gin.H{"message": "Card deleted"})
}
