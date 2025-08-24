package models

import (
	"time"
)

// View represents a project view
type View struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Layout    Layout    `json:"layout"`
	Filter    *string   `json:"filter,omitempty"`
	Sort      []Sort    `json:"sort,omitempty"`
	GroupBy   *string   `json:"group_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Layout represents the layout type of a view
type Layout string

const (
	LayoutTable   Layout = "TABLE_LAYOUT"
	LayoutBoard   Layout = "BOARD_LAYOUT"
	LayoutRoadmap Layout = "ROADMAP_LAYOUT"
)

// Sort represents a sort configuration
type Sort struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction"`
}

// SortDirection represents sort direction
type SortDirection string

const (
	SortDirectionAsc  SortDirection = "ASC"
	SortDirectionDesc SortDirection = "DESC"
)

// CreateViewInput represents input for creating a view
type CreateViewInput struct {
	ProjectID string  `json:"project_id"`
	Name      string  `json:"name"`
	Layout    Layout  `json:"layout"`
	Filter    *string `json:"filter,omitempty"`
}

// UpdateViewInput represents input for updating a view
type UpdateViewInput struct {
	ID      string  `json:"id"`
	Name    *string `json:"name,omitempty"`
	Filter  *string `json:"filter,omitempty"`
	Sort    []Sort  `json:"sort,omitempty"`
	GroupBy *string `json:"group_by,omitempty"`
}
