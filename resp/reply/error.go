package reply

// 未知错误
type UnknownErrReply struct {
}

var unknownErrBytes = []byte("-Err unknown\r\n")

func NewUnknownErrReply() *UnknownErrReply {
	return &UnknownErrReply{}
}

func (u UnknownErrReply) Error() string {
	return "Err Unknown"
}

func (u UnknownErrReply) ToBytes() []byte {
	return unknownErrBytes
}

// 参数个数错误
type ArgsNumErrReply struct {
	Cmd string // 指令
}

func NewArgsNumErrReply(cmd string) *ArgsNumErrReply {
	return &ArgsNumErrReply{
		Cmd: cmd,
	}
}

func (a *ArgsNumErrReply) Error() string {
	return "-Err wrong number of arguments for '" + a.Cmd + "' command\r\n"
}

func (a *ArgsNumErrReply) ToBytes() []byte {
	return []byte("-Err wrong number of arguments for '" + a.Cmd + "' command\r\n")
}

// 语法错误
type SyntaxErrReply struct {
}

var syntaxErrBytes = []byte("-Err syntax error\r\n")

var theSyntaxErrReply = new(SyntaxErrReply)

func NewSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErrReply
}

func (s *SyntaxErrReply) Error() string {
	return "Err Syntax error"
}

func (s *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

// 类型错误
type WrongTypeErrReply struct {
}

var wrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")

func NewWrongTypeErrReply() *WrongTypeErrReply {
	return &WrongTypeErrReply{}
}

func (w *WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value\r\n"
}

func (w *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

// 协议错误 (不符合RESP协议)
type ProtocolErrReply struct {
	Msg string
}

func NewProtocolErrReply() *ProtocolErrReply {
	return &ProtocolErrReply{}
}

func (p *ProtocolErrReply) Error() string {
	return "Err Protolcol error"
}

func (p *ProtocolErrReply) ToBytes() []byte {
	return []byte("-Err Protocol error: '" + p.Msg + "'\r\n")
}
