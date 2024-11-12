package command

type Command interface {
	String() (string, error)
}

type CommandType int

const (
	Ping CommandType = iota + 1
)

func (w CommandType) String() string {
	return [...]string{"ping"}[w-1]
}
