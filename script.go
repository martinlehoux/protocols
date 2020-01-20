package protocols

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Command struct {
	Cmd  string
	Args map[string]string
}

type DeviceData struct {
	Name string
}

type Config struct {
	Create  []DeviceData
	Connect [][2]string
	Script  []Command
}

func Script() {
	devices := make(map[string]*Device)
	filename := os.Args[1]
	var config Config
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	for _, data := range config.Create {
		device := CreateDevice(data.Name)
		devices[data.Name] = &device
	}
	for _, devicePair := range config.Connect {
		device1 := devices[devicePair[0]]
		device2 := devices[devicePair[1]]
		err := Connect(device1, device2)
		if err != nil {
			fmt.Printf("can't connect devices: %v\n", err)
		}
	}
	for _, device := range devices {
		fmt.Printf("Booting up %v\n", device.nickname)
		device.Run()
		device.PrintIfnetsTable()
	}
	for _, command := range config.Script {
		switch command.Cmd {
		case "send":
			devices[command.Args["from"]].SendPacket(devices[command.Args["to"]].MAC, []byte(command.Args["packet"]))
		}
		time.Sleep(time.Second)
	}
}
