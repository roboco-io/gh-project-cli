package models

import (
	"time"
)

// Project represents a GitHub Project v2
type Project struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description *string   `json:"description,omitempty"`
	Owner       Owner     `json:"owner"`
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Fields      []Field   `json:"fields,omitempty"`
	Items       []Item    `json:"items,omitempty"`
	Views       []View    `json:"views,omitempty"`
	Number      int       `json:"number"`
	Closed      bool      `json:"closed"`
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
	Description *string `json:"description,omitempty"`
	OwnerID     string  `json:"owner_id"`
	Title       string  `json:"title"`
}

// UpdateProjectInput represents input for updating a project
type UpdateProjectInput struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Closed      *bool   `json:"closed,omitempty"`
	ID          string  `json:"id"`
}

// ProjectListOptions represents options for listing projects
type ProjectListOptions struct {
	Owner  string   `json:"owner,omitempty"`
	After  string   `json:"after,omitempty"`
	States []string `json:"states,omitempty"`
	Limit  int      `json:"limit,omitempty"`
}
