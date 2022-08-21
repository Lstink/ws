package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"ws/global"
	"ws/logic"
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

// 获取用户列表
func userListHandleFunc(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	userList := logic.Broadcaster.GetUserList()
	b, err := json.Marshal(userList)
	//logic.Broadcaster.BroadCaster(logic.NewErrorMessage("大傻叉"))

	if err != nil {
		fmt.Fprint(w, `[]`)
	} else {
		fmt.Fprint(w, string(b))
	}

}
