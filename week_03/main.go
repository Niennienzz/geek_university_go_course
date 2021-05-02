package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

// 将要在 errgroup.Group 里用到的四个成员
var (
	timer             *Timer
	server1           *Server
	server2           *Server
	interruptListener *InterruptListener
)

// 初始化上面四个成员
func init() {
	timer = NewTimer(time.Minute)

	{
		mux1 := http.NewServeMux()
		mux1.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello from server1"))
		})
		server1 = NewServer("server1", ":8080", mux1)
	}

	{
		mux2 := http.NewServeMux()
		mux2.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello from server2"))
		})
		server2 = NewServer("server2", ":8081", mux2)
	}

	interruptListener = NewInterruptListener()
}

func main() {
	// 给定一个全局上下文，若主函数退出则所有子上下文必须退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 建立一个 errgroup.Group，并把成员上下文绑定在一起
	group, ctx := errgroup.WithContext(ctx)

	// 利用 errgroup.Group 分别并发四个成员
	group.Go(func() error {
		return timer.Run(ctx)
	})

	group.Go(func() error {
		return server1.Run(ctx)
	})

	group.Go(func() error {
		return server2.Run(ctx)
	})

	group.Go(func() error {
		return interruptListener.Run(ctx)
	})

	// 等待所有成员返回
	// 若有成员出错导致返回，打印第一个出错的成员的错误
	if err := group.Wait(); err != nil {
		log.Printf("first error in the group: %v", err)
		return
	}
}

// 以下三个结构体的设计编写均满足以下原则
//  - 如果不是特殊情况，把并发控制权交给调用者
//  - 如果收到 ctx.Done 中的消息，则证明同一个 errgroup.Group 的其他成员出错了，让自己优雅关闭
//  - 反之，如果自己出错需要关闭，则通过返回错误来通知 errgroup.Group 里的其他成员

// Server 监听 ctx.Done
// 若 ctx 取消则自己也取消
// 若自己出错，则通过返回错误来通知其他成员
type Server struct {
	server http.Server
	label  string
}

func NewServer(label, addr string, handler http.Handler) *Server {
	return &Server{
		server: http.Server{
			Addr:    addr,
			Handler: handler,
		},
		label: label,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ch := make(chan error)

	// 这里是一个特殊情况，因为 ListenAndServe 是阻塞的
	// 所以在一个新的 goroutine 里面运行它，并监听它的错误
	go func() {
		ch <- s.server.ListenAndServe()
	}()

	// 若监听到 ctx.Done 取消，则关闭自己的服务
	// 然而关闭服务的过程也有可能比较耗时，并且有可能出错
	// 所以我们给 Shutdown 方法一个新的带有超时的上下文
	// 同时如果 Shutdown 方法出错，此时程序已经在退出的阶段了，我们通过简单输出日志的方式处理这个错误
	//
	// 若监听到自己的 ListenAndServe 出错，则通过返回错误来通知其他成员
	select {
	case <-ctx.Done():
		log.Printf("%s shutdown by context cancellation", s.label)
		ctxShutdown, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := s.server.Shutdown(ctxShutdown); err != nil {
			log.Printf("%s shutdown with error: %v", s.label, err)
		}
	case err := <-ch:
		return fmt.Errorf("%s shutdown with error: %w", s.label, err)
	}

	return nil
}

// Timer 监听 ctx.Done
// 若 ctx 取消则自己也取消
// 若定的时间到了，则自己通过返回错误来通知其他成员
type Timer struct {
	duration time.Duration
}

func NewTimer(duration time.Duration) *Timer {
	return &Timer{duration: duration}
}

func (t Timer) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		log.Printf("timer shutdown by context cancellation")
	case <-time.After(t.duration):
		return fmt.Errorf("timer shutdown after %v", t.duration)
	}
	return nil
}

// InterruptListener 监听 ctx.Done
// 若 ctx 取消则自己也取消
// 若收到系统中断命令 Ctrl+C，则自己通过返回错误来通知其他成员
type InterruptListener struct{}

func NewInterruptListener() *InterruptListener {
	return &InterruptListener{}
}

func (i InterruptListener) Run(ctx context.Context) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	select {
	case <-ctx.Done():
		log.Printf("interrupt listener shutdown by context cancellation")
	case <-ch:
		return errors.New("interrupt listener shutdown by os interrupt")
	}

	return nil
}
