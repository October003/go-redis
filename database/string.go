package database

import (
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/lib/utils"
	"go-redis/resp/reply"
)

func init() {
	RegisterCommand("get", ExecGet, 2)
	RegisterCommand("set", ExecSet, 3)
	RegisterCommand("setnx", ExecSetnx, 3)
	RegisterCommand("getset", ExecGetSet, 3)
	RegisterCommand("strlen", ExecStrlen, 2)
}

// GET k1
func ExecGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exist := db.GetEntity(key)
	if !exist {
		return reply.NewNullBulkReply()
	}
	value := entity.Data.([]byte)
	return reply.NewBulkReply(value)
}

// SET key value
func ExecSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := string(args[1])
	entity := &database.DataEntity{
		Data: value,
	}
	db.PutEntity(key, entity)
	db.addAof(utils.ToCmdLine3("set", args...))
	return reply.NewOkReply()
}

// SETNX k1 v1
func ExecSetnx(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := string(args[1])
	entity := &database.DataEntity{
		Data: value,
	}
	ret := db.PutIfAbsent(key, entity)
	db.addAof(utils.ToCmdLine3("setnx", args...))
	return reply.NewIntReply(int64(ret))
}

// GETSET k1 v1
func ExecGetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := string(args[1])
	entity, exist := db.GetEntity(key)
	db.PutEntity(key, &database.DataEntity{Data: value})
	if !exist {
		return reply.NewNullBulkReply()
	}
	ret := entity.Data.([]byte)
	db.addAof(utils.ToCmdLine3("getset", args...))
	return reply.NewBulkReply(ret)
}

// STRLEN k1
func ExecStrlen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exist := db.GetEntity(key)
	if !exist {
		return reply.NewNullBulkReply()
	}
	bytes := entity.Data.([]byte)
	return reply.NewIntReply(int64(len(bytes)))
}
