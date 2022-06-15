package main

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/CiscoDevNet/intersight-go"
)

type IntersightAlarmPoller struct {
	client  *intersight.APIClient
	authCtx context.Context

	alarmLastCreationTime string
}

func NewIntersightAlarmPoller(keyID, keyFile string) (*IntersightAlarmPoller, error) {
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
		return nil, fmt.Errorf("error creating authentication context for Intersight: %v", err)
	}

	poller := &IntersightAlarmPoller{
		client:  client,
		authCtx: authCtx,
	}

	return poller, nil
}

func (p *IntersightAlarmPoller) PollAlarms() {
	req := p.client.CondApi.GetCondAlarmList(p.authCtx) //FaultApi.GetFaultInstanceList(ctx)
	req = req.Orderby("CreationTime desc")

	// Only get alarms since the time of the last one. If we don't know the time of the last one, just get the last 2.
	if p.alarmLastCreationTime != "" {
		req = req.Filter(fmt.Sprintf("(Severity in ('Critical','Warning','Info')) and (Acknowledge eq 'None') and (CreationTime gt %s)", p.alarmLastCreationTime))
	} else {
		req = req.Filter("(Severity in ('Critical','Warning','Info')) and (Acknowledge eq 'None')").Top(2)
	}

	res, httpRes, err := req.Execute()
	if err != nil {
		log.Errorf("Failed polling faults: %v, HTTP Response: %v: %v", err, httpRes.StatusCode, httpRes.Status)
		return
	}

	alarms := res.CondAlarmList.Results
	for j := len(alarms) - 1; j >= 0; j-- {
		alarm := alarms[j]
		log.Infof("Alarm retreived: %s", alarm.GetDescription())

		msg := fmt.Sprintf(
			`
## New Intersight Alarm

**Severity:** %s

**Affected Object:** %s(%s)

**Message:** %s: %s

**Creation Time:** %s

**Last Transition Time:** %s
`,
			alarm.GetSeverity(),
			alarm.GetAffectedMoDisplayName(),
			alarm.GetAffectedMoType(),
			alarm.GetName(),
			alarm.GetDescription(),
			timeToIntersightString(alarm.GetCreationTime()),
			timeToIntersightString(alarm.GetLastTransitionTime()))
		fmt.Printf("%s\n", msg)
	}

	if len(res.CondAlarmList.Results) > 0 {
		p.alarmLastCreationTime = timeToIntersightString(res.CondAlarmList.Results[0].GetCreationTime())
	}
}

func (p *IntersightAlarmPoller) Run() error {
	for {
		log.Info("Starting poll")
		p.PollAlarms()

		log.Info("Finished poll, sleeping 30 seconds")
		time.Sleep(30 * time.Second)
	}
}

func timeToIntersightString(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.000Z")
}
