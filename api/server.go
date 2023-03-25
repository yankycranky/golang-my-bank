package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/yankycranky/my-bank/db/sqlc"
)

type Server struct {
	router *gin.Engine
	store  *db.Store
}

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

func (server *Server) handleCreateAccount(c *gin.Context) {
	var createAccRequest CreateAccountRequest
	if err := c.ShouldBindJSON(&createAccRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := server.store.CreateAccount(c, db.CreateAccountParams{
		Owner:    createAccRequest.Owner,
		Currency: createAccRequest.Currency,
		Balance:  "0",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create account",
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":   "Account Created Successfully",
		"accountId": account.ID,
		"createdAt": account.CreatedAt,
	})
}

func NewServer(store *db.Store) *Server {

	server := &Server{
		store: store,
	}
	router := gin.Default()
	router.Handle(http.MethodPost, "/create-account", server.handleCreateAccount)
	server.router = router
	return server
}

func (server *Server) Start(address string) {
	server.router.Run(address)
}
