package main

import (
	"fmt"
	"log"

	"github.com/cgascoig/intersight-simple-go/intersight"
	"github.com/icza/dyno"
)

func main() {
	// Setup a new Intersight client. Use the default settings from environment variables.
	client, err := intersight.NewClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Get all NTP Policy objects
	result, err := client.Get("/api/v1/ntp/Policies")
	if err != nil {
		log.Fatalf("Error in API call: %v", err)
	}

	// Extract the results list and loop through each NtpPolicy.
	// Here we use dyno to simplify extracting typed values from the interface{} the client returns
	results, err := dyno.GetSlice(result, "Results")
	if err != nil {
		log.Fatalf("Response does not contain 'Results' key: %v", err)
	}
	for _, result := range results {
		name, err := dyno.GetString(result, "Name")
		if err != nil {
			log.Fatalf("Response does not contain 'Name' key: %v", err)
		}
		enabled, err := dyno.GetBoolean(result, "Enabled")
		if err != nil {
			log.Fatalf("Response does not contain 'Enabled' key: %v", err)
		}
		ntpServers, err := dyno.GetSlice(result, "NtpServers")
		if err != nil {
			log.Fatalf("Response does not contain 'NtpServers' key: %v", err)
		}
		fmt.Printf("NTP Policy: Name=%s Enabled=%v NtpServers=%v\n", name, enabled, ntpServers)
	}
}
