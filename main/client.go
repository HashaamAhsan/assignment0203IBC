package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	a1 "github.com/HashaamAhsan/assignment01IBC"
)

var listenaddress string
var listenaddress1 string
var chainHead *a1.Block
var peerarray [100]string
var maxpeers int = 0

func main() {

	conn, err := net.Dial("tcp", "localhost:6000")

	//var size int = 0
	//var option = 0
	if err != nil {
		//handle error
	}

	message, _ := bufio.NewReader(conn).ReadString('\n')
	//fmt.Print("Message from server: " + message)
	listenaddress = strings.TrimSuffix(message, "\n")
	listenaddress, listenaddress1 = slicestring(listenaddress, "$")
	maxpeers, _ = strconv.Atoi(listenaddress1[len(listenaddress1)-1:])
	listenaddress1 = listenaddress1[:len(listenaddress1)-1]
	fmt.Print("both address: " + listenaddress + " -- " + listenaddress1 + "\n")
	fmt.Print("Max Peers: " + strconv.Itoa(maxpeers) + "\n")
	dec := gob.NewDecoder(conn)
	err = dec.Decode(&chainHead)
	if err != nil {

		log.Println(err)

	}
	conn.Close()

	go update()
	var option int
	m := sync.Mutex{}
	for {
		// read in input from stdin
		fmt.Print("Give your option: ")
		fmt.Scan(&option)
		// option, _ := bufio.NewReader(conn).ReadString('\n')
		// fmt.Print("Message Received:", string(option))

		if option == 1 {
			m.Lock()
			conn2, err2 := net.Dial("tcp", listenaddress1)
			if err2 != nil {
				log.Println(err)
			}

			var text string
			fmt.Print("Give message: ")
			fmt.Scan(&text)

			chainHead = a1.InsertBlock(text, chainHead)

			gobEncoder := gob.NewEncoder(conn2)
			err := gobEncoder.Encode(chainHead)
			if err != nil {

				log.Println(err)

			}
			// send to socket
			//fmt.Fprintf(conn, text+"\n")
			// listen for reply
			//message, _ := bufio.NewReader(conn).ReadString('\n')
			//fmt.Print("Message from server: " + message)
			//fmt.Print("everything okay client side")
			conn2.Close()
			m.Unlock()
		}
		if option == 2 {
			m.Lock()
			a1.ListBlocks(chainHead)
			m.Unlock()
		}
		if option == 3 {
			for i := 0; i < maxpeers; i++ {
				fmt.Print("Peer " + strconv.Itoa(i+1) + "Address: " + peerarray[i] + "\n")
			}
			fmt.Print("Give port no\n")
			var text string
			fmt.Scan(&text)
			connp, err := net.Dial("tcp", "localhost:"+text)
			if err != nil {
				//handle error
				log.Println(err)
			}
			fmt.Fprintf(connp, "msg:hello peer how are you?\n")
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func update() {
	for {
		//fmt.Print("I am in Update!!")
		//fmt.Print("Listen Address: ", listenaddress)
		ln, err := net.Listen("tcp", listenaddress)
		if err != nil {
			log.Fatal(err)

		}
		conn1, err := ln.Accept()
		if err != nil {

			log.Println(err)
			continue

		}
		message, _ := bufio.NewReader(conn1).ReadString('\n')
		message = strings.TrimSuffix(message, "\n")

		time.Sleep(1000 * time.Millisecond)

		if message == "update" {
			dec := gob.NewDecoder(conn1)
			dec.Decode(&chainHead)
			fmt.Print("I have Updated!!\n")
		}
		if message == "mining" {
			fmt.Print("Started Mining .....\n")
			chainHead = a1.InsertBlock("miningcharges", chainHead)
			//time.Sleep(100 * time.Millisecond)
			//gobEncoder := gob.NewEncoder(conn1)
			// err := gobEncoder.Encode(chainHead)
			// if err != nil {

			// 	log.Println(err)

			// }
			// time.Sleep(100 * time.Millisecond)
			conn1.Write([]byte("Mining Complete.....$miningcharges\n"))
		}
		if message == "peerupdate" {
			dec := gob.NewDecoder(conn1)
			dec.Decode(&peerarray)
			fmt.Print("I have Updated my peers!!\n")
		}
		if message[0:3] == "msg" {
			fmt.Printf("Message from the buddy: " + message[4:])
		}

		conn1.Close()
		ln.Close()
	}
}

func slicestring(x string, c string) (string, string) {
	i := strings.Index(x, c)
	//fmt.Println("Index: ", i)
	if i > -1 {
		chars := x[:i]
		arefun := x[i+1:]
		return chars, arefun
	}
	return "", ""
}

// func messagepassing(conn net.Conn, chainHead *a1.Block) {
// 	var option int
// 	m := sync.Mutex{}
// 	for {
// 		// read in input from stdin
// 		fmt.Print("Give your option: ")
// 		fmt.Scan(&option)
// 		if option == 1 {
// 			m.Lock()
// 			var text string
// 			fmt.Print("Give message: ")
// 			fmt.Scan(&text)

// 			chainHead = a1.InsertBlock(text, chainHead)

// 			gobEncoder := gob.NewEncoder(conn)
// 			err := gobEncoder.Encode(chainHead)
// 			if err != nil {

// 				log.Println(err)

// 			}
// 			// send to socket
// 			//fmt.Fprintf(conn, text+"\n")
// 			// listen for reply
// 			//message, _ := bufio.NewReader(conn).ReadString('\n')
// 			//fmt.Print("Message from server: " + message)
// 			m.Unlock()
// 		}
// 		if option == 2 {
// 			m.Lock()
// 			a1.ListBlocks(chainHead)
// 			m.Unlock()
// 		}
// 	}
// }
