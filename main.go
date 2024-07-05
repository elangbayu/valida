/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	// "log"
	"valida/cmd"
	// "valida/internal/apitest"
)

func main() {
	// err := apitest.InitializeLogFile("api_test_log.json")
	// if err != nil {
	// 	log.Fatalf("Failed to initialize log file: %v", err)
	// }
	// defer func() {
	// 	if err := apitest.CloseLogFile(); err != nil {
	// 		log.Printf("Failed to close log file: %v", err)
	// 	}
	// }()

	cmd.Execute()
}
