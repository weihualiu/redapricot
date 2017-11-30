package main

import (
	"flag"
	g "github.com/weihualiu/redapricot/cfg"
	"github.com/weihualiu/redapricot/socket"
	"github.com/weihualiu/redapricot/common/db"

	_ "github.com/weihualiu/redapricot/model"

	_ "net/http/pprof"
	"net/http"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	//version := flag.Bool("v", false, "show version")
	
	flag.Parse()

	//if *version {
	//	//fmt.Println(g.VERSION)
	//	os.Exit(0)
	//}

	g.ParseConfig(*cfg)

	db.Init()

	if g.Config().Debug {
		g.InitLog("debug")
	} else {
		g.InitLog("info")
	}

	go func() {
		http.ListenAndServe("0.0.0.0:8090", nil)
	}()

	go socket.Start()

	select {}
}

