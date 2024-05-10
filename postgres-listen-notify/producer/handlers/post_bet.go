package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func PostBetHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var betData map[string]interface{}
		if err := c.ShouldBindJSON(&betData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		user, _ := betData["username"].(string)
		team, _ := betData["team"].(string)
		amount, _ := betData["amount"].(float64)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := db.ExecContext(ctx, "INSERT INTO bets (timestamp, username, team, amount) VALUES ($1, $2, $3, $4)", time.Now().Format(time.RFC3339), user, team, amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to insert bet data: %s", err)})
			return
		}
		log.Printf("Inserted bet data of user %s for team %s with value %f", user, team, amount)
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Inserted bet data of user %s for team %s with value %f", user, team, amount)})
	}
}
