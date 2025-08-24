package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemQueries(t *testing.T) {
	t.Run("GetIssue query structure", func(t *testing.T) {
		query := &GetIssueQuery{}
		assert.NotNil(t, query)
	})

	t.Run("GetPullRequest query structure", func(t *testing.T) {
		query := &GetPullRequestQuery{}
		assert.NotNil(t, query)
	})

	t.Run("SearchIssues query structure", func(t *testing.T) {
		query := &SearchIssuesQuery{}
		assert.NotNil(t, query)
	})

	t.Run("ListRepositoryIssues query structure", func(t *testing.T) {
		query := &ListRepositoryIssuesQuery{}
		assert.NotNil(t, query)
	})
}

func TestDraftIssueMutations(t *testing.T) {
	t.Run("CreateDraftIssue mutation structure", func(t *testing.T) {
		mutation := &CreateDraftIssueMutation{}
		assert.NotNil(t, mutation)
	})

	t.Run("UpdateDraftIssue mutation structure", func(t *testing.T) {
		mutation := &UpdateDraftIssueMutation{}
		assert.NotNil(t, mutation)
	})

	t.Run("DeleteDraftIssue mutation structure", func(t *testing.T) {
		mutation := &DeleteDraftIssueMutation{}
		assert.NotNil(t, mutation)
	})
}

func TestItemVariableBuilders(t *testing.T) {
	t.Run("BuildGetIssueVariables creates proper variables", func(t *testing.T) {
		variables := BuildGetIssueVariables("owner", "repo", 123)

		assert.NotNil(t, variables)
		assert.Equal(t, "owner", variables["owner"])
		assert.Equal(t, "repo", variables["repo"])
		assert.Equal(t, 123, variables["number"])
	})

	t.Run("BuildSearchIssuesVariables creates proper variables", func(t *testing.T) {
		opts := SearchOptions{
			Query: "is:issue is:open",
			First: 20,
		}

		variables := BuildSearchIssuesVariables(opts)

		assert.NotNil(t, variables)
		assert.Equal(t, "is:issue is:open", variables["query"])
		assert.Equal(t, 20, variables["first"])
		assert.NotContains(t, variables, "after")
	})

	t.Run("BuildSearchIssuesVariables with after cursor", func(t *testing.T) {
		after := "cursor123"
		opts := SearchOptions{
			Query: "is:issue is:open",
			First: 20,
			After: &after,
		}

		variables := BuildSearchIssuesVariables(opts)

		assert.NotNil(t, variables)
		assert.Equal(t, "cursor123", variables["after"])
	})

	t.Run("BuildListIssuesVariables creates proper variables", func(t *testing.T) {
		opts := ListIssueOptions{
			Owner:  "owner",
			Repo:   "repo",
			States: []string{"OPEN", "CLOSED"},
			First:  15,
		}

		variables := BuildListIssuesVariables(opts)

		assert.NotNil(t, variables)
		assert.Equal(t, "owner", variables["owner"])
		assert.Equal(t, "repo", variables["repo"])
		assert.Equal(t, 15, variables["first"])
		assert.Equal(t, []string{"OPEN", "CLOSED"}, variables["states"])
	})

	t.Run("BuildListIssuesVariables with default states", func(t *testing.T) {
		opts := ListIssueOptions{
			Owner: "owner",
			Repo:  "repo",
			First: 10,
		}

		variables := BuildListIssuesVariables(opts)

		assert.NotNil(t, variables)
		assert.Equal(t, []string{"OPEN"}, variables["states"])
	})

	t.Run("BuildCreateDraftIssueVariables creates proper variables", func(t *testing.T) {
		body := "Draft issue body"
		input := CreateDraftIssueInput{
			ProjectID: "project-id",
			Title:     "Draft Issue Title",
			Body:      &body,
		}

		variables := BuildCreateDraftIssueVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")

		inputVar := variables["input"].(map[string]interface{})
		assert.Equal(t, "project-id", inputVar["projectId"])
		assert.Equal(t, "Draft Issue Title", inputVar["title"])
		assert.Equal(t, "Draft issue body", inputVar["body"])
	})

	t.Run("BuildCreateDraftIssueVariables without body", func(t *testing.T) {
		input := CreateDraftIssueInput{
			ProjectID: "project-id",
			Title:     "Draft Issue Title",
		}

		variables := BuildCreateDraftIssueVariables(input)

		assert.NotNil(t, variables)
		inputVar := variables["input"].(map[string]interface{})
		assert.NotContains(t, inputVar, "body")
	})

	t.Run("BuildUpdateDraftIssueVariables creates proper variables", func(t *testing.T) {
		newTitle := "Updated Title"
		newBody := "Updated Body"
		input := UpdateDraftIssueInput{
			DraftIssueID: "draft-issue-id",
			Title:        &newTitle,
			Body:         &newBody,
		}

		variables := BuildUpdateDraftIssueVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")

		inputVar := variables["input"].(map[string]interface{})
		assert.Equal(t, "draft-issue-id", inputVar["draftIssueId"])
		assert.Equal(t, "Updated Title", inputVar["title"])
		assert.Equal(t, "Updated Body", inputVar["body"])
	})
}

func TestSearchOptions(t *testing.T) {
	t.Run("Default first value is set", func(t *testing.T) {
		opts := SearchOptions{
			Query: "test",
			First: 0, // Should be set to default
		}

		variables := BuildSearchIssuesVariables(opts)
		assert.Equal(t, 10, variables["first"])
	})
}

func TestListOptions(t *testing.T) {
	t.Run("Default first value is set for issues", func(t *testing.T) {
		opts := ListIssueOptions{
			Owner: "owner",
			Repo:  "repo",
			First: 0, // Should be set to default
		}

		variables := BuildListIssuesVariables(opts)
		assert.Equal(t, 10, variables["first"])
	})

	t.Run("Default first value is set for PRs", func(t *testing.T) {
		opts := ListPullRequestOptions{
			Owner: "owner",
			Repo:  "repo",
			First: 0, // Should be set to default
		}

		variables := BuildListPullRequestsVariables(opts)
		assert.Equal(t, 10, variables["first"])
	})
}
