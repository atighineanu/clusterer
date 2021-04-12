package data

var (
Inventorytmpl = `[{{.StackName}}]
{{.AllIps}}
[masters]
{{.IpMasters}}
[workers]
{{.IpWorkers}}`
)