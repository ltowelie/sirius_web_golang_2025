package build

import (
	"runtime"
	"sync"
)

var (
	once sync.Once
)

var (
	Version    string
	CommitHash string
	BuildDate  string
)

type Info struct {
	Version    string `json:"version,omitempty"`
	CommitHash string `json:"commit_hash,omitempty"`
	BuildDate  string `json:"build_date,omitempty"`
	GoVersion  string `json:"go_version,omitempty"`
}

func New() *Info {
	var info Info

	once.Do(func() {
		info = Info{
			Version:    Version,
			CommitHash: CommitHash,
			BuildDate:  BuildDate,
			GoVersion:  runtime.Version(),
		}
	})

	return &info
}
