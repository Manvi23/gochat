package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	pb "github.com/gautamrege/gochat/api"
	"net"
	"os"
	"sync"
	"time"
)

var (
	name = flag.String("name", "", "The name you want to chat as")
	port = flag.Int("port", 12345, "Port that your server will run on.")
)

func main() {
	// Parse flags for host, port and name
	flag.Parse()

	if *name == "" {
		fmt.Println("fuck off if you don't have a name")
		os.Exit(1)
	}
	// Create your own Global Handle ME
	var wg sync.WaitGroup
	wg.Add(4)

	// exit channel is a buffered channel for 5 exit patterns
	exit := make(chan bool, 5)

	// Listener for is-alive broadcasts from other hosts. Listening on 33333
	go registerHandle(&wg, exit)

	// Broadcast for is-alive on 33333 with own Handle.
	go isAlive(&wg, exit)

	// Cleanup Dead Handles
	go cleanupDeadHandles(&wg, exit)

	// gRPC listener
	go listen(&wg, exit)

	//for {
	//	// Loop indefinitely and render Term
	//	// When we need to exit, send true 3 times on exit channel!
	//}
	time.Sleep(1 * time.Second)

	h := pb.Handle{
		Name: "Anuj",
		Host: "192.168.1.18",
		Port: int32(10000),
	}
	sendChat(h, "wtf")

	// exit cleanly on waitgroup
	wg.Wait()
	close(exit)
}

// Broadcast Listener
// Listens on 33333 and updates the Global Handles list
func registerHandle(wg *sync.WaitGroup, exit chan bool) {
	defer wg.Done()
	// Check if the handle is already in HANDLES. If not, add a new one!

	localAddress, _ := net.ResolveUDPAddr("udp", ":333333")
	connection, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		fmt.Println(err)
	}
	defer connection.Close()
	fmt.Println("listening")

	h := Handle{}
	for {
		inputBytes := make([]byte, 4096)
		length, _, _ := connection.ReadFromUDP(inputBytes)
		buffer := bytes.NewBuffer(inputBytes[:length])
		decoder := gob.NewDecoder(buffer)
		decoder.Decode(&h)
		fmt.Println("listened data", h)
	}
}

// isAlive go-routine that publishes it's Handle on 33333
const listenerPort = 5000

func isAlive(wg *sync.WaitGroup, exit chan bool) {
	defer wg.Done()

	for {
		select {
		case <-exit:
			break
		default:
			conn, err := net.Dial("udp", "192.168.1.255:33333")
			if err != nil {
				fmt.Println(err)
			}
			defer conn.Close()
			handle := Handle{
				Name:       "Varsha",
				Port:       listenerPort,
				Host:       "192.168.1.135",
				Created_at: time.Now(),
			}
			handleJson, err := json.Marshal(handle)
			if err != nil {
				fmt.Println(err)
			}
			conn.Write(handleJson)
			fmt.Println("Brodcast: ", handle)
			time.Sleep(time.Second * 10)
		}
	}
}

// cleanup Dead Handlers
func cleanupDeadHandles(wg *sync.WaitGroup, exit chan bool) {
	defer wg.Done()
	// wait for DEAD_HANDLE_INTERVAL seconds before removing them from chatrooms and handle list
}
