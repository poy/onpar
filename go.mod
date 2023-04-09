module github.com/poy/onpar

go 1.18

require (
	git.sr.ht/~nelsam/hel v0.4.2
	github.com/fatih/color v1.15.0
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/poy/onpar/v3 v3.0.0-20230319125606-a35e84953e8f // indirect
	golang.org/x/sys v0.7.0 // indirect
)


// go modules make anything above v0 extremely hard to support, so we're
// scrapping any versions above v0. From now on, hel will follow v0.x.x
// versioning.
//
// This retract directive retracts all releases beyond v0.
retract [v1.0.0, v1.1.3]
