package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/go-redis/redis"
)

const (
	HOST = "localhost"
	PORT = "9000"
	TYPE = "tcp"
)

type sampleData struct {
	time   string
	symbol string
	open   float64
	high   float64
	low    float64
	close  float64
	volume int
}

func randString(n int) string {
	const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func handleIncomingRequest(conn net.Conn, client *redis.Client) {
	// store incoming data
	buffer := new(bytes.Buffer)
	_, err := conn.Read(buffer.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	// respond
	conn.Write([]byte("Hi back!\n"))
	generatedArray := makeArrayOfObjects(client)

	fmt.Println(generatedArray)
	go tickFunction(conn, generatedArray)

}

func makeArrayOfObjects(client *redis.Client) []sampleData {

	var myArray = []sampleData{}

	k := 1
	for ; k <= 10; k++ {
		item1 := sampleData{
			time:   time.Now().String(),
			symbol: randString(3),
			open:   100.00,
			high:   100.00,
			low:    100.00,
			close:  100.00,
			volume: 10000,
		}

		//Init a map[string]interface{}
		var m = make(map[string]interface{})
		m["time"] = item1.time
		m["symbol"] = item1.symbol
		m["open"] = item1.open
		m["high"] = item1.high
		m["low"] = item1.low
		m["close"] = item1.close
		m["volume"] = item1.volume
		hash, err := client.HMSet(item1.symbol, m).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println(hash)
		myArray = append(myArray, item1)
	}

	return myArray
}

func tickFunction(conn net.Conn, generatedArray []sampleData) {

	for range time.Tick(time.Second * 1) {
		conn.Write([]byte("\n1"))
	}
}

func main() {

	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println("Server is started on port", PORT)

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	// we can call set with a `Key` and a `Value`.
	// if there has been an error setting the value
	// handle the error
	if err != nil {
		fmt.Println(err)
	}

	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleIncomingRequest(conn, client)
	}
}
