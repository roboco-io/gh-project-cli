package models

import (
	"time"
)

// Item represents a project item (issue, PR, or draft)
type Item struct {
	ID        string            `json:"id"`
	Type      ItemType          `json:"type"`
	Title     string            `json:"title"`
	Body      *string           `json:"body,omitempty"`
	URL       *string           `json:"url,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Fields    map[string]string `json:"fields,omitempty"`
	Content   *ItemContent      `json:"content,omitempty"`
}

// ItemType represents the type of project item
type ItemType string

const (
	ItemTypeDraftIssue  ItemType = "DRAFT_ISSUE"
	ItemTypeIssue       ItemType = "ISSUE"
	ItemTypePullRequest ItemType = "PULL_REQUEST"
)

// ItemContent represents the content of an item (issue or PR)
type ItemContent struct {
	ID         string      `json:"id"`
	Number     int         `json:"number"`
	Title      string      `json:"title"`
	URL        string      `json:"url"`
	State      string      `json:"state"`
	Closed     bool        `json:"closed"`
	Repository *Repository `json:"repository,omitempty"`
}

// Repository represents a GitHub repository
type Repository struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    Owner  `json:"owner"`
	URL      string `json:"url"`
}

// CreateItemInput represents input for creating a project item
type CreateItemInput struct {
	ProjectID string   `json:"project_id"`
	Type      ItemType `json:"type"`
	Title     string   `json:"title,omitempty"`
	Body      *string  `json:"body,omitempty"`
	ContentID *string  `json:"content_id,omitempty"` // For existing issues/PRs
}

// UpdateItemInput represents input for updating a project item
type UpdateItemInput struct {
	ID     string            `json:"id"`
	Fields map[string]string `json:"fields,omitempty"`
}

// ItemListOptions represents options for listing items
type ItemListOptions struct {
	ProjectID string            `json:"project_id"`
	Filter    map[string]string `json:"filter,omitempty"`
	Limit     int               `json:"limit,omitempty"`
	After     string            `json:"after,omitempty"`
}
