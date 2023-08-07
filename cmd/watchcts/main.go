package main

import (
	"log"
	"os"

	"go.bug.st/serial"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("ERROR: you must provide a serial device")
	}

	portName := os.Args[1]
	port, err := serial.Open(portName, &serial.Mode{})
	if err != nil {
		log.Fatalf("ERROR: failed to open serial port: %v", err)
	}

	for {
		status, err := port.GetModemStatusBits()
		if err != nil {
			log.Fatalf("ERROR: unable to read status bits from %s", portName)
		}

		if !status.CTS {
			log.Print("lost CTS")
			break
		}
	}
}
