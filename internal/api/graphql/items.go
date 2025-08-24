package graphql

import "time"

// Repository represents a GitHub repository
type Repository struct {
	ID            string `graphql:"id"`
	Name          string `graphql:"name"`
	NameWithOwner string `graphql:"nameWithOwner"`
	Owner         struct {
		Login string `graphql:"login"`
	} `graphql:"owner"`
}

// Issue represents a GitHub issue
type Issue struct {
	ID     string `graphql:"id"`
	Title  string `graphql:"title"`
	Number int    `graphql:"number"`
	URL    string `graphql:"url"`
	State  string `graphql:"state"`
	Closed bool   `graphql:"closed"`

	Body      string    `graphql:"body"`
	CreatedAt time.Time `graphql:"createdAt"`
	UpdatedAt time.Time `graphql:"updatedAt"`

	Author struct {
		Login string `graphql:"login"`
	} `graphql:"author"`

	Repository Repository `graphql:"repository"`

	Labels struct {
		Nodes []struct {
			ID    string `graphql:"id"`
			Name  string `graphql:"name"`
			Color string `graphql:"color"`
		} `graphql:"nodes"`
	} `graphql:"labels(first: 10)"`

	Assignees struct {
		Nodes []struct {
			ID    string `graphql:"id"`
			Login string `graphql:"login"`
		} `graphql:"nodes"`
	} `graphql:"assignees(first: 10)"`
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	ID     string `graphql:"id"`
	Title  string `graphql:"title"`
	Number int    `graphql:"number"`
	URL    string `graphql:"url"`
	State  string `graphql:"state"`
	Closed bool   `graphql:"closed"`
	Merged bool   `graphql:"merged"`

	Body      string    `graphql:"body"`
	CreatedAt time.Time `graphql:"createdAt"`
	UpdatedAt time.Time `graphql:"updatedAt"`

	Author struct {
		Login string `graphql:"login"`
	} `graphql:"author"`

	Repository Repository `graphql:"repository"`

	Labels struct {
		Nodes []struct {
			ID    string `graphql:"id"`
			Name  string `graphql:"name"`
			Color string `graphql:"color"`
		} `graphql:"nodes"`
	} `graphql:"labels(first: 10)"`

	Assignees struct {
		Nodes []struct {
			ID    string `graphql:"id"`
			Login string `graphql:"login"`
		} `graphql:"nodes"`
	} `graphql:"assignees(first: 10)"`

	ReviewRequests struct {
		Nodes []struct {
			RequestedReviewer struct {
				User struct {
					ID    string `graphql:"id"`
					Login string `graphql:"login"`
				} `graphql:"... on User"`
			} `graphql:"requestedReviewer"`
		} `graphql:"nodes"`
	} `graphql:"reviewRequests(first: 10)"`
}

// DraftIssue represents a draft issue in a project
type DraftIssue struct {
	ID        string    `graphql:"id"`
	Title     string    `graphql:"title"`
	Body      *string   `graphql:"body"`
	CreatedAt time.Time `graphql:"createdAt"`
	UpdatedAt time.Time `graphql:"updatedAt"`

	Assignees struct {
		Nodes []struct {
			ID    string `graphql:"id"`
			Login string `graphql:"login"`
		} `graphql:"nodes"`
	} `graphql:"assignees(first: 10)"`
}

// Queries

