package protocol

type Command byte

const (
	CommandSize = 1

	CmdInvalid Command = 0x00

	CmdGet   Command = 0x01
	CmdSet   Command = 0x02
	CmdErase Command = 0x03
)
