package main

type buildkiteCreator struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type build struct {
	ID          string           `json:"id,omitempty"`
	URL         string           `json:"url,omitempty"`
	WebURL      string           `json:"web_url,omitempty"`
	Number      int              `json:"number,omitempty"`
	State       string           `json:"state,omitempty"`
	Blocked     bool             `json:"blocked,omitempty"`
	Message     string           `json:"message,omitempty"`
	Commit      string           `json:"commit,omitempty"`
	Branch      string           `json:"branch,omitempty"`
	Tag         *string          `json:"tag,omitempty"`
	Source      *string          `json:"source,omitempty"`
	Creator     buildkiteCreator `json:"creator,omitempty"`
	CreatedAt   *string          `json:"created_at,omitempty"`
	ScheduledAt *string          `json:"scheduled_at,omitempty"`
	StartedAt   *string          `json:"started_at,omitempty"`
	FinishedAt  *string          `json:"finished_at,omitempty"`
	RebuiltFrom string           `json:"rebuilt_from,omitempty"`
}

type pipeline struct {
	ID                   string  `json:"id,omitempty"`
	URL                  string  `json:"url,omitempty"`
	WebURL               string  `json:"web_url,omitempty"`
	Name                 string  `json:"name,omitempty"`
	Description          *string `json:"description,omitempty"`
	Slug                 string  `json:"slug,omitempty"`
	Repository           string  `json:"repository,omitempty"`
	ScheduledBuildsCount int     `json:"scheduled_builds_count,omitempty"`
	RunningBuildsCount   int     `json:"running_builds_count,omitempty"`
	ScheduledJobsCount   int     `json:"scheduled_jobs_count,omitempty"`
	RunningJobsCount     int     `json:"running_jobs_count,omitempty"`
	WaitingJobsCount     int     `json:"waiting_jobs_count,omitempty"`
	Visibility           string  `json:"visibility,omitempty"`
}

type job struct {
	ID           string  `json:"id,omitempty"`
	Type         string  `json:"type,omitempty"`
	Name         *string `json:"name,omitempty"`
	State        string  `json:"state,omitempty"`
	Command      *string `json:"command,omitempty"`
	WebURL       string  `json:"web_url,omitempty"`
	ExitStatus   *int    `json:"exit_status,omitempty"`
	CreatedAt    *string `json:"created_at,omitempty"`
	ScheduledAt  *string `json:"scheduled_at,omitempty"`
	StartedAt    *string `json:"started_at,omitempty"`
	FinishedAt   *string `json:"finished_at,omitempty"`
	Retried      *bool   `json:"retried,omitempty"`
	RetriesCount *int    `json:"retries_count,omitempty"`
	Agent        *agent  `json:"agent,omitempty"`
}

type agent struct {
	ID              string   `json:"id,omitempty"`
	WebURL          string   `json:"web_url,omitempty"`
	Name            string   `json:"name,omitempty"`
	ConnectionState string   `json:"connection_state,omitempty"`
	Version         string   `json:"version,omitempty"`
	Priority        int      `json:"priority,omitempty"`
	Metadata        []string `json:"meta_data,omitempty"`
}

type buildkiteEvent struct {
	Type     string    `json:"event,omitempty"`
	Build    *build    `json:"build,omitempty"`
	Pipeline *pipeline `json:"pipeline,omitempty"`
	Job      *job      `json:"job,omitempty"`
}
