package constants

const (
	CONTAINER_IMAGE     = "cloudrunr-image"
	STATUS_FAILED       = "FAILED"
	STATUS_QUEUED       = "QUEUED"
	GITHUB_URL_PREFIX_1 = "https://github.com/"
	GITHUB_URL_PREFIX_2 = "github.com/"

	INVALID_GITHUB_URL_MESSAGE = "Invalid Url, Please provide a valid github url"

	REACT   = "REACT"
	ANGULAR = "ANGULAR"
)

var VALID_FRAMEWORKS = []string{
	REACT, ANGULAR,
}
