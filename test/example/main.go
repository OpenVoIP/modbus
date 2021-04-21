package main

import (
	"log"
	"os"
	"time"

	modbus "github.com/OpenVoIP/modbus/pkg"
	modbusTCP "github.com/OpenVoIP/modbus/pkg/tcp"
	"github.com/OpenVoIP/modbus/pkg/utils"
)

func main() {
	logger := log.New(os.Stdout, "app: ", log.LstdFlags)

	// Modbus TCP
	handler := modbusTCP.NewTCPClientHandler("192.168.12.239:502")
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 1
	handler.Logger = utils.GetLogger()
	handler.Handle = func(data []byte) {
		// 主动上传数据
		logger.Printf("handle %+v\n", data)
	}

	// Connect manually so that multiple requests are handled in one connection session
	go func() {
	reconnect:
		err := handler.Connect()
		if err != nil {
			logger.Printf("Connect have error %+v\n", err)
		}

		time.Sleep(3 * time.Second)
		goto reconnect

	}()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		client := modbus.NewClient(handler)

		defer ticker.Stop()
		for {
			<-ticker.C
			results, err := client.ReadDiscreteInputs(200, 1)
			if err != nil {
				logger.Printf("have error %+v\n", err)
			} else {
				logger.Printf("get data %+v\n", results)
			}

		}
	}()

	select {}
}
