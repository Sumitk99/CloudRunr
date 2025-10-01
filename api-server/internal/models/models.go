package models

import "time"

type DeployReq struct {
	ProjectID string `json:"project_id"`
}

type DeployRes struct {
	Status       string `json:"status"`
	DeploymentID string `json:"deployment_id,omitempty"`
	Error        string `json:"error,omitempty"`
	Url          string `json:"url,omitempty"`
}

type NewDeployment struct {
	DeploymentID *string `json:"deployment_id"`
	GitUrl       *string `json:"git_url"`
	Framework    *string `json:"framework"`
	DistFolder   *string `json:"dist_folder"`
	ProjectID    *string `json:"project_id"`
	RunCommand   *string `json:"run_command"`
	RootFolder   *string `json:"root_folder"`
}

type NewProjectReq struct {
	RootFolder string  `json:"root"`
	GitUrl     *string `json:"git_url"`
	Framework  *string `json:"framework"`
	DistFolder *string `json:"dist_folder"`
	ProjectID  *string `json:"project_id"`
	Name       *string `json:"name"`
	RunCommand *string `json:"run_command"`
}

type ProjectDetails struct {
	UserID          string `json:"user_id"`
	GitUrl          string `json:"git_url"`
	Framework       string `json:"framework"`
	DistFolder      string `json:"dist_folder"`
	ProjectID       string `json:"project_id"`
	Name            string `json:"name"`
	RootFolder      string `json:"root"`
	RunCommand      string `json:"run_command"`
	SubDomain       any    `json:"subdomain,omitempty"`
	CustomSubDomain any    `json:"custom_subdomain,omitempty"`
}

type NewProjectRes struct {
	DeploymentId string `json:"deployment_id"`
}

type UserProjectListContent struct {
	ProjectID string `json:"project_id"`
	Name      string `json:"name"`
	GitUrl    string `json:"git_url"`
	Framework string `json:"framework"`
}

type LogData struct {
	ID           int64     `json:"id,omitempty"`
	LogStatement string    `json:"log_statement"`
	Time         time.Time `json:"time"`
}

type LogRetrievalResponse struct {
	Data       []LogData `json:"data"`
	HasMore    bool      `json:"has_more"`
	NextCursor *int64    `json:"next_cursor,omitempty"`
	Status     string    `json:"status"`
}

type DeploymentDetails struct {
	DeploymentID string `json:"deployment_id"`
	ProjectID    string `json:"project_id"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}

type DeploymentListResponse struct {
	Deployments []DeploymentDetails `json:"deployments"`
}
