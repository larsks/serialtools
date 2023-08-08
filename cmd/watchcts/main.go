package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"serialtools/version"
	"time"

	"go.bug.st/serial"
)

var (
	upScript, downScript string
	missingOkay          bool
	checkInterval        time.Duration
	cooldownInterval     time.Duration
	showVersion          bool
)

func init() {
	flag.StringVar(&upScript, "up", "", "run this script when CTS goes high")
	flag.StringVar(&downScript, "down", "", "run this script when CTS goes low")
	flag.BoolVar(&missingOkay, "missingok", false, "do not error out if script is missing")
	flag.DurationVar(&checkInterval, "checkInterval", 1*time.Second, "port check interval")
	flag.DurationVar(&cooldownInterval, "cooldownInterval", 10*time.Second, "cool down interval")
	flag.BoolVar(&showVersion, "version", false, "show version and exit")
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

	if showVersion {
		fmt.Println(version.GetVersionString())
		return 0
	}

	if flag.NArg() < 1 {
		panic(fmt.Errorf("you must provide a serial device"))
	}

	portName := flag.Arg(0)
	port, err := serial.Open(portName, &serial.Mode{})
	if err != nil {
		panic(fmt.Errorf("failed to open serial port %s: %v", portName, err))
	}
	defer port.Close()

	log.Printf("watching serial port %s", portName)

	status, err := port.GetModemStatusBits()
	if err != nil {
		panic(fmt.Errorf("unable to read status bits from %s", portName))
	}
	ctsState := status.CTS

	for {
		status, err := port.GetModemStatusBits()
		if err != nil {
			panic(fmt.Errorf("unable to read status bits from %s", portName))
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
					panic(fmt.Errorf("ERROR: %s", err))
				}
			}

			time.Sleep(cooldownInterval)
		}

		time.Sleep(checkInterval)
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("ERROR: %s", r)
			os.Exit(1)
		}
	}()

	os.Exit(realMain())
}
