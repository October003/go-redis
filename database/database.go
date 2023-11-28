package database

import (
	"go-redis/aof"
	"go-redis/config"
	"go-redis/interface/resp"
	"go-redis/lib/logger"
	"go-redis/resp/reply"
	"strconv"
	"strings"
)

type DataBase struct {
	dbSet      []*DB
	aofHandler *aof.AofHandler
}

func NewDataBase() *DataBase {
	database := &DataBase{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	database.dbSet = make([]*DB, config.Properties.Databases)
	for i := range database.dbSet {
		db := NewDB()
		db.index = i
		database.dbSet[i] = db
	}
	if config.Properties.AppendOnly {
		aofHandler, err := aof.NewAofHandler(database)
		if err != nil {
			panic(err)
		}
		database.aofHandler = aofHandler
		for _, db := range database.dbSet {
			db.addAof = func(line aof.CmdLine) {
				database.aofHandler.AddAof(db.index, line)
			}
		}
	}
	return database
}

type Database interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply
	Close()
	AfterClientClose(c resp.Connection)
}

// select x
func (db *DataBase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) != 2 {
			return reply.NewArgsNumErrReply("select")
		}
		ExecSelect(client, db, args[1:])
	}
	dbIndex := client.GetDBIndex()
	database := db.dbSet[dbIndex]
	return database.Exec(client, args)
}

func (db *DataBase) Close() {

}

func (db *DataBase) AfterClientClose(c resp.Connection) {

}

// select x
func ExecSelect(c resp.Connection, db *DataBase, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.NewStandardErrReply("ERR invalid db index")
	}
	if dbIndex > len(db.dbSet) {
		return reply.NewStandardErrReply("Err db index is out of range")
	}
	c.SelectDB(dbIndex)
	return reply.NewOkReply()
}
