// Package api provide the server side Goed API
// via RPC over local socket.
// See client/ for the client implementation.
package api

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

type Api struct {
}

func (a *Api) Start() {
	r := new(GoedRpc)
	rpc.Register(r)
	rpc.HandleHTTP()
	l, err := net.Listen("unix", core.Socket)
	if err != nil {
		log.Fatalf("Socket listen error %s : \n", err.Error())
	}

	go func() {
		err = http.Serve(l, nil)
		if err != nil {
			panic(err)
		}
	}()
}

// Goed RPC functions holder
type GoedRpc struct{}

type RpcStruct struct {
	Data []string
}

func (r *GoedRpc) Action(args RpcStruct, res *RpcStruct) error {
	results, err := actions.Exec(args.Data[0], args.Data[1:])
	for _, r := range results {
		res.Data = append(res.Data, r)
	}
	return err
}

func (r *GoedRpc) Open(args []interface{}, _ *struct{}) error {
	vid := actions.Ar.EdOpen(args[1].(string), -1, args[0].(string), true)
	actions.Ar.EdActivateView(vid)
	actions.Ar.EdRender()
	return nil
}

func (r *GoedRpc) Edit(args []interface{}, _ *struct{}) error {
	prevView := actions.Ar.EdCurView()
	vid := actions.Ar.EdOpen(args[1].(string), -1, args[0].(string), true)
	actions.Ar.EdActivateView(vid)
	actions.Ar.EdRender()
	// Wait til file closed
	for {
		idx, _ := actions.Ar.EdViewIndex(vid)
		if idx == -1 { // view no longer present
			// switch back to the original view
			actions.Ar.EdActivateView(prevView)
			actions.Ar.EdRender()
			return nil
		}
		time.Sleep(250 * time.Millisecond)
	}
}
