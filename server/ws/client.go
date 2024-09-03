package ws

// import (
// 	"bytes"
// 	"fmt"
// 	"time"

// 	"github.com/gofiber/websocket/v2"
// )

// const (
// 	//time allowed to write a meesage to peer
// 	writeWait = 10 * time.Second
// 	//time allowed to read the next pong message from the peer
// 	pongWait = 60 * time.Second
// 	//send pings to peer with this period. must be less that pongwait
// 	pingPeriod = pongWait * 9 / 10

// 	maxMessageSize = 512
// )

// var (
// 	newLine = []byte{'\n'}
// 	space   = []byte{' '}
// )

// type Client struct {
// 	Hub *Hub
// 	//the websocket connection
// 	Conn *websocket.Conn
// 	//buffered channel of outbound messages
// 	Send chan []byte
// }

// // readPump pumps messages from the websocket connection to the hub
// //
// // The application runs readPump in a per-connection goroutine. The application
// // ensures that there is at most one reader on a connection by executing all
// // reads from this goroutine.
// func (c *Client) ReadPump() {
// 	defer func() {
// 		c.Hub.Unregister <- c
// 		c.Conn.Close()
// 	}()

// 	c.Conn.SetReadLimit(maxMessageSize)
// 	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
// 	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

// 	for {
// 		_, message, err := c.Conn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				fmt.Println("error : %", err)
// 			}
// 			break
// 		}
// 		message = bytes.TrimSpace(bytes.Replace(message, newLine, space, -1))
// 		c.Hub.Broadcast <- message
// 	}
// }

// // writePump pumps messages from the hub to the websocket connection.
// //
// // A goroutine running writePump is started for each connection. The
// // application ensures that there is at most one writer to a connection by
// // executing all writes from this goroutine.
// func (c *Client) WritePump() {
// 	ticker := time.NewTicker(pingPeriod)
// 	defer func() {
// 		ticker.Stop()
// 		c.Conn.Close()
// 	}()
// 	for {
// 		select {
// 		case message, ok := <-c.Send:
// 			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if !ok {
// 				// The hub closed the channel.
// 				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}

// 			w, err := c.Conn.NextWriter(websocket.TextMessage)
// 			if err != nil {
// 				return
// 			}
// 			w.Write(message)

// 			// Add queued chat messages to the current websocket message.
// 			n := len(c.Send)
// 			for i := 0; i < n; i++ {
// 				w.Write(newLine)
// 				w.Write(<-c.Send)
// 			}

// 			if err := w.Close(); err != nil {
// 				return
// 			}
// 		case <-ticker.C:
// 			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				return
// 			}
// 		}
// 	}
// }
