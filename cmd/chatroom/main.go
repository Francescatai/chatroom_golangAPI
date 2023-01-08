package main

import(
	"fmt"
	"net/http"
	"log"

	_ "net/http/pprof"

	"chatsystem/server"
	"chatsystem/global"

)

var (
	addr   = ":2023"
	banner = `
    ____              _____
   |    |    |   /\     |
   |    |____|  /  \    | 
   |    |    | /----\   |
   |____|    |/      \  |

Go 語言學習項目：ChatRoom，start on：%s
`
)

func init() {
	global.Init()
}

func main() {
	fmt.Printf(banner+"\n", addr)

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}