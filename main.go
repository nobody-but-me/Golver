
package main

import (
	"net/http"
	"context"
	"errors"
	"fmt"
	"net"
	"io"
)
const _key_server_address = "serveraddress"

func _get_root(_w http.ResponseWriter, _r *http.Request) {
	_ctx := _r.Context()
	fmt.Printf("[INFO from %s] GOT / REQUEST \n", _ctx.Value(_key_server_address))
	io.WriteString(_w, "[INFO] That's my WebSite \n")
}
func _get_hello(_w http.ResponseWriter, _r *http.Request) {
	_ctx := _r.Context()
	fmt.Printf("[INFO from %s] GOT / HELLO REQUEST \n", _ctx.Value(_key_server_address))
	io.WriteString(_w, "[INFO] Hello from Golang! \n")
}

func main() {
	_mux := http.NewServeMux()
	_mux.HandleFunc("/hello", _get_hello)
	_mux.HandleFunc("/", _get_root)
	
	_ctx, _cancel_ctx := context.WithCancel(context.Background())
	_server_one := &http.Server{
		Addr: ":4242",
		Handler: _mux,
		BaseContext: func(_l net.Listener) context.Context {
			_ctx = context.WithValue(_ctx, _key_server_address, _l.Addr().String())
			return _ctx
		},
	}
	_server_two := &http.Server{
		Addr: ":2424",
		Handler: _mux,
		BaseContext: func(_l net.Listener) context.Context {
			_ctx = context.WithValue(_ctx, _key_server_address, _l.Addr().String())
			return _ctx
		},
	}

	go func() {
		_err := _server_one.ListenAndServe()
		if errors.Is(_err, http.ErrServerClosed) {
			fmt.Printf("[INFO] Server one have been closed. \n")
		} else if _err != nil {
			fmt.Printf("[FAILED] Error server one: %s. \n", _err)
		}
		fmt.Printf("[INFO] Server one have been started. \n")
		_cancel_ctx()
	}()
	
	go func() {
		_err := _server_two.ListenAndServe()
		if errors.Is(_err, http.ErrServerClosed) {
			fmt.Printf("[INFO] Server two have been closed. \n")
		} else if _err != nil {
			fmt.Printf("[FAILED] Error server two: %s. \n", _err)
		}
		fmt.Printf("[INFO] Server two have been started. \n")
		_cancel_ctx()
	}()
	<-_ctx.Done()
}
