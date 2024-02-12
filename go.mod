module github.com/poy/onpar

go 1.21.4

require (
	git.sr.ht/~nelsam/hel v0.6.6
	github.com/fatih/color v1.16.0
)

require (
	git.sr.ht/~nelsam/correct v0.0.5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/exp v0.0.0-20240205201215-2c58cdc269a3 // indirect
	golang.org/x/sys v0.17.0 // indirect
)

// go modules make anything above v0 extremely hard to support, so we're
// scrapping any versions above v0. From now on, hel will follow v0.x.x
// versioning.
//
// This retract directive retracts all releases beyond v0.
retract [v1.0.0, v1.1.3]
