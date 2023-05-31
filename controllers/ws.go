package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/adamlahbib/gitaz/lib"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:      func(r *http.Request) bool { return true },
		HandshakeTimeout: time.Duration(time.Second * 5),
}



func HandleLogsStreamingSocket(c *gin.Context) {
	// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	// c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	app := c.Param("app")

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Println(err)
		return
	}

	defer ws.Close()

	streamer := lib.Kafka.Reader(app)
	defer streamer.Close()

	streamer.Read(func(key string, msg []byte) {
		err = ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println(err)
		}
	}, func(err error, msg string) {
		// ws.Close()
		// streamer.Close()
	})

	// for {
	// 	_, msg, err := ws.ReadMessage()
	// 	if err != nil {
	// 		log.Println(err)
	// 		break
	// 	}

	// 	fmt.Println(string(msg))
		
	// 	if string(msg) == "ping" {
	// 		msg = []byte("pong" + app)
	// 	}

	// 	err = ws.WriteMessage(websocket.TextMessage, msg)
	// 	if err != nil {
	// 		log.Println(err)
	// 		break
	// 	}
	// }

}