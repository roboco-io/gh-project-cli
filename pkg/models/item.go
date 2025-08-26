package models

import (
	"time"
)

// Item represents a project item (issue, PR, or draft)
type Item struct {
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Body      *string           `json:"body,omitempty"`
	URL       *string           `json:"url,omitempty"`
	Fields    map[string]string `json:"fields,omitempty"`
	Content   *ItemContent      `json:"content,omitempty"`
	ID        string            `json:"id"`
	Type      ItemType          `json:"type"`
	Title     string            `json:"title"`
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
	Repository *Repository `json:"repository,omitempty"`
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	URL        string      `json:"url"`
	State      string      `json:"state"`
	Number     int         `json:"number"`
	Closed     bool        `json:"closed"`
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
	Body      *string  `json:"body,omitempty"`
	ContentID *string  `json:"content_id,omitempty"`
	ProjectID string   `json:"project_id"`
	Type      ItemType `json:"type"`
	Title     string   `json:"title,omitempty"`
}

// UpdateItemInput represents input for updating a project item
type UpdateItemInput struct {
	Fields map[string]string `json:"fields,omitempty"`
	ID     string            `json:"id"`
}

// ItemListOptions represents options for listing items
type ItemListOptions struct {
	Filter    map[string]string `json:"filter,omitempty"`
	ProjectID string            `json:"project_id"`
	After     string            `json:"after,omitempty"`
	Limit     int               `json:"limit,omitempty"`
}
