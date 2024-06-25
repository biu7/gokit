package env

import "os"

const (
	StageKey = "STAGE"

	StageProd  = "prod"
	StageDebug = "debug"
	StageLocal = "local"
)

func Stage() string {
	return os.Getenv(StageKey)
}

func Prod() bool {
	return Stage() == StageProd
}

func Debug() bool {
	return Stage() == StageDebug
}

func Local() bool {
	var stage = Stage()
	return stage == "" || stage == StageLocal
}
