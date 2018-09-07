package main

import (
	"bufio"
	//	"encoding/binary"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gitdlam/telnet"
)

const timeout = 10 * time.Second

func checkErr(err error) {
	if err != nil {
		log.Println("error:", err.Error())
		return
	}
}

func expect(t *telnet.Conn, d ...string) {

	checkErr(t.SetReadDeadline(time.Now().Add(timeout)))
	checkErr(t.SkipUntil(d...))
}

func sendln(t *telnet.Conn, s string) {
	checkErr(t.SetWriteDeadline(time.Now().Add(timeout)))
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'
	_, err := t.Write(buf)
	checkErr(err)
}

func sendTab(t *telnet.Conn, s string) {
	checkErr(t.SetWriteDeadline(time.Now().Add(timeout)))
	buf := make([]byte, len(s))
	copy(buf, s)
	_, err := t.Write(buf)
	checkErr(err)
	time.Sleep(500 * time.Millisecond)
	buf2 := make([]byte, 1)
	copy(buf2, "\t")
	_, err = t.Write(buf2)
	checkErr(err)

}

func sendX(t *telnet.Conn) {

	checkErr(t.SetWriteDeadline(time.Now().Add(timeout)))
	//binary.PutVarint(buf, 24)
	_, err := t.Write([]byte{0x18})
	checkErr(err)
	time.Sleep(100 * time.Millisecond)
}

func sendA(t *telnet.Conn) {
	checkErr(t.SetWriteDeadline(time.Now().Add(timeout)))
	_, err := t.Write([]byte{0x01})
	checkErr(err)
	time.Sleep(100 * time.Millisecond)
}

func sendO(t *telnet.Conn) {
	checkErr(t.SetWriteDeadline(time.Now().Add(timeout)))
	_, err := t.Write([]byte{0x0F})
	checkErr(err)
	time.Sleep(100 * time.Millisecond)
}
func pause1() {
	time.Sleep(time.Second)
}

func tinyPause() {
	time.Sleep(50 * time.Millisecond)
}

func skipNextResponse(t *telnet.Conn) {
	checkErr(t.SetReadDeadline(time.Now().Add(1500 * time.Millisecond)))
	err := t.SkipNextResponse()
	if err != nil {
		if !strings.Contains(err.Error(), "timeout") {
			log.Println("skip error:", err.Error())
		}
	}
}

func main() {

	//	file, err := os.Open(os.Args[1])
	file, err := os.Open("login.txt")
	if err != nil {
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ":")
		host, user, passwd, user2, passwd2, menu := s[0], s[1], s[2], s[3], s[4], s[5]
		t, err := telnet.Dial("tcp", host+":23")
		checkErr(err)
		t.SetUnixWriteMode(false)

		//		var data []byte

		t.WriteRaw([]byte{0xFF, 0xFD, 0x03, 0xFF, 0xFB, 0x18, 0xFF, 0xFB, 0x1F, 0xFF, 0xFB, 0x20, 0xFF, 0xFB, 0x21, 0xFF, 0xFB, 0x22, 0xFF, 0xFB, 0x27, 0xFF, 0xFD, 0x05, 0xFF, 0xFB, 0x23})

		skipNextResponse(t)

		t.WriteRaw([]byte{0xFF, 0xFA, 0x1F, 0x00, 0x50, 0x00, 0x18, 0xFF, 0xF0, 0xFF, 0xFA, 0x18, 0x00, 0x78, 0x74, 0x65, 0x72, 0x6D, 0xFF, 0xF0})
		log.Println("before login: ", user)
		expect(t, "Login:")
		log.Println("found login: ")
		sendln(t, user)
		expect(t, "Password:")
		log.Println("found password")
		sendln(t, passwd)
		pause1()
		skipNextResponse(t)

		sendTab(t, user2)
		skipNextResponse(t)
		sendln(t, passwd2)
		log.Println("logged in")
		expect(t, "S=>")
		pause1()
		sendln(t, menu)
		expect(t, "Plt #:")
		pause1()
		sendln(t, "PLT33454")
		skipNextResponse(t)
		sendA(t)
		//		expect(t, "LPN #:")
		//		log.Println("LPN prompt")
		pause1()

		LPNs := []string{"12876293", "12827685", "12943166", "12033059", "12903995", "12950112", "12949357", "12950111"}

		for _, v := range LPNs {
			sendln(t, v)
			log.Println(v)
			skipNextResponse(t)
			sendA(t)
		}

		pause1()
		sendO(t)
		skipNextResponse(t)
		sendA(t)

		skipNextResponse(t)
		sendX(t)
		skipNextResponse(t)
		sendX(t)
		skipNextResponse(t)
		sendX(t)
		skipNextResponse(t)
		sendX(t)
		log.Println("done")

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
