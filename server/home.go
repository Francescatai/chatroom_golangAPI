package server

import (
	// "encoding/json"
	"fmt"
	"html/template"
	"net/http"

	// "chatsystem/global"
	// "chatsystem/logic"
)


func homeHandleFunc(w http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles(rootDir + "/template/home.html")
	if err != nil {
		fmt.Fprint(w, "頁面解析異常！")
		return
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Fprint(w, "頁面執行錯誤！")
		return
	}
}