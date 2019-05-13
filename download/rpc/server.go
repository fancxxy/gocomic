package rpc

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

// Download is defined for rpc server
type Download int

// Start rpc server
func Start() error {
	download := new(Download)
	rpc.Register(download)
	rpc.HandleHTTP()

	listen, err := net.Listen("tcp", ":7700")
	if err != nil {
		return err
	}

	fmt.Println("listen and serveing rpc on :7700")
	return http.Serve(listen, nil)
}
