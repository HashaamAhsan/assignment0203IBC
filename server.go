package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	a1 "github.com/HashaamAhsan/assignment01IBC"
)

var chainHead *a1.Block
var msgarray [100]string
var addarray [100]string
var size int = 0

func main() {

	var nodes int
	fmt.Print("Give X no of nodes to connect to: ")
	fmt.Scan(&nodes)

	ln, err1 := net.Listen("tcp", "localhost:6000")
	if err1 != nil {

		log.Fatal(err1)

	}

	var s string
	var s1 string
	chainHead = a1.InsertBlock("Satoshis100", nil)
	//chainHead = a1.InsertBlock("Satoshi50Alice50", chainHead)
	go servertalk()
	for {
		if nodes > 0 {
			fmt.Print("No of nodes remaining: ", nodes)
			nodes = nodes - 1
		}
		conn, err2 := ln.Accept()

		log.Println("A client has connected", conn.RemoteAddr())
		//m := sync.Mutex{}

		//conarray[size] = conn
		if err2 != nil {

			log.Println(err2)
			continue

		}

		s = "localhost:600"
		s = s + strconv.Itoa(size+1)
		s1 = "$localhost:500"
		s1 = s1 + strconv.Itoa(size+1)
		addarray[size] = s
		msgarray[size] = s1[1:]
		s = s + s1 + strconv.Itoa(size+1)

		fmt.Fprintf(conn, s+"\n")

		size = size + 1
		println("Size: ", size)

		//message, _ := bufio.NewReader(conn).ReadString('\n')
		//fmt.Print("Message from client: " + message)
		time.Sleep(1000 * time.Millisecond)
		gobEncoder := gob.NewEncoder(conn)
		err3 := gobEncoder.Encode(chainHead)

		if err3 != nil {

			log.Println(err3)

		}
		conn.Close()
		//printconn(conarray, size)
		chainHead = a1.InsertBlock("Satoshis200", chainHead)
		broadcast()
		if nodes == 0 {
			updatingaboutpeers()
		}
		go handleConnection(s1[1:])

	}

}

func handleConnection(messageadrr string) {

	//var miner int
	//var text string

	// text = "2"
	// carr[index].Write([]byte(text + "\n"))
	// // will listen for message to process ending in newline (\n)
	// text = "3"
	//carr[index].Write([]byte(text + "\n"))
	var mesadrr string = messageadrr
	fmt.Print("server message address: ", mesadrr+"\n")
	for true {
		ln, err1 := net.Listen("tcp", mesadrr)
		if err1 != nil {
			//handle error
			log.Println(err1)
		}

		conn2, err2 := ln.Accept()
		if err2 != nil {
			//handle error
			log.Println(err2)
		}

		dec := gob.NewDecoder(conn2)
		dec.Decode(&chainHead)

		//m.Lock()
		// output message received
		//fmt.Print("Message Received:", string(message))
		//*chainHead = a1.InsertBlock(message, *chainHead)

		//miner = 2
		//fmt.Fprintf(carr[miner], text+"\n")
		//carr[miner].Write([]byte(text + "\n"))
		//message, _ := bufio.NewReader(carr[miner]).ReadString('\n')
		//fmt.Print("Message from client: " + message)
		// sample process for string received
		//newmessage := strings.ToUpper("updated server blockchain")
		//carr[index].Write([]byte(newmessage + "\n"))
		fmt.Print("every thing okay server side\n")
		var x int = 0
		for true {
			x = rand.Intn(size)
			if msgarray[x] != mesadrr {
				break
			}
		}
		var b bool = miner(x)
		time.Sleep(1000 * time.Millisecond)
		if b == true {
			broadcast()
		} else {
			chainHead = chainHead.Prevpointer
		}
		a1.ListBlocks(chainHead)
		conn2.Close()
		ln.Close()
		// send new string back to client
		//m.Unlock()
	}
}

func broadcast() {
	for i := 0; i < size; i++ {
		conn1, err := net.Dial("tcp", addarray[i])
		if err != nil {
			//handle error
			log.Println(err)
		}
		fmt.Fprintf(conn1, "update\n")

		time.Sleep(1000 * time.Millisecond)

		gobEncoder := gob.NewEncoder(conn1)
		err1 := gobEncoder.Encode(chainHead)
		if err1 != nil {

			log.Println(err)

		}
		conn1.Close()
	}
}

func printconn(carr [100]net.Conn, size int) {
	for i := 0; i < size; i++ {
		log.Println("A client is in the list", carr[i].RemoteAddr())
	}
}

func miner(index int) bool {
	conn1, err := net.Dial("tcp", addarray[index])
	if err != nil {
		//handle error
		log.Println(err)
	}
	fmt.Fprintf(conn1, "mining\n")

	// time.Sleep(100 * time.Millisecond)

	// dec := gob.NewDecoder(conn1)
	// dec.Decode(&chainHead)

	// time.Sleep(100 * time.Millisecond)
	var minerblock string
	message, _ := bufio.NewReader(conn1).ReadString('\n')
	fmt.Print(message)
	message = strings.TrimSuffix(message, "\n")
	message, minerblock = slicestring(message, "$")

	conn1.Close()
	if message == "Mining Complete....." {
		chainHead = a1.InsertBlock(minerblock, chainHead)
		return true
	}
	return false
}

func updatingaboutpeers() {
	for i := 0; i < size; i++ {
		conn1, err := net.Dial("tcp", addarray[i])
		if err != nil {
			//handle error
			log.Println(err)
		}
		fmt.Fprintf(conn1, "peerupdate\n")

		time.Sleep(1000 * time.Millisecond)
		var peerarr [100]string
		fmt.Print("Node num: " + strconv.Itoa(i) + "\n")
		var queue []int
		for k := 0; k <= i; k++ {
			for true {
				var a int = rand.Intn(size)
				if a != i {
					//peerarr[k] = addarray[a]
					if i < size-1 {
						if searchq(queue, a) == true {
							queue = append(queue, a)
							peerarr[k] = addarray[a]
							break
						}
					}
					if i == size-1 {
						peerarr[k] = addarray[a]
						break
					}
					//break
				}
			}
		}
		// fmt.Print("Printing peerarr\n")
		// for j := 0; j <= i; j++ {
		// 	fmt.Print(peerarr[j] + ",")
		// }
		// fmt.Print("\n")

		gobEncoder := gob.NewEncoder(conn1)
		err1 := gobEncoder.Encode(peerarr)
		if err1 != nil {

			log.Println(err)

		}
		conn1.Close()
	}
}

func servertalk() {
	var option int
	m := sync.Mutex{}
	for {
		fmt.Print("Give your option: ")
		fmt.Scan(&option)
		if option == 1 {
			m.Lock()
			var text string
			fmt.Print("Give message: ")
			fmt.Scan(&text)
			chainHead = a1.InsertBlock(text, chainHead)
			var b bool = miner(rand.Intn(size - 1))
			if b == true {
				broadcast()
			} else {
				chainHead = chainHead.Prevpointer
			}
			m.Unlock()
		}
		if option == 2 {
			m.Lock()
			a1.ListBlocks(chainHead)
			m.Unlock()
		}
	}
}

func searchq(queue []int, x int) bool {
	for i := 0; i < len(queue); i++ {
		if queue[i] == x {
			return false
		}
	}
	return true
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
