package websocket

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func Run() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(w, "Http, Hello")
		fmt.Printf("%T", w)
		w.Write([]byte("this is response"))
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		conn, err := websocket.Accept(w, req, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close(websocket.StatusInternalError, "the sky is falling")
		ctx, cancel := context.WithTimeout(req.Context(), time.Second*10)
		defer cancel()

		var v interface{}
		err = wsjson.Read(ctx, conn, &v)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("received: %v\n", v)
		err = wsjson.Write(ctx, conn, "Hello websocket Client")
		if err != nil {
			log.Println(err)
			return
		}
		conn.Close(websocket.StatusNormalClosure, "")

	})

	log.Fatal(http.ListenAndServe(":2022", nil))
}
