package handlers

import (
	"net/http"
	"ticket-system/internal/database"
	"ticket-system/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateTicketRequest struct {
	Title  string `json:"title" binding:"required"`
	Detail string `json:"detail"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// Valid status transitions - assignment ka rule
var validTransitions = map[string][]string{
	"open":        {"in_progress"},
	"in_progress": {"closed"},
	"closed":      {}, // closed se kuch nahi ho sakta
}

func CreateTicket(c *gin.Context) {
	userID := c.GetString("userID")

	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket := models.Ticket{
		ID:     uuid.New().String(),
		Title:  req.Title,
		Detail: req.Detail,
		Status: "open",
		UserID: userID,
	}

	if result := database.DB.Create(&ticket); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ticket"})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func ListTickets(c *gin.Context) {
	userID := c.GetString("userID")

	var tickets []models.Ticket
	database.DB.Where("user_id = ?", userID).Find(&tickets)

	// Return empty array instead of null
	if tickets == nil {
		tickets = []models.Ticket{}
	}

	c.JSON(http.StatusOK, tickets)
}

func GetTicket(c *gin.Context) {
	userID := c.GetString("userID")
	ticketID := c.Param("id")

	var ticket models.Ticket
	if result := database.DB.Where("id = ? AND user_id = ?", ticketID, userID).First(&ticket); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func UpdateTicketStatus(c *gin.Context) {
	userID := c.GetString("userID")
	ticketID := c.Param("id")

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find ticket - ownership check bhi yahi ho raha hai
	var ticket models.Ticket
	if result := database.DB.Where("id = ? AND user_id = ?", ticketID, userID).First(&ticket); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	// Check if transition is valid
	allowed, exists := validTransitions[ticket.Status]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid current status"})
		return
	}

	isAllowed := false
	for _, s := range allowed {
		if s == req.Status {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid status transition from " + ticket.Status + " to " + req.Status,
		})
		return
	}

	// Update status
	database.DB.Model(&ticket).Update("status", req.Status)

	// Return updated ticket
	database.DB.Where("id = ?", ticketID).First(&ticket)
	c.JSON(http.StatusOK, ticket)
}
