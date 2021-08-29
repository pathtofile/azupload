# AzUpload

Simple project that takes arbitrary JSON from stdin (1 per line) and uploads it to
Azure Monitor. This can be used to pipe JSON for a programs stdout into Azure Monitor, without needing
to 'cache' it to disk.

I don't think the Microsoft-supplied tools can easily read from stdin, however this was mostly
an execuse to learn Go.

I'm using this to send [Tracee events](https://github.com/aquasecurity/tracee/tree/main/tracee-ebpf) to Azure.

# Setup
In Azure, first create a Log Analytics workspace.
Then follow the steps in these articles to set the required environment variables:
- https://docs.microsoft.com/en-us/azure/azure-monitor/logs/data-collector-api
- https://docs.microsoft.com/en-us/azure/azure-arc/data/upload-logs

tl;dr:
1. In the Azure portal, locate your Log Analytics workspace.
2. Select Agents management.
3. Copy Workspace ID, and set as `WORKSPACE_ID`
4. Copy Primary Key, and set as `WORKSPACE_SHARED_KEY`

# Build
Either grab the latest pre-build from [the releases page](https://github.com/pathtofile/azupload/releases)
or build it yourself:
```bash
git checkout https://github.com/pathtofile/azupload.git
cd azupload
go build
```

# Run
```bash
# Set environment variables
$> export WORKSPACE_ID="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
$> export WORKSPACE_SHARED_KEY="aaaaaaaaaaaaa=="

# Get a program that outputs JSON events 1 per line
# If data is in a file, just use cat
$> cat events.json
{"timestamp":1630205061359686618,"processId":9778}
{"timestamp":1630205061359686618,"processId":9779}

# Then just pipe output into azupload
$> cat events.json | ./azupload
POST Successful at Sun, 29 Aug 2021 08:53:45 GMT
POST Successful at Sun, 29 Aug 2021 08:53:45 GMT
```
