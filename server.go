package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/natefinch/npipe"
)

type Progress struct {
	direction string
	bytes     uint64
}

func main() {
	var pipeName string
	if len(os.Args) < 2 {
		pipeName = `\\.\pipe\mypipe`
	} else {
		pipeName = os.Args[1]
	}

	// Pipe erstellen
	listener, err := npipe.Listen(pipeName)
	if err != nil {
		log.Fatalf("Fehler beim Erstellen der Named Pipe: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Named Pipe %s wartet auf Verbindungen ...\n", pipeName)
	for {
		// Verbindung akzeptieren
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Fehler beim Akzeptieren der Verbindung: %v", err)
		}

		peer := conn.RemoteAddr().String()
		fmt.Println("Verbindung aufgebaut von ", peer)
		go processClient(conn)
	}

}

func processClient(conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	c := make(chan Progress)

	copy := func(r io.ReadCloser, w io.WriteCloser) {
		defer func() {
			r.Close()
			w.Close()
		}()
		n, _ := io.Copy(w, r)
		c <- Progress{bytes: uint64(n)}
	}

	go copy(conn, os.Stdout)
	go copy(os.Stdin, conn)

	p := <-c
	log.Printf("[%s]: Verbindung wurde von der Gegenstelle beendet, %d bytes wurden empfangen\n", conn.RemoteAddr(), p.bytes)
	p = <-c
	log.Printf("[%s]: %d bytes wurden gesendet\n", conn.RemoteAddr(), p.bytes)
	os.Exit(0)
}
