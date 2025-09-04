package constants

const (
	BUCKET_NAME           = "cloudrunr"
	REACT                 = "REACT"
	REACT_BUILD_COMMAND   = "npm run build"
	ANGULAR_BUILD_COMMAND = "npx ng build --configuration=production"
	ANGULAR               = "ANGULAR"
	DEFAULT_DIST_FOLDER   = "dist"

	STATUS_DEPLOYED = "SUCCESS"
	STATUS_FAILED   = "FAILED"

	LOG_KAFKA_TOPIC          = "log_data"
	BUILD_STATUS_KAFKA_TOPIC = "build_status"
)
