package protocol

type Flag byte

const (
	FlagSetNoResponse Flag = 1 << iota
	FlagSetCompressed
)
