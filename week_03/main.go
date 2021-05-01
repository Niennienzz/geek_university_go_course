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

var (
	timer             *Timer
	server1           *Server
	server2           *Server
	interruptListener *InterruptListener
)

func init() {
	timer = NewTimer(time.Second * 10)

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

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

	if err := group.Wait(); err != nil {
		log.Printf("first error in the group: %v", err)
		return
	}
}

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
	go func() {
		<-ctx.Done()
		log.Printf("%s shutdown by context cancellation", s.label)
		s.server.Shutdown(context.Background())
	}()

	return s.server.ListenAndServe()
}

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
