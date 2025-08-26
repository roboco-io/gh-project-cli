package models

import (
	"time"
)

// Field represents a project field
type Field struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      FieldType `json:"type"`
	Options   []string  `json:"options,omitempty"`
}

// FieldType represents the type of field
type FieldType string

const (
	FieldTypeText         FieldType = "TEXT"
	FieldTypeNumber       FieldType = "NUMBER"
	FieldTypeDate         FieldType = "DATE"
	FieldTypeSingleSelect FieldType = "SINGLE_SELECT"
	FieldTypeIteration    FieldType = "ITERATION"
)

// CreateFieldInput represents input for creating a field
type CreateFieldInput struct {
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	Type      FieldType `json:"type"`
	Options   []string  `json:"options,omitempty"`
}

// UpdateFieldInput represents input for updating a field
type UpdateFieldInput struct {
	ID      string   `json:"id"`
	Name    *string  `json:"name,omitempty"`
	Options []string `json:"options,omitempty"`
}
