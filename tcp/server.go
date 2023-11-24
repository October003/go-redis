package tcp

import (
	"context"
	"go-redis/interface/tcp"
	"go-redis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

// ListenAndServeWithSignal 
func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	signalChan := make(chan os.Signal,1)
	// SIGHUP     终止进程     终端线路挂断
	// SIGQUIT   建立CORE文件终止进程，并且生成core文件
	// SIGTERM   终止进程     软件终止信号
	// SIGINT     终止进程     中断进程
	// signal.Notufy 将监听到的系统信号转发到signalChan中
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		signal := <-signalChan
		switch signal {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()
	logger.Info("start listen...")
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	ListenAndServe(listener, handler, closeChan)
	return nil
}

// ListenAndServe
func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) error {
	ctx := context.Background()
	// 监听到关闭信号时关闭链接
	go func() {
		<-closeChan
		_ = listener.Close()
		_ = handler.Close()
	}()
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()
	var waitDone sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("link accepted...")
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
	return nil
}
