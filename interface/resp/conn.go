package resp

// Redis Serializtion Protocol(RESP) redis序列化协议
// 正常回复  +

// Connection Reids连接
type Connection interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
}


