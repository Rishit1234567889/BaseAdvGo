package handlers

import (
	"database/sql"

	"github.com/Rishit1234567889/baseToAdvGo/internal/store"
	"github.com/redis/go-redis/v9"
)

type Handler struct { // 1.7(a)
	DB      *sql.DB        // DB instance // 3.5
	Redis   *redis.Client  //7.1
	Queries *store.Queries // Query stores //3.5
}

func NewHandlers(db *sql.DB, queries *store.Queries, redisClient *redis.Client) *Handler {
	return &Handler{
		DB:      db,
		Queries: queries,
		Redis:   redisClient,
	}

}
