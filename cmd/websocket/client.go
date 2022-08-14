package websocket

import (
	"context"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, h, err := websocket.Dial(ctx, "ws://localhost:2022/ws", nil)
	log.Println(conn, h, err)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close(websocket.StatusInternalError, "the sky is falling")
	err = wsjson.Write(ctx, conn, "hi")
	if err != nil {
		log.Println(err)
		return
	}

	conn.Close(websocket.StatusNormalClosure, "")
}
