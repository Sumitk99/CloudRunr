package models

type DeployReq struct {
	GitUrl     *string `json:"git_url"`
	ProjectID  *string `json:"project_id"`
	Framework  *string `json:"framework"`
	DistFolder *string `json:"dist_folder"`
}

type DeployRes struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Url    string `json:"url"`
}
