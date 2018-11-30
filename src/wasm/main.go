package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/weave-blueprint/src/blueprint"
	"os"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			fmt.Fprintf(os.Stdout, "exit err %v", err)
			os.Exit(1)
		}
	}()
	reader := bufio.NewReader(os.Stdin)
	var totalSize, methodStrLen int16
	err = binary.Read(reader, binary.BigEndian, &totalSize)
	if err != nil {
		return
	}

	err = binary.Read(reader, binary.BigEndian, &methodStrLen)
	if err != nil {
		return
	}

	sb := make([]byte, methodStrLen)
	_, err = reader.Read(sb)
	if err != nil {
		return
	}

	var argc int16
	err = binary.Read(reader, binary.BigEndian, &argc)
	if err != nil {
		return
	}
	for  i := int16(0); i < argc; i++ {
		var al int16
		err = binary.Read(reader, binary.BigEndian, &al)
		if err != nil {
			return
		}
		ab := make([]byte, al)
		_, err = reader.Read(ab)
		if err != nil {
			return
		}

	}

	meta, err := blueprint.Contract.Meta()
	if err != nil {
		return
	}
	data, err := proto.Marshal(&meta)
	if err != nil {
		return
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		return
	}
}
