package handler

import (
	"github.com/micro/go-micro/v2/server"
	pb "github.com/x-community/mail-service/proto"
)

// Options represents handler options
type Options struct {
	ServiceName string
	Config      Config
}

// Register will register service handler
func Register(s server.Server, opts Options) error {
	err := pb.RegisterMailServiceHandler(s, &service{id: opts.ServiceName, cfg: opts.Config})
	if err != nil {
		return err
	}
	return nil
}
