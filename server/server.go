package server

import (
	"log"
	"net/http"

	"github.com/mocheer/golib/cmd"
	"github.com/mocheer/golib/format"
	"github.com/mocheer/golib/web/urlserver"
)

const (
	HTML  = "html"
	JSON  = "json"
	PJSON = "pjson"
)

type Server interface {
	Start() error
}

// FileServer
// URLServer
type webserver struct {
	config *WebConfig
	// cookie []*http.Cookie
}

func NewWebServer(filePath string) (Server, error) {
	var webConfig, err = NewWebConfig(filePath)
	if err != nil {
		return nil, err
	}
	s := &webserver{webConfig}
	return s, nil
}

func NewDefaultServer(p string) Server {
	var webConfig = NewDefaultConfig()
	if p != "" {
		webConfig.Port = p
	}
	return &webserver{webConfig}
}

func (this *webserver) Start() error {
	go this.StartCMD()
	return this.StartListen()
}

func (this *webserver) StartCMD() {
	c := this.config
	if c.Start != nil {
		switch c.Start.(type) {
		case string:
			cmd.Open([]string{c.Start.(string)})
		case []interface{}:
			start := c.Start.([]interface{})
			for _, param := range start {
				switch param.(type) {
				case string:
					cmd.Open([]string{param.(string)})
				case []interface{}:
					params := []string{}
					for _, p := range param.([]interface{}) {
						params = append(params, p.(string))
					}
					cmd.Start(params)
				}
			}
		}
	}
}

func (this *webserver) StartListen() error {
	c := this.config
	if c.FileServer != nil {
		this.StartFileServer()
	}
	if c.URLServer != "" {
		this.StartURLServer()
	}
	err := http.ListenAndServe(":"+c.Port, nil)
	if err != nil {
		return err
	}
	return nil
}

func (this *webserver) StartFileServer() {
	c := this.config
	for d, k := range c.FileServer {
		log.Println("FileServer:", k, d)
		http.Handle(k, http.StripPrefix(k, http.FileServer(http.Dir(d)))) //如果不用http.StripPrefix，则k只能为"/"
	}
}

func (this *webserver) StartURLServer() {
	c := this.config
	log.Println("UrlServer:", c.URLServer)
	http.HandleFunc(c.URLServer, urlHandle)
}

func urlHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	u := r.FormValue("u")
	if u == "" {
		w.Write([]byte("error of the url"))
		return
	}
	m := r.FormValue("m")
	values := r.Form
	values.Del("u")
	values.Del("m")
	log.Println(u)
	urlInfo := &urlserver.UrlInfo{Url: u, Values: values}
	var s string
	switch m {
	case "get":
		s = urlserver.Get(urlInfo)
	case "post":
		s = urlserver.Post(urlInfo)
	case "test":
		s = urlserver.Test(urlInfo)
	default:
		s = urlserver.Get(urlInfo)
	}
	log.Println(u)
	resBytes := []byte(s)
	f := r.FormValue("f")
	switch f {
	case PJSON:
		jsonFormatter := &format.JSONFormatter{}
		b := jsonFormatter.Format(resBytes)
		w.Write(b)
	default:
		w.Write(resBytes)
	}
}
