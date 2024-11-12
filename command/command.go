package command

type Command interface {
	String() (string, error)
}

type CommandType int

const (
	AsteriskInfoCommand CommandType = iota + 1
	Ping
)

func (w CommandType) String() string {
	return [...]string{"asteriskinfo", "ping"}[w-1]
}
