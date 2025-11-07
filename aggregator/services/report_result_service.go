package services

import (
	"aggregator/entities"
	"encoding/json"
	"fmt"
	"log"
)

func ReportResultService(taskAttempt entities.TaskAttempt) error {
	//TODO: Implement the report result service

	// Pretty-print the TaskAttempt struct
	prettyJSON, err := json.MarshalIndent(taskAttempt, "", "  ")
	if err != nil {
		log.Printf("Error marshaling TaskAttempt: %v", err)
		return err
	}

	fmt.Println("ReportResultService received TaskAttempt:")
	fmt.Println(string(prettyJSON))

	return nil
}
