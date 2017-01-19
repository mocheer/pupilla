package server

import (
	"encoding/json"
	"io/ioutil"
)

//WebConfig web config
type WebConfig struct {
	Start      interface{}       `json:"start"`      // 打开浏览器链接到start指定的url
	Port       string            `json:"port"`       // 监听端口
	FileServer map[string]string `json:"fileServer"` // 文件服务配置
	URLServer  string            `json:"urlServer"`  // get/post 代理 防止跨域
}

//NewWebConfig new WebConfig by json file
func NewWebConfig(filePath string) (*WebConfig, error) {
	c := &WebConfig{}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return c, nil
}

//NewDefaultConfig new default web config
func NewDefaultConfig() *WebConfig {
	return &WebConfig{
		Port: "8575",
		FileServer: map[string]string{
			"./": "/s",
		},
		URLServer: "/u/"}
}
