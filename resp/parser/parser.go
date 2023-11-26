package parser

import (
	"bufio"
	"errors"
	"go-redis/interface/resp"
	"go-redis/lib/logger"
	"go-redis/resp/reply"
	"io"
	"runtime/debug"
	"strconv"
	"strings"
)

// Bulk String有两行，第一行为 $+正文长度，第二行为实际内容
// $3\r\nSET\r\n

type Payload struct {
	Data resp.Reply
	Err  error
}

type readState struct {
	readingMultiLine  bool
	expectedArgsCount int
	msgType           byte
	args              [][]byte
	bulkLen           int64 // 正文长度
}

func (s *readState) finished() bool {
	return s.expectedArgsCount > 0 && len(s.args) == s.expectedArgsCount
}

// ParseStream 解析字节流
func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

func parse0(reader io.Reader, ch chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	bufReader := bufio.NewReader(reader)
	var state readState
	var err error
	var msg []byte
	for {
		var ioErr bool
		msg, ioErr, err = readLine(bufReader, &state)
		if err != nil {
			if ioErr {
				ch <- &Payload{
					Err: err,
				}
				close(ch)
				return
			}
			ch <- &Payload{
				Err: err,
			}
			state = readState{}
			continue
		}
		// 判断是不是多行解析模式
		// 我们简单的将 Reply 分为两类:
		// 单行: StatusReply, IntReply, ErrorReply
		// 多行: BulkReply, MultiBulkRepl
		if !state.readingMultiLine {
			logger.Info("[parser parse0] state.readingMultiLine")
			if msg[0] == '*' {
				err := parseMultiBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: errors.New("protocol error:" + string(msg)),
					}
					state = readState{}
					continue
				}
				if state.expectedArgsCount == 0 {
					ch <- &Payload{
						Data: &reply.EmptyMultiBulkReply{},
					}
					state = readState{}
					continue
				}
			} else if msg[0] == '$' { // $4\r\nPING\r\n
				err := parseBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: errors.New("protocol error:" + string(msg)),
					}
					state = readState{}
					continue
				}
				if state.bulkLen == -1 {
					ch <- &Payload{
						Data: &reply.NullBulkReply{},
					}
					state = readState{}
					continue
				}
			} else {
				result, err := parseSingleLineReply(msg)
				ch <- &Payload{
					Data: result,
					Err:  err,
				}
				state = readState{}
				continue
			}
		} else {
			err := readBody(msg, &state)
			if err != nil {
				ch <- &Payload{
					Err: errors.New("protocol error:" + string(msg)),
				}
				state = readState{}
				continue
			}
			if state.finished() {
				var result resp.Reply
				if state.msgType == '*' {
					result = reply.NewMultiBulkReply(state.args)
				} else if state.msgType == '$' {
					result = reply.NewBulkReply(state.args[0])
				}
				ch <- &Payload{
					Data: result,
					Err:  err,
				}
				state = readState{}
			}
		}
	}
}

// readLine 读取一行数据
func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	var msg []byte
	var err error
	// 1.read simple line 直接\r\n切分
	if state.bulkLen == 0 {
		msg, err = bufReader.ReadBytes('\n')
		// logger.Info("[bulklen == 0 ]msg:" + string(msg))
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error:" + string(msg))
		}
	} else {
		// 2. read bulk line (binary safe)
		// 读取到 $5\r\nk\r\ney\r\n时 ，会将其误认为两行  应该读取指定长度的内容
		msg = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(bufReader, msg)
		// logger.Info("[bulklen != 0 ]msg:" + string(msg))
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-1] != '\n' || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error:" + string(msg))
		}
		state.bulkLen = 0
	}	
	logger.Infof("[parser readLine]state.bulkLen = %d , msg = %s", state.bulkLen, string(msg))
	return msg, false, nil
}

// parseMultiBulkHeader 初始化开始解析的过程
// *3\r\n$3\r\nSET\r\n$3\r\nKEY\r\n$5\r\nvalue\r\n
func parseMultiBulkHeader(msg []byte, state *readState) error {
	var err error
	var expectedLine uint64
	logger.Info("[parser parseMultiBulkHeader] msg = " + string(msg))
	expectedLine, err = strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 32)
	logger.Infof("[parser parseMultiBulkHeader] expectedLine = %d", expectedLine)
	if err != nil {
		return errors.New("protocol error:" + string(msg))
	}
	if expectedLine == 0 {
		state.expectedArgsCount = 0
		return nil
	} else if expectedLine > 0 {
		state.readingMultiLine = true
		state.expectedArgsCount = int(expectedLine)
		state.msgType = msg[0]
		state.args = make([][]byte, 0, expectedLine)
		return nil
	} else {
		return errors.New("protocol error:" + string(msg))
	}
}

// $4\r\nPING\r\n
// parseBulkHeader
func parseBulkHeader(msg []byte, state *readState) error {
	var err error
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 32)
	if err != nil {
		return errors.New("protocol error:" + string(msg))
	}
	if state.bulkLen == -1 {
		return nil
	} else if state.bulkLen > 0 {
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.expectedArgsCount = 1
		state.args = make([][]byte, 0, 1)
		return nil
	} else {
		return errors.New("protocol error:" + string(msg))
	}
}

// +OK\r\n	-Err\r\n	:5\r\n
func parseSingleLineReply(msg []byte) (resp.Reply, error) {
	str := strings.TrimSuffix(string(msg), "\r\n")
	var result resp.Reply
	switch msg[0] {
	case '+':
		result = reply.NewStatusReply(str[1:])
	case '-':
		result = reply.NewStandardErrReply(str[1:])
	case ':':
		val, err := strconv.ParseInt(str[1:], 10, 32)
		if err != nil {
			return nil, errors.New("protocol error:" + string(msg))
		}
		result = reply.NewIntReply(val)
	}
	return result, nil
}

// $4\r\n	PING\r\n
// *3\r\n	$3\r\nSET\r\n$3\r\nKEY\r\n$5\r\nvalue\r\n
func readBody(msg []byte, state *readState) error {
	var err error
	line := msg[0 : len(msg)-2]
	// $3
	if line[0] == '$' {
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return errors.New("protocol error:" + string(msg))
		}
		// $0\r\n
		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, []byte{})
	}
	return nil
}
