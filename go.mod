module github.com/poy/onpar/v2

go 1.18

require (
	git.sr.ht/~nelsam/hel/v4 v4.1.0
	github.com/fatih/color v1.9.0
)

require (
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.11 // indirect
	github.com/poy/onpar v1.1.2 // indirect
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
)

// go modules make anything above v0 extremely hard to support, so we're
// scrapping any versions above v0. From now on, onpar will follow v0.x.x
// versioning.
//
// This retract directive retracts all releases beyond v0.
retract [v2.0.0, v2.0.4]
