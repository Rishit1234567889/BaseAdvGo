package handlers

import (
	"database/sql"

	"github.com/Rishit1234567889/baseToAdvGo/internal/store"
)

type Handler struct { // 1.7(a)
	DB *sql.DB // DB instance // 3.5

	Queries *store.Queries // Query stores //3.5
}

func NewHandlers(db *sql.DB, queries *store.Queries) *Handler {
	return &Handler{
		DB:      db,
		Queries: queries,
	}

}
