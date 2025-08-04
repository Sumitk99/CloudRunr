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
)

func Script(srv *server.Server, projectID, framework, buildDestination string) {

	fmt.Print("Executing build-server script...\n")
	currPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	fmt.Printf("Current directory: %s\n", currPath)
	OutputDirPath := filepath.Join(currPath, "output")
	fmt.Printf("Output directory: %s\n", OutputDirPath)

	buildCommand := helper.DetectBuildCommand(framework)
	fullCommand := fmt.Sprintf("npm install && %s", buildCommand)
	log.Println("Build Command : ", buildCommand)
	process := exec.Command("bash", "-c", fullCommand)
	process.Dir = OutputDirPath
	stdout, err := process.StdoutPipe()
	if err != nil {
		fmt.Println("Error getting stdout pipe:", err)
		return
	}
	stderr, err := process.StderrPipe()
	if err != nil {
		fmt.Println("Error getting stderr pipe:", err)
		return
	}

	if err := process.Start(); err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	go func() {
		io.Copy(log.Writer(), stdout)
	}()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				fmt.Fprintf(os.Stderr, "ERROR: %s", buf[:n])
			}
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading stderr: %v", err)
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

	if len(buildDestination) == 0 {
		buildDestination = constants.DEFAULT_DIST_FOLDER
	}

	DistFolderPath := path.Join(OutputDirPath, buildDestination)

	files, err := helper.GetFilePaths(DistFolderPath)
	err = server.UploadToS3(srv.S3Client, DistFolderPath, projectID, files)
	if err != nil {
		log.Println("Error uploading files to S3:", err)
		return
	}

}
