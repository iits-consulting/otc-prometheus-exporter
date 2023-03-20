package internal

type OtcNamespaces string

const (
	EcsNamespace = "SYS.ECS"
	RdsNamespace = "SYS.RDS"
	DmsNamespace = "SYS.DMS"
	NatNamespace = "SYS.NAT"
	ElbNamespace = "SYS.ELB"
)

var OtcNamespacesMapping = map[string]string{
	"ECS": EcsNamespace,
	"BMS": "SERVICE.BMS",
	// find more stuff
}
