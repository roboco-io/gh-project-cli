package models

import (
	"time"
)

// Project represents a GitHub Project v2
type Project struct {
	ID          string    `json:"id"`
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	URL         string    `json:"url"`
	Closed      bool      `json:"closed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Owner       Owner     `json:"owner"`
	Fields      []Field   `json:"fields,omitempty"`
	Items       []Item    `json:"items,omitempty"`
	Views       []View    `json:"views,omitempty"`
}

// Owner represents the owner of a project (user or organization)
type Owner struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Type  string `json:"type"` // "User" or "Organization"
	URL   string `json:"url"`
}

// CreateProjectInput represents input for creating a project
type CreateProjectInput struct {
	OwnerID     string  `json:"owner_id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
}

// UpdateProjectInput represents input for updating a project
type UpdateProjectInput struct {
	ID          string  `json:"id"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Closed      *bool   `json:"closed,omitempty"`
}

// ProjectListOptions represents options for listing projects
type ProjectListOptions struct {
	Owner  string   `json:"owner,omitempty"`
	States []string `json:"states,omitempty"` // "OPEN", "CLOSED"
	Limit  int      `json:"limit,omitempty"`
	After  string   `json:"after,omitempty"`
}
