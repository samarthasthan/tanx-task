package handlers

import "github.com/samarthasthan/tanx-task/internal/database"

type Handler struct {
	database.Database
}

func NewHandler(db database.Database) *Handler {
	return &Handler{db}
}
