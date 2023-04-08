module github.com/poy/onpar

go 1.14

require (
	github.com/fatih/color v1.9.0
	github.com/nelsam/hel/v2 v2.3.2
)

// go modules make anything above v0 extremely hard to support, so we're
// scrapping any versions above v0. From now on, onpar will follow v0.x.x
// versioning.
//
// This retract directive retracts all releases beyond v0.
retract [v1.0.0, v1.1.3]
