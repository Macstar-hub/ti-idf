package httppost

import (
	"flag"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	port                       string
	clients                    = make(map[*websocket.Conn]bool)
	broadcast                  = make(chan []byte)
	mutex                      = &sync.Mutex{}
	path                       string
	fileLastSeekPostitionSlice []int
)

func init() {
	flag.StringVar(&port, "a", ":9080", "Server Run Port")
	flag.StringVar(&path, "p", "/Users/Shared/codes.dir/go.dir/git.dir/ti-idf/logs/redis/debug.log", "Path to watch for watchdog")
}

func conncetionUpgrader(readBuffer int, writeBuffer int, enableCompression bool) websocket.Upgrader {
	upgrader := websocket.Upgrader{
		ReadBufferSize:    readBuffer,
		WriteBufferSize:   writeBuffer,
		EnableCompression: false,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins for simplicity, adjust as needed
			return true
		},
	}
	return upgrader
}

// Make a listener and then produce to html ui
func Producer(redisChannle chan int, wg *sync.WaitGroup) {
	for status := range redisChannle {
		time.Sleep(100 * time.Microsecond)
		if status == 100 {
			wg.Done()
			close(redisChannle)
		}
	}
}

func sendMessage(router *gin.Engine) {
	// Load client page:
	router.GET("/api/v1/uploadfile", func(c *gin.Context) {
		c.File("../../web/landingPages/upload.html")
	})

	// Make Upgrader for WebSocket Connections
	upgrader := conncetionUpgrader(1024, 1024, false)

	// Make get and sent message.
	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Error upgrading connection:", err)
			return
		}
		defer conn.Close()

		// Make register connection:
		mutex.Lock()
		clients[conn] = true
		mutex.Unlock()

		for {
			_, message, err := conn.ReadMessage()
			log.Printf("Message is %s", string(message))
			if err != nil {
				log.Println("Cannot read message from websocket ...", err)
				break
			}
			broadcast <- message
		}
	})
}

// Make broadcasting message to all clients:
func broadcastMessage() {
	for {
		message := <-broadcast
		// log.Println("Read Line: ", message)
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Cannot write message with error :", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
