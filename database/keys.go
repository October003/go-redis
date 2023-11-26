package database

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

func init() {
	RegisterCommand("del", execDel, -2)
	RegisterCommand("exists", execExists, -2)
	RegisterCommand("flushdb", execFlushDB, -1)
}

// DEL K1 K2 K2
func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleted := db.Removes(keys...)
	return reply.NewIntReply(int64(deleted))
}

// EXISTS K1 K2 K3
func execExists(db *DB, args [][]byte) resp.Reply {
	var result int64
	for _, arg := range args {
		key := string(arg)
		_, exists := db.GetEntity(key)
		if exists {
			result++
		}
	}
	return reply.NewIntReply(result)
}

// KEYS
func execKeys(db *DB, args [][]byte) resp.Reply {
	
}

// FLUSHDB
func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	return reply.NewOkReply()
}
