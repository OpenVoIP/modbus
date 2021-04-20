// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package ascii

import (
	"bytes"
	"log"
	"os"
	"testing"

	modbus "github.com/OpenVoIP/modbus/internal/ascii"
)

const (
	asciiDevice = "/dev/pts/6"
)

func TestASCIIClient(t *testing.T) {
	// Diagslave does not support broadcast id.
	handler := modbus.NewASCIIClientHandler(asciiDevice)
	handler.SlaveId = 17
	ClientTestAll(t, modbus.NewClient(handler))
}

func TestASCIIClientAdvancedUsage(t *testing.T) {
	handler := modbus.NewASCIIClientHandler(asciiDevice)
	handler.BaudRate = 19200
	handler.DataBits = 8
	handler.Parity = "E"
	handler.StopBits = 1
	handler.SlaveId = 12
	handler.Logger = log.New(os.Stdout, "ascii: ", log.LstdFlags)
	err := handler.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadDiscreteInputs(15, 2)
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
	results, err = client.ReadWriteMultipleRegisters(0, 2, 2, 2, []byte{1, 2, 3, 4})
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
}

func TestASCIIEncoding(t *testing.T) {
	encoder := asciiPackager{}
	encoder.SlaveId = 17

	pdu := ProtocolDataUnit{}
	pdu.FunctionCode = 3
	pdu.Data = []byte{0, 107, 0, 3}

	adu, err := encoder.Encode(&pdu)
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte(":1103006B00037E\r\n")
	if !bytes.Equal(expected, adu) {
		t.Fatalf("adu actual: %v, expected %v", adu, expected)
	}
}

func TestASCIIDecoding(t *testing.T) {
	decoder := asciiPackager{}
	decoder.SlaveId = 247
	adu := []byte(":F7031389000A60\r\n")

	pdu, err := decoder.Decode(adu)
	if err != nil {
		t.Fatal(err)
	}

	if 3 != pdu.FunctionCode {
		t.Fatalf("Function code: expected %v, actual %v", 15, pdu.FunctionCode)
	}
	expected := []byte{0x13, 0x89, 0, 0x0A}
	if !bytes.Equal(expected, pdu.Data) {
		t.Fatalf("Data: expected %v, actual %v", expected, pdu.Data)
	}
}

func BenchmarkASCIIEncoder(b *testing.B) {
	encoder := asciiPackager{
		SlaveId: 10,
	}
	pdu := ProtocolDataUnit{
		FunctionCode: 1,
		Data:         []byte{2, 3, 4, 5, 6, 7, 8, 9},
	}
	for i := 0; i < b.N; i++ {
		_, err := encoder.Encode(&pdu)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkASCIIDecoder(b *testing.B) {
	decoder := asciiPackager{
		SlaveId: 10,
	}
	adu := []byte(":F7031389000A60\r\n")
	for i := 0; i < b.N; i++ {
		_, err := decoder.Decode(adu)
		if err != nil {
			b.Fatal(err)
		}
	}
}
