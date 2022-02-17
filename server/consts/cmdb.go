package consts

const (
	Outer = iota
	Inner
)

const (
	ArchitectureX86 = iota + 1
	ArchitectureArm
)

const (
	OsRedhat = iota + 1
	OsSuse
	OsCentos
	OsKylin
)

var ArchitectureMap = map[int64]string{
	ArchitectureX86: "x86",
	ArchitectureArm: "arm",
}

var ArchitectureMapReverse = map[string]int64{
	"x86": ArchitectureX86,
	"arm": ArchitectureArm,
}

var OsMap = map[int64]string{
	OsRedhat: "redhat",
	OsSuse:   "suse",
	OsCentos: "centos",
	OsKylin:  "kylin",
}

var OsMapReverse = map[string]int64{
	"redhat": OsRedhat,
	"suse":   OsSuse,
	"centos": OsCentos,
	"kylin":  OsKylin,
}
