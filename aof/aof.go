package aof

import (
	"go-redis/config"
	"go-redis/interface/database"
	"go-redis/lib/logger"
	"go-redis/lib/utils"
	"go-redis/resp/reply"
	"os"
	"strconv"
)

// append only file

type CmdLine [][]byte

const aofBufferSize = 1 << 16

type payload struct {
	cmdLine CmdLine
	dbIndex int
}

type AofHandler struct {
	database    database.Database
	aofChan     chan *payload
	aofFile     *os.File
	aofFileName string
	currentDB   int
}

// NewAofHandler
func NewAofHandler(db database.Database) (*AofHandler, error) {
	handler := &AofHandler{}
	handler.database = db
	handler.aofFileName = config.Properties.AppendFilename
	// LOAD AOF
	handler.LoadAof()
	aofile, err := os.OpenFile(handler.aofFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofFile = aofile
	handler.aofChan = make(chan *payload, aofBufferSize)
	go func() {
		handler.HandleAof()
	}()
	return handler, nil
}

// Add payload(set k v) --> aofChan
func (handler *AofHandler) AddAof(dbIndex int, cmd CmdLine) {
	if config.Properties.AppendOnly && handler.aofChan != nil {
		handler.aofChan <- &payload{
			cmdLine: cmd,
			dbIndex: dbIndex,
		}
	}
}

// HandleAof payload(set k v) <-- aofChan (落盘)
func (handler *AofHandler) HandleAof() {
	handler.currentDB = 0
	for p := range handler.aofChan {
		if p.dbIndex != handler.currentDB {
			// 添加select
			data := reply.NewMultiBulkReply(utils.ToCmdLine("select" + strconv.Itoa(p.dbIndex))).ToBytes()
			_, err := handler.aofFile.Write(data)
			if err != nil {
				logger.Error(err)
				continue
			}
			handler.currentDB = p.dbIndex
		}
		data := reply.NewMultiBulkReply(p.cmdLine).ToBytes()
		_, err := handler.aofFile.Write(data)
		if err != nil {
			logger.Error(err)
		}
	}
}

func (handler *AofHandler) LoadAof() {

}
