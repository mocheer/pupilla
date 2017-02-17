package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/mocheer/golib/file"
	"github.com/mocheer/pupilla/server"
)

func main() {
	var (
		wc  = flag.String("c", "./web.config.json", "web.config.json")
		p   = flag.String("p", "", "port")
		s   server.Server
		err error
	)
	flag.Parse()
	if file.Exist(*wc) {
		log.Println("WebServer:" + *wc)
		s, err = server.NewWebServer(*wc)
		if err != nil {
			log.Println("error on NewWebServer:", err)
			os.Exit(1)
		}
	} else {
		s = server.NewDefaultServer(*p)
	}
	err = s.Start()
	if err != nil {
		log.Println("error on server start:", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		flag, _, _ := reader.ReadLine()
		if string(flag) == "exit" {
			os.Exit(1)
		}
	}
}
