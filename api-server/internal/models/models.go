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
}

type NewProjectReq struct {
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
	RunCommand      string `json:"run_command"`
	SubDomain       any    `json:"subdomain,omitempty"`
	CustomSubDomain any    `json:"custom_subdomain,omitempty"`
}

type NewProjectRes struct {
	DeploymentId string `json:"deployment_id"`
}

type LogData struct {
	LogStatement string    `json:"log_statement"`
	Time         time.Time `json:"time"`
}

type LogRetrievalResponse struct {
	Data []LogData `json:"data"`
}
