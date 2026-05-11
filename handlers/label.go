package handlers

import (
	"net/http"
	"strconv"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func ListLabels(c *gin.Context) {
	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}
	if _, ok := boardOwnerCheck(c, uint(boardID)); !ok {
		return
	}
	var labels []models.Label
	database.DB.Where("board_id = ?", boardID).Find(&labels)
	c.JSON(http.StatusOK, labels)
}

func CreateLabel(c *gin.Context) {
	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}
	if _, ok := boardOwnerCheck(c, uint(boardID)); !ok {
		return
	}
	var input struct {
		Name  string `json:"name" binding:"required"`
		Color string `json:"color" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	label := models.Label{BoardID: uint(boardID), Name: input.Name, Color: input.Color}
	if err := database.DB.Create(&label).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create label"})
		return
	}
	c.JSON(http.StatusCreated, label)
}

func UpdateLabel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID"})
		return
	}
	var label models.Label
	if err := database.DB.First(&label, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Label not found"})
		return
	}
	if _, ok := boardOwnerCheck(c, label.BoardID); !ok {
		return
	}
	var input struct {
		Name  *string `json:"name"`
		Color *string `json:"color"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Color != nil {
		updates["color"] = *input.Color
	}
	database.DB.Model(&label).Updates(updates)
	c.JSON(http.StatusOK, label)
}

func DeleteLabel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID"})
		return
	}
	var label models.Label
	if err := database.DB.First(&label, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Label not found"})
		return
	}
	if _, ok := boardOwnerCheck(c, label.BoardID); !ok {
		return
	}
	// Remove from all card associations then delete
	database.DB.Exec("DELETE FROM card_labels WHERE label_id = ?", label.ID)
	database.DB.Delete(&label)
	c.JSON(http.StatusOK, gin.H{"message": "Label deleted"})
}

func AttachLabel(c *gin.Context) {
	cardID, err1 := strconv.Atoi(c.Param("id"))
	labelID, err2 := strconv.Atoi(c.Param("labelId"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	card, ok := cardOwnerCheck(c, uint(cardID))
	if !ok {
		return
	}
	var label models.Label
	if err := database.DB.First(&label, labelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Label not found"})
		return
	}
	if label.BoardID != card.BoardID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Label does not belong to this board"})
		return
	}
	database.DB.Model(card).Association("Labels").Append(&label)
	c.JSON(http.StatusOK, gin.H{"message": "Label attached"})
}

func DetachLabel(c *gin.Context) {
	cardID, err1 := strconv.Atoi(c.Param("id"))
	labelID, err2 := strconv.Atoi(c.Param("labelId"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	card, ok := cardOwnerCheck(c, uint(cardID))
	if !ok {
		return
	}
	var label models.Label
	if err := database.DB.First(&label, labelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Label not found"})
		return
	}
	database.DB.Model(card).Association("Labels").Delete(&label)
	c.JSON(http.StatusOK, gin.H{"message": "Label detached"})
}
