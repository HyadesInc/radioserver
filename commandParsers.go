package main

import (
	"bytes"
	"encoding/binary"
)

func ParseCmdHelloBody(data []uint8) (Version, string) {
	var protocolVersion uint32
	var clientName string
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &protocolVersion)

	clientName = string(data[4:])

	return SplitProtocolVersion(protocolVersion), clientName
}

func ParseCmdGetSettingBody(data []uint8) {

}

func ParseCmdSetSettingBody(data []uint8) (setting uint32, args []uint32) {
	buf := bytes.NewReader(data)

	var numArgs = (len(data) - 4) / 4

	args = make([]uint32, numArgs)

	binary.Read(buf, binary.LittleEndian, &setting)
	for i := 0; i < numArgs; i++ {
		binary.Read(buf, binary.LittleEndian, &args[i])
	}

	return setting, args
}