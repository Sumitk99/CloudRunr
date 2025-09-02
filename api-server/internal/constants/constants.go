package constants

const (
	CONTAINER_IMAGE     = "cloudrunr-image"
	STATUS_FAILED       = "FAILED"
	STATUS_QUEUED       = "QUEUED"
	GITHUB_URL_PREFIX_1 = "https://github.com/"
	GITHUB_URL_PREFIX_2 = "github.com/"

	INVALID_GITHUB_URL_MESSAGE  = "Invalid Url, Please provide a valid github url"
	TOKEN_EXPIRED               = "Token is expired, Please login Again"
	INVALID_TOKEN               = "The token is no longer valid"
	NO_TOKEN                    = "No Authorization Token found. Please Log in"
	NO_PROJECT_FOUND            = "No Project found with this project id"
	UNAUTHORIZED_PROJECT_ACCESS = "You are not authorized to use this project"

	USER_NOT_FOUND = "No user found with this email"
	REACT          = "REACT"
	ANGULAR        = "ANGULAR"

	ROOT_DOMAIN = ".cloudrunr.micro-scale.software"
)

var VALID_FRAMEWORKS = []string{
	REACT, ANGULAR,
}
