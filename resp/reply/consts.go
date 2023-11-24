package reply

// ping 回复
type PongReply struct {
}

var pongBytes = []byte("+PONG\r\n")

func NewPongReply() *PongReply {
	return &PongReply{}
}

func (p PongReply) ToBytes() []byte {
	return pongBytes
}

// Ok 回复
type OkReply struct {
}

var okBytes = []byte("+OK\r\n")

var theOkReply = new(OkReply)

func NewOkReply() *OkReply {
	return theOkReply
}

func (o OkReply) ToBytes() []byte {
	return okBytes
}

// NullBulkReply 空字符串回复
type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n")

func NewNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

func (n NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

// EmptyMultiBulkReply 空数组回复
type EmptyMultiBulkReply struct {
}

var emptyMultiBulkBytes = []byte("*0\r\n")

func NewEmptyMutilBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

func (e EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

// NoReply 空回复
type NoReply struct {
}

var noBytes = []byte("")

func NewNoReply() *NoReply {
	return &NoReply{}
}

func (n NoReply) ToBytes() []byte {
	return noBytes
}
