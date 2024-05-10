package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lib/pq"
)

const channelName = "new_bet"

type dbNotification struct {
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
	Team      string    `json:"team"`
	Amount    float64   `json:"amount"`
}

type Listener struct {
	db          *sql.DB
	channelName string
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := connectDB(ctx)
	if err != nil {
		panic(err)
	}

	listener := NewListener(db)
	err = listener.Start(ctx)
	if err != nil {
		panic(err)
	}
}

func connectDB(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set up connection pool parameters
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func NewListener(db *sql.DB) *Listener {
	return &Listener{
		db:          db,
		channelName: channelName,
	}
}

func (l *Listener) Start(ctx context.Context) error {
	slog.Info("Starting listener")

	listener := pq.NewListener(os.Getenv("DATABASE_URL"), 10*time.Second, time.Minute, nil)
	err := listener.Listen(channelName)
	if err != nil {
		return fmt.Errorf("failed to listen to channel: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			notification := listener.NotificationChannel()
			for n := range notification {
				l.handleNotifications(n)
			}
		}
	}
}

func (d *Listener) handleNotifications(notification *pq.Notification) {
	payload, err := d.parsePayload(notification.Extra)
	if err != nil {
		slog.With("error", err).Error("Could not parse PostgreSQL notification payload.")
		return
	}
	logCallback(payload)
}

func logCallback(payload dbNotification) {
	slog.Info(fmt.Sprintf("notification received: Timestamp: %s, UserName: %s, Team: %s, Bet Amount: %f", payload.Timestamp, payload.Username, payload.Team, payload.Amount))
}

func (d *Listener) parsePayload(rawPayload string) (dbNotification, error) {
	var payload dbNotification

	err := json.Unmarshal([]byte(rawPayload), &payload)
	if err != nil {
		return payload, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return payload, nil
}
