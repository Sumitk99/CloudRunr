package script

import (
	"fmt"
	"github.com/Sumitk99/build-server/constants"
	"github.com/Sumitk99/build-server/helper"
	"github.com/Sumitk99/build-server/server"
	_ "github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

type BuildConfig struct {
	ProjectID    string
	DeploymentID string
	Framework    string
	BuildFolder  string
	RunCommand   string
}

func Script(srv *server.Server, cfg BuildConfig) {

	srv.LogProducer(constants.LOG_KAFKA_TOPIC, server.Log{
		DeploymentID: cfg.DeploymentID,
		ProjectID:    cfg.ProjectID,
		LogStatement: fmt.Sprintf("Build process started for projectId: %s", cfg.ProjectID),
		Timestamp:    time.Now(),
	})

	fmt.Print("Executing build-server script...\n")
	currPath, err := os.Getwd()
	if err != nil {
		srv.LogProducer(constants.LOG_KAFKA_TOPIC, server.Log{
			DeploymentID: cfg.DeploymentID,
			ProjectID:    cfg.ProjectID,
			LogStatement: fmt.Sprintf("ERROR: Failed to get current directory: %v", err),
			Timestamp:    time.Now(),
		})

		fmt.Println("Error getting current directory:", err)
		return
	}
	fmt.Printf("Current directory: %s\n", currPath)
	OutputDirPath := filepath.Join(currPath, "output")
	fmt.Printf("Output directory: %s\n", OutputDirPath)

	buildCommand := helper.DetectBuildCommand(cfg.Framework)
	fullCommand := fmt.Sprintf("npm install && %s", buildCommand)
	log.Println("Build Command : ", buildCommand)

	srv.LogProducer(constants.LOG_KAFKA_TOPIC, server.Log{
		DeploymentID: cfg.DeploymentID,
		ProjectID:    cfg.ProjectID,
		LogStatement: fmt.Sprintf("Detected build command: %s", buildCommand),
		Timestamp:    time.Now(),
	})

	process := exec.Command("bash", "-c", fullCommand)
	process.Dir = OutputDirPath
	stdout, err := process.StdoutPipe()
	if err != nil {
		fmt.Println("Error getting stdout pipe:", err)
		srv.LogProducer(constants.LOG_KAFKA_TOPIC, server.Log{
			DeploymentID: cfg.DeploymentID,
			ProjectID:    cfg.ProjectID,
			LogStatement: fmt.Sprintf("ERROR: Failed to get stdout pipe: %v", err),
			Timestamp:    time.Now(),
		})

		return
	}
	stderr, err := process.StderrPipe()
	if err != nil {
		fmt.Println("Error getting stderr pipe:", err)
		srv.LogProducer(constants.LOG_KAFKA_TOPIC, server.Log{
			DeploymentID: cfg.DeploymentID,
			ProjectID:    cfg.ProjectID,
			LogStatement: fmt.Sprintf("ERROR: Failed to get stderr pipe: %v", err),
			Timestamp:    time.Now(),
		})

		return
	}

	if err = process.Start(); err != nil {
		srv.LogProducer(constants.LOG_KAFKA_TOPIC, server.Log{
			DeploymentID: cfg.DeploymentID,
			ProjectID:    cfg.ProjectID,
			LogStatement: fmt.Sprintf("ERROR: Failed to start build command: %v", err),
			Timestamp:    time.Now(),
		})

		log.Fatalf("Failed to start command: %v", err)
	}
	srv.LogProducer(constants.LOG_KAFKA_TOPIC, server.Log{
		DeploymentID: cfg.DeploymentID,
		ProjectID:    cfg.ProjectID,
		LogStatement: fmt.Sprintf("Build command execution started"),
		Timestamp:    time.Now(),
	})

	go func() {
		io.Copy(log.Writer(), stdout)
	}()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				output := string(buf[:n])
				srv.LogProducer(constants.LOG_KAFKA_TOPIC, server.Log{
					DeploymentID: cfg.DeploymentID,
					ProjectID:    cfg.ProjectID,
					LogStatement: fmt.Sprintf("STDOUT: %s", output),
					Timestamp:    time.Now(),
				})

				fmt.Fprintf(os.Stderr, "ERROR: %s", buf[:n])
			}
			if err != nil {
				if err != io.EOF {
					srv.LogProducer(constants.LOG_KAFKA_TOPIC, server.Log{
						DeploymentID: cfg.DeploymentID,
						ProjectID:    cfg.ProjectID,
						LogStatement: fmt.Sprintf("ERROR: Reading stdout: %v", err),
						Timestamp:    time.Now(),
					})
				}
				break
			}
		}
	}()

	if err := process.Wait(); err != nil {
		log.Printf("ERROR: Build process failed: %v", err)
		return
	}
	fmt.Print("Build process completed successfully.\n")
	srv.LogProducer(constants.LOG_KAFKA_TOPIC, server.Log{
		DeploymentID: cfg.DeploymentID,
		ProjectID:    cfg.ProjectID,
		LogStatement: fmt.Sprintf("Build process completed successfully"),
		Timestamp:    time.Now(),
	})

	if len(cfg.BuildFolder) == 0 {
		cfg.BuildFolder = constants.DEFAULT_DIST_FOLDER
	}

	DistFolderPath := path.Join(OutputDirPath, cfg.BuildFolder)

	files, err := helper.GetFilePaths(DistFolderPath)
	err = srv.UploadToS3(DistFolderPath, cfg.ProjectID, files)
	if err != nil {
		log.Println("Error uploading files to S3:", err)
		return
	}

}
