package main

import (
	"context"
	"fmt"
	"log"
	"os"

	intersight "github.com/CiscoDevNet/intersight-go"
)

func setupIntersightClient(keyID, keyFile string) (*intersight.APIClient, context.Context, error) {
	config := intersight.NewConfiguration()

	// Uncomment this line if you want to see the Intersight API requests/responses
	// config.Debug = true

	client := intersight.NewAPIClient(config)

	// Set up the authentication configuration struct
	authConfig := intersight.HttpSignatureAuth{
		KeyId:          keyID,
		PrivateKeyPath: keyFile,

		SigningScheme: intersight.HttpSigningSchemeRsaSha256,
		SignedHeaders: []string{
			intersight.HttpSignatureParameterRequestTarget, // The special (request-target) parameter expresses the HTTP request target.
			"Host",   // The Host request header specifies the domain name of the server, and optionally the TCP port number.
			"Date",   // The date and time at which the message was originated.
			"Digest", // A cryptographic digest of the request body.
		},
		SigningAlgorithm: intersight.HttpSigningAlgorithmRsaPKCS1v15,
	}

	// Get a context that includes the authentication config
	authCtx, err := authConfig.ContextWithValue(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("error creating authentication context: %v", err)
	}

	return client, authCtx, nil
}

func main() {
	keyID := os.Getenv("IS_KEY_ID")
	keyFile := os.Getenv("IS_KEY_FILE")

	intersightClient, intersightAuthCtx, err := setupIntersightClient(keyID, keyFile)
	if err != nil {
		log.Fatalf("Error setting up Intersight client: %v", err)
	}

	getNtpPolicyListRequest := intersightClient.NtpApi.GetNtpPolicyList(intersightAuthCtx)
	ntpResult, httpResult, err := getNtpPolicyListRequest.Execute()
	if err != nil {
		log.Fatalf("Error making Intersight API call: %v", err)
	}

	log.Print("HTTP Response: %v", httpResult)

	for _, ntpPolicy := range ntpResult.NtpPolicyList.Results {
		fmt.Printf("NTP Policy: Name=%v Enabled=%v NtpServers=%v\n", *ntpPolicy.Name, *ntpPolicy.Enabled, ntpPolicy.NtpServers)
	}
}
