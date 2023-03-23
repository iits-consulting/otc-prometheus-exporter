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
	"AS" : "SYS.AS",
	"EVS" : "SYS.EVS",
	"SFS" : "SYS.SFS",
	"EFS" : "SYS.EFS",
	"CBR" : "SYS.CBR",
	"VPS" : "SYS.VPS",
	"ELB" : ElbNamespace,
	"NAT" : NatNamespace,
	"WAF" : "SYS.WAF",
	"DMS" : DmsNamespace,
	"DCS" : "SYS.DCS",
	"RDS" : RdsNamespace,
	"DDS" : "SYS.DDS",
	"NoSQL" : "SYS.NoSQL",
	"GAUSSDB" : "SYS.GAUSSDB",
	"GAUSSDBV5" : "SYS.GAUSSDBV5",
	"DWS" : "SYS.DWS",
	"ES" : "SYS.ES",

}
