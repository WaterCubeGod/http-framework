package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/**
	server需要实现的功能
	1.启动
	2.关闭
	不同网站的协议不一致，有的是http，有的是https，不同的协议启动方式不一样，需要抽象server
**/
// server接口，内置http中的Handler
type server interface {
	//必须要组合http.Handler
	http.Handler
	//启动
	Start(adder string) error
	//关闭
	Stop() error
}

type HTTPOption func(h *HTTPServer)

//http实现
type HTTPServer struct {
	srv  *http.Server
	stop func() error
}

//优雅停机
func WithHTTPServerStop(fn func() error) HTTPOption {
	return func(h *HTTPServer) {
		if fn == nil {
			fn = func() error {
				fmt.Println("--------------------")
				quit := make(chan os.Signal, 1)
				signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
				<-quit
				log.Println("停止服务")

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := h.srv.Shutdown(ctx); err != nil {
					log.Println("服务关闭失败", err)
					return err
				}

				<-ctx.Done()
				log.Println("服务关闭成功")
				return nil
			}
		}
		h.stop = fn
	}
}

func NewHTTP(opts ...HTTPOption) server {
	h := &HTTPServer{}
	for _, opt := range opts {
		opt(h)
	}

	return h
}

//接收、转发请求
func (h *HTTPServer) ServeHTTP(http.ResponseWriter, *http.Request) {

}

//启动服务
func (h *HTTPServer) Start(adder string) error {
	h.srv = &http.Server{
		Addr:    adder,
		Handler: h,
	}
	return h.srv.ListenAndServe()
}

//停止服务
func (h *HTTPServer) Stop() error {
	return h.stop()
}

func main() {
	h := NewHTTP(WithHTTPServerStop(nil))
	go func() {
		err := h.Start(":8080")
		if err != nil && err != http.ErrServerClosed {
			panic("启动失败")
		}
	}()
	err := h.Stop()
	if err != nil {
		panic("关闭失败")
	}
}
