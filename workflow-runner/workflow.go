package main

import (
	"context"
	"fmt"
	"time"

	intersight "github.com/CiscoDevNet/intersight-go"
)

type WorkflowClient struct {
	client *intersight.APIClient
	ctx    context.Context
}

func NewWorkflowClient(keyID, keyFile string) (*WorkflowClient, error) {

	config := intersight.NewConfiguration()

	// Uncomment this line if you want to see the Intersight API requests/responses
	config.Debug = true

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
		return nil, fmt.Errorf("error creating authentication context: %v", err)
	}

	return &WorkflowClient{client: client, ctx: authCtx}, nil
}

func (wfClient *WorkflowClient) GetWorkflowMoid(wfName string) (*intersight.WorkflowWorkflowDefinition, error) {
	res, _, err := wfClient.client.WorkflowApi.GetWorkflowWorkflowDefinitionList(wfClient.ctx).Filter(fmt.Sprintf("Name eq '%s'", wfName)).Orderby("Version desc").Top(1).Execute()
	if err != nil {
		return nil, fmt.Errorf("error making Intersight API call: %v", err)
	}

	if len(res.WorkflowWorkflowDefinitionList.Results) != 1 {
		return nil, fmt.Errorf("workflow query didn't return exactly 1 result (%d)", len(res.WorkflowWorkflowDefinitionList.Results))
	}

	return &res.WorkflowWorkflowDefinitionList.Results[0], nil
}

func (wfClient *WorkflowClient) ExecuteWorkflow(wf *intersight.WorkflowWorkflowDefinition, params map[string]interface{}) error {
	t := time.Now().Format("0102T030405")

	req := wfClient.client.WorkflowApi.CreateWorkflowWorkflowInfo(wfClient.ctx)
	req = req.WorkflowWorkflowInfo(intersight.WorkflowWorkflowInfo{
		ClassId:    "workflow.WorkflowInfo",
		ObjectType: "workflow.WorkflowInfo",
		Name:       optStr(fmt.Sprintf("workflow-runner-%s", t)),
		Organization: &intersight.OrganizationOrganizationRelationship{
			MoMoRef: &intersight.MoMoRef{
				ObjectType: "organization.Organization",
				Moid:       optStr("5ddec4226972652d33548943"),
			},
		},
		Action: optStr("Start"),
		Input: map[string]interface{}{
			"VMName": optStr(fmt.Sprintf("workflow-runner-%s", t)),
		},
		WorkflowDefinition: &intersight.WorkflowWorkflowDefinitionRelationship{
			MoMoRef: &intersight.MoMoRef{
				Moid:       optStr(wf.GetMoid()),
				ObjectType: "workflow.WorkflowDefinition",
			},
		},
		WorkflowCtx: *intersight.NewNullableWorkflowWorkflowCtx(&intersight.WorkflowWorkflowCtx{
			InitiatorCtx: *intersight.NewNullableWorkflowInitiatorContext(&intersight.WorkflowInitiatorContext{
				InitiatorMoid: optStr(wf.GetMoid()),
				InitiatorName: optStr(wf.GetName()),
				InitiatorType: optStr("workflow.WorkflowDefinition"),
			})}),
	})

	_, _, err := req.Execute()
	return err
}

func optStr(s string) *string {
	return &s
}
