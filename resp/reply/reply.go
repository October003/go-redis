package reply

import (
	"bytes"
	"go-redis/interface/resp"
	"strconv"
)

var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n"
)

type BulkReply struct {
	Arg []byte
}

func NewBulkReply(arg []byte) *BulkReply {
	return &BulkReply{Arg: arg}
}

func (b BulkReply) ToBytes() []byte {
	if len(b.Arg) == 0 {
		return nullBulkReplyBytes
	}
	return []byte("$" + strconv.Itoa(len(b.Arg)) + CRLF + string(b.Arg) + CRLF)
}

type MultiBulkReply struct {
	Args [][]byte
}

func NewMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{Args: args}
}

func (m *MultiBulkReply) ToBytes() []byte {
	argLen := len(m.Args)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(argLen) + CRLF)
	for _, arg := range m.Args {
		if arg == nil {
			buf.WriteString(string(nullBulkReplyBytes) + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + string(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

// 回复状态
type StatusReply struct {
	Status string
}

func NewStatusReply(status string) *StatusReply {
	return &StatusReply{Status: status}
}

func (s *StatusReply) ToBytes() []byte {
	return []byte("+" + s.Status + CRLF)
}

type IntReply struct {
	Code int64
}

// 正常数字回复
func NewIntReply(code int64) *IntReply {
	return &IntReply{Code: code}
}

func (i *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(i.Code, 10) + CRLF)
}

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

type StandardErrReply struct {
	Status string
}

func (s *StandardErrReply) Error() string {
	return s.Status
}

func (s *StandardErrReply) ToBytes() []byte {
	return []byte("-" + s.Status + CRLF)
}

func isErrReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}
