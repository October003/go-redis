package tcp

import (
	"context"
	"net"
)

// Handler 业务逻辑的抽象接口
type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}
