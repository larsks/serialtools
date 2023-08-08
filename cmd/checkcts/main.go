package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"serialtools/version"

	"go.bug.st/serial"
)

var (
	showVersion bool
)

func init() {
	flag.BoolVar(&showVersion, "version", false, "show version and exit")
}

func realMain() int {
	flag.Parse()

	if showVersion {
		fmt.Println(version.GetVersionString())
		return 0
	}

	if len(os.Args) <= 1 {
		log.Printf("ERROR: you must provide a serial device")
		return 2
	}

	portName := os.Args[1]
	log.Printf("opening %s", portName)
	port, err := serial.Open(portName, &serial.Mode{})
	if err != nil {
		log.Printf("ERROR: failed to open serial port: %v", err)
		return 3
	}
	defer port.Close()

	status, err := port.GetModemStatusBits()
	if err != nil {
		log.Printf("ERROR: unable to read status bits from %s", portName)
		return 3
	}

	if status.CTS {
		log.Printf("CTS is set")
		return 0
	} else {
		log.Printf("CTS is not set")
		return 1
	}
}

func main() {
	os.Exit(realMain())
}
