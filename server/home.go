package server

import (
	"fmt"
	"html/template"
	"net/http"
	"ws/global"
)

func homeHandleFunc(w http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles(global.RootDir + "/template/home.html")
	if err != nil {
		fmt.Println(w, "模版解析错误")
		return
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Println(w, "模版执行错误！")
		return
	}

}
