package database

import (
	"go-redis/interface/resp"
	"go-redis/lib/wildcard"
	"go-redis/resp/reply"
)

func init() {
	RegisterCommand("del", ExecDel, -2)          // del k1 k2 k3  -2 为多个参数
	RegisterCommand("exists", ExecExists, -2)    // exists k1 k2 k3
	RegisterCommand("flushdb", ExecFlushDB, -1)  // flushdb ...   -1 为可变长参数
	RegisterCommand("type", ExecType, 2)         // type k1h
	RegisterCommand("rename", ExecRename, 3)     // rename k1 k2
	RegisterCommand("renamenx", ExecRenameNx, 3) //renamenx k1 k2
	RegisterCommand("keys", ExecKeys, 2)         // keys *
}

// DEL K1 K2 K2
func ExecDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleted := db.Removes(keys...)
	return reply.NewIntReply(int64(deleted))
}

// EXISTS K1 K2 K3
func ExecExists(db *DB, args [][]byte) resp.Reply {
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

// KEYS *
func ExecKeys(db *DB, args [][]byte) resp.Reply {
	pattern, _ := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val any) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return reply.NewMultiBulkReply(result)
}

// FLUSHDB
func ExecFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	return reply.NewOkReply()
}

// TYPE k1
func ExecType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exist := db.GetEntity(key)
	if !exist {
		return reply.NewStatusReply("none")
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.NewStatusReply("string")
	}
	// TODO:
	return reply.NewUnknownErrReply()
}

// RENAME k1 k2
func ExecRename(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])
	entity, exist := db.GetEntity(src)
	if !exist {
		return reply.NewStandardErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	_ = db.Remove(src)
	return reply.NewOkReply()
}

// RENAMENX k1 k2
func ExecRenameNx(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])
	_, ok := db.GetEntity(dest)
	if ok {
		return reply.NewIntReply(0)
	}
	entity, exist := db.GetEntity(src)
	if !exist {
		return reply.NewStandardErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	_ = db.Remove(src)
	return reply.NewIntReply(1)
}
