package main

import (
	"bytes"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/plugin"
	"github.com/loomnetwork/go-loom/plugin/contractpb"
	"github.com/loomnetwork/weave-blueprint/src/blueprint"
	"github.com/loomnetwork/weave-blueprint/src/types"
	"log"
)

func main() {
	addr1 := loom.MustParseAddress("chain:0xb16a379ec18d4093666f8f38b11a3071c920207d")
	user := "test"
	encoding := plugin.EncodingType_JSON

	marshaler, err := contractpb.MarshalerFactory(encoding)
	if err != nil {
		log.Fatal(err)
	}

	var argsBuffer bytes.Buffer
	err = marshaler.Marshal(&argsBuffer, &types.BluePrintCreateAccountTx{
		Version: 1,
		Owner:   user,
		Data:    []byte(user),
	})
	if err != nil {
		log.Fatal(err)
	}

	msg := &plugin.ContractMethodCall{
		Method: "CreateAccount",
		Args:   argsBuffer.Bytes(),
	}

	var msgBuffer bytes.Buffer
	err = marshaler.Marshal(&msgBuffer, msg)
	if err != nil {
		log.Fatal(err)
	}

	req := &plugin.Request{
		ContentType: encoding,
		Accept:      encoding,
		Body:        msgBuffer.Bytes(),
	}
	resp, err := blueprint.Contract.Call(plugin.CreateFakeContext(addr1, addr1), req)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("response: %#v", resp)

	resp, err = blueprint.Contract.Call(plugin.CreateFakeContext(addr1, addr1), req)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("response: %#v", resp)
}
