package database

import "go-redis/interface/resp"

type CmdLine = [][]byte

type Database interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply
	Close()
	AfterClientClose(c resp.Connection)
}

// Redis 数据结构
type DataEntity struct {
	Data interface{}
}
