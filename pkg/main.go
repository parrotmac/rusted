package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/parrotmac/rusted/pkg/modem"
	"github.com/parrotmac/rusted/pkg/server"
	"github.com/parrotmac/rusted/pkg/utils"
)

type Rusted struct {
	*utils.Features
	remote *server.Remote
	dev    *modem.Device
}

func (r *Rusted) initApp() {
	logrus.Debugln("[APP INIT] Starting Init")

	portNames, err := modem.GetPortNames()
	if err != nil {
		logrus.Printf("Error listing serial ports: %v", err)
	}

	for _, port := range portNames {
		fmt.Printf("Found serial port: %v\n", port)
	}

	dev, err := modem.FindLowSpeedHuaweiModemPort()
	if err != nil {
		logrus.Warnf("Failure finding Huawei modem: %v", err)
	} else {
		logrus.Println("Found Huawei device!")
	}
	r.dev = dev

	r.remote = &server.Remote{}

	wrapper := r.remote.ConnectMqttWrapper("tcp://mqtt.stag9.com:1883")

	wrapper.ATCommandHandler = func(atCommand string) string {
		resp, err := r.dev.SendModemCommandWithDeadline(atCommand, time.Second*1)
		if err != nil {
			logrus.Errorf("Got err: %v", err)
			return ""
		}
		return resp
	}

	r.Features = &utils.Features{
		MockSerial: utils.GetEnvBool("RUSTED_MOCK_SERIAL"),
	}
	logrus.Debugln("[APP INIT] Init Finished")
}

func (r *Rusted) runLoop() {
	logrus.Warnln("[MAIN LOOP] Not yet implemented")
	for {
		time.Sleep(time.Second + 1)
	}
}

func main() {
	logrus.Info("Starting Rusted")

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	r := &Rusted{}

	r.initApp()

	r.runLoop()

	logrus.Warnln("Exiting...")
}
