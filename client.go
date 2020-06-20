// Copyright (c) 2020, RetailNext, Inc.
// This material contains trade secrets and confidential information of
// RetailNext, Inc.  Any use, reproduction, disclosure or dissemination
// is strictly prohibited without the explicit written permission
// of RetailNext, Inc.
// All rights reserved.

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/goburrow/modbus"
)

var port = flag.String("port", "/dev/tty.usbserial-1420", "")
var setCurrentTemp = flag.String("set-temp", "", "15.5")

func main() {
	flag.Parse()
	fmt.Println("Connecting to port", *port)
	handler := modbus.NewRTUClientHandler(*port)
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second

	client := modbus.NewClient(handler)
	result, err := client.ReadHoldingRegisters(0, 7)
	if err != nil {
		fmt.Println("ERROR READ:", err)
		os.Exit(1)
	}
	fmt.Println("RESULT", result)
	fmt.Println("Enabled", result[1])
	fmt.Println("Temperature", float64(int(result[2])*256+int(result[3]))/10.0)
	fmt.Println("Manual Setting", result[5] == 1)
	fmt.Println("Heating", result[7] == 1)
	fmt.Println("Set Temperature", float64(int(result[8])*256+int(result[9]))/10.0)
	fmt.Println("Weekly mode set Temperature", float64(int(result[10])*256+int(result[11]))/10.0)
	fmt.Println("Lock Enabled", result[13] == 1)

	if *setCurrentTemp != "" {
		t, err := strconv.ParseFloat(*setCurrentTemp, 32)
		if err != nil {
			fmt.Println("Incorrect temperature")
			os.Exit(1)
		}
		result, err = client.WriteSingleRegister(4, uint16(t*10))
		if err != nil {
			fmt.Println("Error setting temp", err)
			os.Exit(1)
		}
		fmt.Println("Set Temp result:", result)
	}
}
