package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"go.bug.st/serial"
)

var (
	upScript, downScript string
	missingOkay          bool
	checkInterval        time.Duration
)

func init() {
	flag.StringVar(&upScript, "up", "", "run this script when CTS goes high")
	flag.StringVar(&downScript, "down", "", "run this script when CTS goes low")
	flag.BoolVar(&missingOkay, "missingok", false, "do not error out if script is missing")
	flag.DurationVar(&checkInterval, "interval", 10*time.Second, "port check interval")
}

func runScript(script, portName, action string) error {
	cmd := exec.Command(script, action, portName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s script failed: %w", action, err)
	}

	return nil
}

func realMain() int {
	flag.Parse()

	if flag.NArg() < 1 {
		log.Printf("ERROR: you must provide a serial device")
		return 2
	}

	portName := flag.Arg(0)
	port, err := serial.Open(portName, &serial.Mode{})
	if err != nil {
		log.Printf("ERROR: failed to open serial port %s: %v", portName, err)
		return 1
	}
	defer port.Close()

	log.Printf("watching serial port %s", portName)

	status, err := port.GetModemStatusBits()
	if err != nil {
		log.Printf("ERROR: unable to read status bits from %s", portName)
		return 1
	}
	ctsState := status.CTS

	for {
		status, err := port.GetModemStatusBits()
		if err != nil {
			log.Printf("ERROR: unable to read status bits from %s", portName)
			return 1
		}

		if status.CTS != ctsState {
			ctsState = status.CTS
			if ctsState {
				log.Printf("CTS is high")
				if upScript != "" {
					if err := runScript(upScript, "up", portName); err != nil && !missingOkay {
						log.Printf("ERROR: %s", err)
						return 1
					}
				}
			} else {
				log.Printf("CTS is low")
				if downScript != "" {
					if err := runScript(downScript, "up", portName); err != nil && !missingOkay {
						log.Printf("ERROR: %s", err)
						return 1
					}
				}
			}

			time.Sleep(checkInterval)
		}
	}
}

func main() {
	os.Exit(realMain())
}
