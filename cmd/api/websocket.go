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

// func fileWatchDog(path string) {

// 	// Make new watcher instance:
// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		log.Println("Cannot make new file watchdog with error: ", err)
// 	}

// 	// Send event throw channel:
// 	go func() {
// 		for {
// 			select {
// 			case event := <-watcher.Event:
// 				log.Printf("Event is: %s", event)

// 				// send event throw channel when file has been changed:
// 				if event.IsModify() {
// 					// Trigger only when modification not attribute change:
// 					if event.IsAttrib() == false {

// 						// Make last seek position:
// 						fileLastSeekPostitionSlice = append(fileLastSeekPostitionSlice, fileSeekPosition(path))
// 						sendLinesTrigger(path, fileLastSeekPostitionSlice[len(fileLastSeekPostitionSlice)-3])
// 					}
// 				}
// 			}
// 		}
// 	}()

// 	// Make watcher:
// 	err = watcher.Watch(path)
// 	if err != nil {
// 		log.Panicln("Cannot watch file: ", err)
// 	}
// }

// func sendLinesTrigger(path string, lastSize int) {
// 	// Open file from os:
// 	file, err := os.Open(path)
// 	if err != nil {
// 		log.Println("Cannot open file with error: ", err)
// 	}
// 	defer file.Close()

// 	// Set Seek Position:
// 	_, err = file.Seek(int64(lastSize), io.SeekStart)
// 	if err != nil {
// 		log.Println("Cannot set seek postion on file: ", err)
// 	}

// 	// Sync Input from file:
// 	scanner := bufio.NewScanner(file)

// 	// Read line by line:
// 	scanner.Split(bufio.ScanLines)

// 	// Send Line throw channel:
// 	for scanner.Scan() {
// 		broadcast <- scanner.Bytes()
// 	}
// 	// Make close file and clean up all attribute releated.
// 	file.Close()
// }

// func fileSeekPosition(path string) int {
// 	// Make file stats file:
// 	fileInfo, err := os.Stat(path)
// 	if err != nil {
// 		log.Println("Cannot stats file with error: ", err)
// 	}

// 	// Make current file size:
// 	currentFileSize := fileInfo.Size()
// 	fileLastSeekPostitionSlice = append(fileLastSeekPostitionSlice, int(currentFileSize))

// 	return int(currentFileSize)
// }

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