// GetIssueQuery gets a specific issue
type GetIssueQuery struct {
	Repository struct {
		Issue Issue `graphql:"issue(number: $number)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

// GetPullRequestQuery gets a specific pull request
type GetPullRequestQuery struct {
	Repository struct {
		PullRequest PullRequest `graphql:"pullRequest(number: $number)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

// SearchIssuesQuery searches for issues
type SearchIssuesQuery struct {
	Search struct {
		IssueCount int `graphql:"issueCount"`
		Nodes      []struct {
			Issue Issue `graphql:"... on Issue"`
		} `graphql:"nodes"`
		PageInfo PageInfo `graphql:"pageInfo"`
	} `graphql:"search(query: $query, type: ISSUE, first: $first, after: $after)"`
}

// SearchPullRequestsQuery searches for pull requests
type SearchPullRequestsQuery struct {
	Search struct {
		IssueCount int `graphql:"issueCount"`
		Nodes      []struct {
			PullRequest PullRequest `graphql:"... on PullRequest"`
		} `graphql:"nodes"`
		PageInfo PageInfo `graphql:"pageInfo"`
	} `graphql:"search(query: $query, type: ISSUE, first: $first, after: $after)"`
}

// ListRepositoryIssuesQuery lists issues in a repository
type ListRepositoryIssuesQuery struct {
	Repository struct {
		Issues struct {
			Nodes    []Issue  `graphql:"nodes"`
			PageInfo PageInfo `graphql:"pageInfo"`
		} `graphql:"issues(first: $first, after: $after, states: $states, orderBy: {field: UPDATED_AT, direction: DESC})"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

// ListRepositoryPullRequestsQuery lists pull requests in a repository
type ListRepositoryPullRequestsQuery struct {
	Repository struct {
		PullRequests struct {
			Nodes    []PullRequest `graphql:"nodes"`
			PageInfo PageInfo      `graphql:"pageInfo"`
		} `graphql:"pullRequests(first: $first, after: $after, states: $states, orderBy: {field: UPDATED_AT, direction: DESC})"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

// Mutations

// CreateDraftIssueMutation creates a draft issue in a project
type CreateDraftIssueMutation struct {
	AddProjectV2DraftIssue struct {
		ProjectItem ProjectV2Item `graphql:"projectItem"`
	} `graphql:"addProjectV2DraftIssue(input: $input)"`
}

// UpdateDraftIssueMutation updates a draft issue
type UpdateDraftIssueMutation struct {
	UpdateProjectV2DraftIssue struct {
		DraftIssue DraftIssue `graphql:"draftIssue"`
	} `graphql:"updateProjectV2DraftIssue(input: $input)"`
}

// DeleteDraftIssueMutation deletes a draft issue
type DeleteDraftIssueMutation struct {
	DeleteProjectV2Item struct {
		DeletedItemID string `graphql:"deletedItemId"`
	} `graphql:"deleteProjectV2Item(input: $input)"`
}

// Input Types

// CreateDraftIssueInput represents input for creating a draft issue
type CreateDraftIssueInput struct {
	ProjectID string  `json:"projectId"`
	Title     string  `json:"title"`
	Body      *string `json:"body,omitempty"`
}

// UpdateDraftIssueInput represents input for updating a draft issue
type UpdateDraftIssueInput struct {
	DraftIssueID string  `json:"draftIssueId"`
	Title        *string `json:"title,omitempty"`
	Body         *string `json:"body,omitempty"`
}

// DeleteDraftIssueInput represents input for deleting a draft issue
type DeleteDraftIssueInput struct {
	ItemID string `json:"itemId"`
}

// SearchOptions represents search options for issues/PRs
type SearchOptions struct {
	Query string
	First int
	After *string
}

// ListIssueOptions represents options for listing repository issues
type ListIssueOptions struct {
	Owner  string
	Repo   string
	States []string
	First  int
	After  *string
}

// ListPullRequestOptions represents options for listing repository pull requests
type ListPullRequestOptions struct {
	Owner  string
	Repo   string
	States []string
	First  int
	After  *string
}

// Variable Builders

// BuildGetIssueVariables builds variables for getting an issue
func BuildGetIssueVariables(owner, repo string, number int) map[string]interface{} {
	return map[string]interface{}{
		"owner":  owner,
		"repo":   repo,
		"number": number,
	}
}

// BuildGetPullRequestVariables builds variables for getting a pull request
func BuildGetPullRequestVariables(owner, repo string, number int) map[string]interface{} {
	return map[string]interface{}{
		"owner":  owner,
		"repo":   repo,
		"number": number,
	}
}

// BuildSearchIssuesVariables builds variables for searching issues
func BuildSearchIssuesVariables(opts SearchOptions) map[string]interface{} {
	if opts.First <= 0 {
		opts.First = 10
	}

	variables := map[string]interface{}{
		"query": opts.Query,
		"first": opts.First,
	}

	if opts.After != nil {
		variables["after"] = *opts.After
	}

	return variables
}

// BuildSearchPullRequestsVariables builds variables for searching pull requests
func BuildSearchPullRequestsVariables(opts SearchOptions) map[string]interface{} {
	if opts.First <= 0 {
		opts.First = 10
	}

	variables := map[string]interface{}{
		"query": opts.Query,
		"first": opts.First,
	}

	if opts.After != nil {
		variables["after"] = *opts.After
	}

	return variables
}

// BuildListIssuesVariables builds variables for listing repository issues
func BuildListIssuesVariables(opts ListIssueOptions) map[string]interface{} {
	if opts.First <= 0 {
		opts.First = 10
	}

	variables := map[string]interface{}{
		"owner": opts.Owner,
		"repo":  opts.Repo,
		"first": opts.First,
	}

	if len(opts.States) > 0 {
		variables["states"] = opts.States
	} else {
		variables["states"] = []string{"OPEN"}
	}

	if opts.After != nil {
		variables["after"] = *opts.After
	}

	return variables
}

// BuildListPullRequestsVariables builds variables for listing repository pull requests
func BuildListPullRequestsVariables(opts ListPullRequestOptions) map[string]interface{} {
	if opts.First <= 0 {
		opts.First = 10
	}

	variables := map[string]interface{}{
		"owner": opts.Owner,
		"repo":  opts.Repo,
		"first": opts.First,
	}

	if len(opts.States) > 0 {
		variables["states"] = opts.States
	} else {
		variables["states"] = []string{"OPEN"}
	}

	if opts.After != nil {
		variables["after"] = *opts.After
	}

	return variables
}

// BuildCreateDraftIssueVariables builds variables for creating a draft issue
func BuildCreateDraftIssueVariables(input CreateDraftIssueInput) map[string]interface{} {
	inputMap := map[string]interface{}{
		"projectId": input.ProjectID,
		"title":     input.Title,
	}

	if input.Body != nil {
		inputMap["body"] = *input.Body
	}

	return map[string]interface{}{
		"input": inputMap,
	}
}

// BuildUpdateDraftIssueVariables builds variables for updating a draft issue
func BuildUpdateDraftIssueVariables(input UpdateDraftIssueInput) map[string]interface{} {
	inputMap := map[string]interface{}{
		"draftIssueId": input.DraftIssueID,
	}

	if input.Title != nil {
		inputMap["title"] = *input.Title
	}

	if input.Body != nil {
		inputMap["body"] = *input.Body
	}

	return map[string]interface{}{
		"input": inputMap,
	}
}

// BuildDeleteDraftIssueVariables builds variables for deleting a draft issue
func BuildDeleteDraftIssueVariables(input DeleteDraftIssueInput) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"itemId": input.ItemID,
		},
	}
}
