package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	abm "github.com/romero96/golang-abm"
)

func main() {
	funcframework.RegisterHTTPFunction("/create", abm.CreateRecord)
	funcframework.RegisterHTTPFunction("/list", abm.ListRecords)
	funcframework.RegisterHTTPFunction("/update", abm.UpdateRecordByID)

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
