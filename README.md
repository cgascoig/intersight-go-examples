# Intersight Go Examples

This repository contains example code to show how the [Intersight API](https://intersight.com/apidocs/introduction/overview/) can be used with the [Intersight Go SDK](https://github.com/CiscoDevNet/intersight-go). 

## Examples

All the examples below will look for the IS_KEY_ID and IS_KEY_FILE environment variables for the Intersight API Key ID and file name for Intersight API Key (Secret).

```
export IS_KEY_ID=123456789012345678901234/123456789012345678901234/123456789012345678901234
export IS_KEY_FILE=~/intersight-api-key.pem
```

You can build all the examples using the Makefile:
```
gmake
```

Or build them manually:

```
go build -o "build/list-ntp-policies" ./list-ntp-policies
go build -o "build/alarm-streamer" ./alarm-streamer
go build -o "build/workflow-runner" ./workflow-runner
```

### List NTP Policies

This is the simplest example, it simply gets all the NTP Policy managed objects from the Intersight API and prints out their name, enabled state and NtpServer list.

```
$ ./build/list-ntp-policies
NTP Policy: Name=NTP_ESL Enabled=true NtpServers=[ntp.esl.cisco.com]
NTP Policy: Name=CiscoNTP Enabled=true NtpServers=[ntp.esl.cisco.com]
NTP Policy: Name=BA-NTP Enabled=true NtpServers=[ntp.esl.cisco.com]
NTP Policy: Name=se-cimc-6ntp-policy Enabled=false NtpServers=[]
```

### Alarm Streamer

This example will poll the Intersight cond.Alarm API every 30s and display in new alarms.

```
$ ./build/alarm-streamer
INFO[0000] Starting poll
INFO[0000] Alarm retreived: Capacity analysis has found that 3 volumes on your cluster have potentially 6.131 TB of inactive data.

## New Intersight Alarm

**Severity:** Warning

**Affected Object:** dc-netapp-aiq.nsd5.ciscolabs.com/Cluster-nsd5_netapp()

**Message:** StorageNetAppClusterWarningEvent: Capacity analysis has found that 3 volumes on your cluster have potentially 6.131 TB of inactive data.

**Creation Time:** 2022-06-13T11:11:20.084Z

**Last Transition Time:** 2022-06-13T11:11:20.084Z

INFO[0000] Finished poll, sleeping 30 seconds
```

### Workflow Runner

This example shows how to start an Intersight Cloud Orchestrator workflow. 

```
$ ./build/workflow-runner NewVMCmd
Got workflow Moid: 62a95e18696f6e2d318aa1ab
Workflow successfully started ...

```



