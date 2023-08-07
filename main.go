package main

import (
	"embed"
	"fmt"
	"github.com/freeconf/restconf"
	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/source"
	"net/http"
	"reflect"
)

//go:embed yang/*.yang
var Yang embed.FS

type Model struct {
	Democontainer DataProvider
}

type DataProvider struct {
	Nodeid NodeID
}

type NodeID string

func main() {
	model := Model{
		Democontainer: DataProvider{
			Nodeid: "initial",
		},
	}
	dev := device.New(source.EmbedDir(Yang, "yang"))

	// not using nodeutil.ReflectChild because I want to set members of Reflect in the real code
	rf := nodeutil.Reflect{}
	if err := dev.Add("demomodule", rf.Child(reflect.ValueOf(&model))); err != nil {
		fmt.Printf("cannot add demomodule module: %v\n", err)
		return
	}
	handler := restconf.NewHttpServe(dev)

	srv := http.Server{Addr: "127.0.0.1:8089", Handler: wrapper(&model, handler)}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("http server returned error: %v\n", err)
	}
}

func wrapper(model *Model, inner http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("before: %#v\n", *model)
		inner.ServeHTTP(w, r)
		fmt.Printf("after: %#v\n", *model)
	}
}
