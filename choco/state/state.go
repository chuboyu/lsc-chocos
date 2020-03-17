package state

type State string

const (
	NIL      State = ""
	CREATED  State = "created"
	RUNNING  State = "running"
	STANDBY  State = "standby"
	DISABLED State = "disabled"
)
