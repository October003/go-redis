package resp

// Redis Serializtion Protocol(RESP) redis序列化协议
// 正常回复   以“+”开头，以“\r\n”结尾的字符串形式  +OK\r\n
// 错误回复   以“-”开头，以“\r\n”结尾的字符串形式  -Error message\r\n
// 整数      以":"开头，以“\r\n”结尾的字符串形式  :123456\r\n
// 多行字符串 以"$"开头，后跟实际发送字节数，以"\r\n"结尾  $9\r\nimooc.com\r\n  $0\r\n
// 数组      以"*"开头，后跟成员个数  *3\r\n$3\r\nSET\r\n$3\r\nKEY\r\n$5\r\nvalue\r\n

type Connection interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
}
