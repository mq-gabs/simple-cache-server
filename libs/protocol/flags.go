package protocol

const FlagsSize = 1

type Flag byte

const (
	FlagSetNoResponse Flag = 1 << iota
	FlagSetCompressed
)
