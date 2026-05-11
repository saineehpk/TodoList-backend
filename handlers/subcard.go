package handlers

import (
	"net/http"
	"strconv"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func subCardOwnerCheck(c *gin.Context, subCardID uint) (*models.SubCard, bool) {
	var sub models.SubCard
	if err := database.DB.First(&sub, subCardID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SubCard not found"})
		return nil, false
	}
	if _, ok := cardOwnerCheck(c, sub.CardID); !ok {
		return nil, false
	}
	return &sub, true
}

func nextSubPosition(cardID uint) int {
	var count int64
	database.DB.Model(&models.SubCard{}).Where("card_id = ?", cardID).Count(&count)
	return int(count)
}

func CreateSubCard(c *gin.Context) {
	cardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}
	if _, ok := cardOwnerCheck(c, uint(cardID)); !ok {
		return
	}
	var input struct {
		Title string `json:"title" binding:"required,max=200"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub := models.SubCard{
		CardID:   uint(cardID),
		Title:    input.Title,
		Position: nextSubPosition(uint(cardID)),
	}
	if err := database.DB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subcard"})
		return
	}
	c.JSON(http.StatusCreated, sub)
}

func UpdateSubCard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subcard ID"})
		return
	}
	sub, ok := subCardOwnerCheck(c, uint(id))
	if !ok {
		return
	}
	var input struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]any{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Completed != nil {
		updates["completed"] = *input.Completed
	}
	database.DB.Model(sub).Updates(updates)
	database.DB.First(sub, sub.ID)
	c.JSON(http.StatusOK, sub)
}

func DeleteSubCard(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subcard ID"})
		return
	}
	sub, ok := subCardOwnerCheck(c, uint(id))
	if !ok {
		return
	}
	database.DB.Delete(sub)
	c.JSON(http.StatusOK, gin.H{"message": "SubCard deleted"})
}
