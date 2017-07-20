package flag

import "github.com/giantswarm/microkit/flag"

type Flag struct {
	Excludes string
}

func New() *Flag {
	f := &Flag{}
	flag.Init(f)
	return f
}
