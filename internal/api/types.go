package api

import "time"

type Entry struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	Status    string    `json:"status"`
	Private   bool      `json:"private"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateEntryRequest struct {
	ProjectId int64  `json:"project_id"`
	TaskId    int64  `json:"task_id,omitempty"`
	Date      string `json:"spent_date,omitempty"`
	Hours     int    `json:"hours,omitempty"`
}

type UpdateEntryRequest struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
	Tags    *string `json:"tags,omitempty"`
	Status  *string `json:"status,omitempty"`
	Private *bool   `json:"private,omitempty"`
}

type ListParams struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Filter string `json:"filter,omitempty"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserAssignment struct {
	ID               int      `json:"id"`
	IsProjectManager bool     `json:"is_project_manager"`
	IsActive         bool     `json:"is_active"`
	Budget           *float64 `json:"budget"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
	HourlyRate       float64  `json:"hourly_rate"`
}

type CreateEntryResponse struct {
	ID             int            `json:"id"`
	SpentDate      string         `json:"spent_date"`
	User           User           `json:"user"`
	Client         Client         `json:"client"`
	Project        Project        `json:"project"`
	Task           Task           `json:"task"`
	UserAssignment UserAssignment `json:"user_assignment"`
	TaskAssignment TaskAssignment `json:"task_assignment"`
	Hours          float64        `json:"hours"`
	RoundedHours   float64        `json:"rounded_hours"`
	Notes          *string        `json:"notes"`
	CreatedAt      string         `json:"created_at"`
	UpdatedAt      string         `json:"updated_at"`
	IsLocked       bool           `json:"is_locked"`
	LockedReason   *string        `json:"locked_reason"`
	IsClosed       bool           `json:"is_closed"`
	ApprovalStatus string         `json:"approval_status"`
	IsBilled       bool           `json:"is_billed"`
	TimerStartedAt *string        `json:"timer_started_at"`
	StartedTime    *string        `json:"started_time"`
	EndedTime      *string        `json:"ended_time"`
	IsRunning      bool           `json:"is_running"`
	Invoice        *interface{}   `json:"invoice"`
	ExternalRef    *string        `json:"external_reference"`
	Billable       bool           `json:"billable"`
	Budgeted       bool           `json:"budgeted"`
	BillableRate   float64        `json:"billable_rate"`
	CostRate       float64        `json:"cost_rate"`
}

type ListEntriesResponse struct {
	Entries []*Entry `json:"entries"`
	Total   int      `json:"total"`
}

type ListAssignedProjectsResponse struct {
	ProjectAssignments []*ProjectAssignment `json:"project_assignments"`
	PerPage            int                  `json:"per_page"`
	TotalPages         int                  `json:"total_pages"`
	TotalEntries       int                  `json:"total_entries"`
	NextPage           *int                 `json:"next_page"`
	PreviousPage       *int                 `json:"previous_page"`
	Page               int                  `json:"page"`
	Links              Links                `json:"links"`
}

type ProjectAssignment struct {
	ID               int64             `json:"id"`
	IsProjectManager bool              `json:"is_project_manager"`
	IsActive         bool              `json:"is_active"`
	UseDefaultRates  bool              `json:"use_default_rates"`
	Budget           *float64          `json:"budget"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	HourlyRate       float64           `json:"hourly_rate"`
	Project          Project           `json:"project"`
	Client           ClientData        `json:"client"`
	TaskAssignments  []*TaskAssignment `json:"task_assignments"`
}

type Project struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type ClientData struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TaskAssignment struct {
	ID         int64     `json:"id"`
	Billable   bool      `json:"billable"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	HourlyRate float64   `json:"hourly_rate"`
	Budget     *float64  `json:"budget"`
	Task       Task      `json:"task"`
}

type Task struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Links struct {
	First    string  `json:"first"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Last     string  `json:"last"`
}
