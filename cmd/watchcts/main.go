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
	cooldownInterval     time.Duration
)

func init() {
	flag.StringVar(&upScript, "up", "", "run this script when CTS goes high")
	flag.StringVar(&downScript, "down", "", "run this script when CTS goes low")
	flag.BoolVar(&missingOkay, "missingok", false, "do not error out if script is missing")
	flag.DurationVar(&checkInterval, "checkInterval", 1*time.Second, "port check interval")
	flag.DurationVar(&cooldownInterval, "cooldownInterval", 10*time.Second, "cool down interval")
}

func runScript(script, portName, action string) error {
	log.Printf("running script %s %s %s", script, portName, action)
	cmd := exec.Command(script, portName, action)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
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
			var action, script string
			ctsState = status.CTS
			if ctsState {
				log.Printf("CTS is high")
				action = "up"
				script = upScript
			} else {
				log.Printf("CTS is low")
				action = "down"
				script = downScript
			}

			if script != "" {
				if err := runScript(script, portName, action); err != nil && !missingOkay {
					log.Printf("ERROR: %s", err)
					return 1
				}
			}

			time.Sleep(cooldownInterval)
		}

		time.Sleep(checkInterval)
	}
}

func main() {
	os.Exit(realMain())
}
