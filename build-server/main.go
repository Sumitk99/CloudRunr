package main

import (
	"github.com/Sumitk99/build-server/script"
	"github.com/Sumitk99/build-server/server"
	"log"
)

func main() {

	s3, err := server.ConnectToS3("", "", "", "ap-south-1")
	if err != nil {
		log.Fatal(err)
	}
	Server := &server.Server{
		S3Client: s3,
	}
	script.Script(Server)
}
