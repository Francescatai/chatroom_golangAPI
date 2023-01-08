package server

import (
	"log"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"chatsystem/logic"
)

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	// Accept 從客户端接受 WebSocket 握手，並將連接切換到 WebSocket。
	// 如果 Origin 域與主機不同，Accept 將拒絕握手，除非設置了 InsecureSkipVerify 選項（通過第三個參數 AcceptOptions 設置）。
	// 換句話說，默認情況下，它不允許跨源請求。如果發生錯誤，Accept 將始終寫入適當的回應
	conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		log.Println("websocket accept error:", err)
		return
	}

	// 1. 新用戶進來，構建該用戶實例
	token := req.FormValue("token")
	nickname := req.FormValue("nickname")
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal: ", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("非個人暱稱長度必須為：2-20字"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal!")
		return
	}
	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("個人暱稱已重名：", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("該暱稱已重名！"))
		conn.Close(websocket.StatusUnsupportedData, "nickname exists!")
		return
	}

	userHasToken := logic.NewUser(conn, token, nickname, req.RemoteAddr)

	// 2. 開啟給用戶發送消息的 goroutine
	go userHasToken.SendMessage(req.Context())

	// 3. 給當前用戶發送歡迎消息
	userHasToken.MessageChannel <- logic.NewWelcomeMessage(userHasToken)

	// 避免 token 泄露
	tmpUser := *userHasToken
	user := &tmpUser
	user.Token = ""

	// 通知所有用户新用户加入
	msg := logic.NewUserEnterMessage(user)
	logic.Broadcaster.Broadcast(msg)

	// 4. 將該用戶加入到廣播器的用戶列表中
	logic.Broadcaster.UserEntering(user)
	log.Println("user:", nickname, "joins chat")

	// 5. 接收用户訊息
	err = user.ReceiveMessage(req.Context())

	// 6. 用户下線
	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewUserLeaveMessage(user)
	logic.Broadcaster.Broadcast(msg)
	log.Println("user:", nickname, "leaves chat")

	// 根據讀取時的錯誤執行不同方式的 Close
	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client error:", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}
}