package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const usage = "Usage: workflow-runner <workflow_name> [<args> ...]"

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		log.Fatal(usage)
	}

	workflowName := args[0]

	keyID := os.Getenv("IS_KEY_ID")
	keyFile := os.Getenv("IS_KEY_FILE")

	wfClient, err := NewWorkflowClient(keyID, keyFile)
	if err != nil {
		log.Fatalf("Error setting up Intersight client: %v", err)
	}

	wf, err := wfClient.GetWorkflowMoid(workflowName)
	if err != nil {
		log.Fatalf("Error finding workflow: %v", err)
	}

	fmt.Printf("Got wfMoid: %s\n", wf.GetMoid())

	err = wfClient.ExecuteWorkflow(wf, map[string]interface{}{})
	if err != nil {
		log.Fatalf("Error executing workflow: %v", err)
	}

	log.Print("Success")
}
