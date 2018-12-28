package main

import (
	"fmt"
	"github.com/parrotmac/rusted/pkg/modem"
	log "github.com/sirupsen/logrus"

	"github.com/parrotmac/rusted/pkg/utils"
)

type Rusted struct {
	*utils.Features
}

func (r *Rusted) initApp() {
	log.Debugln("[APP INIT] Starting Init")

	err, portNames := modem.GetPortNames()
	if err != nil {
		log.Printf("Error listing serial ports: %v", err)
	}

	for _, port := range portNames {
		fmt.Printf("Found serial port: %v\n", port)
	}

	r.Features = &utils.Features{
		MockSerial: utils.GetEnvBool("RUSTED_MOCK_SERIAL"),
	}
	log.Debugln("[APP INIT] Init Finished")
}

func (r *Rusted) runLoop() {
	log.Warnln("[MAIN LOOP] Not yet implemented")
}

func main() {
	log.Info("Starting Rusted")

	r := &Rusted{}

	r.initApp()

	r.runLoop()

	log.Warnln("Exiting...")
}
