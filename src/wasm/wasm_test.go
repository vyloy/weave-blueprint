package main_test

import (
	"context"
	"encoding/json"
	"github.com/loomnetwork/loomchain/log"
	. "github.com/loomnetwork/loomchain/plugin"
	"github.com/loomnetwork/weave-blueprint/src/types"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/go-loom"
	loom_plugin "github.com/loomnetwork/go-loom/plugin"
	"github.com/loomnetwork/loomchain"
	"github.com/loomnetwork/loomchain/eth/subs"
	registry "github.com/loomnetwork/loomchain/registry/factory"
	"github.com/loomnetwork/loomchain/store"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Implements loomchain.EventHandler interface
type fakeEventHandler struct {
}

func (eh *fakeEventHandler) Post(height uint64, e *loomchain.EventData) error {
	return nil
}

func (eh *fakeEventHandler) EmitBlockTx(height uint64) error {
	return nil
}

func (eh *fakeEventHandler) SubscriptionSet() *loomchain.SubscriptionSet {
	return nil
}

func (eh *fakeEventHandler) EthSubscriptionSet() *subs.EthSubscriptionSet {
	return nil
}

func TestWASMPluginVMContract(t *testing.T) {
	log.Setup("debug", "file://-")
	loader := NewWASMLoader("testdata")
	block := abci.Header{
		ChainID: "chain",
		Height:  int64(34),
		Time:    time.Unix(123456789, 0),
	}
	state := loomchain.NewStoreState(context.Background(), store.NewMemStore(), block, nil)
	createRegistry, err := registry.NewRegistryFactory(registry.LatestRegistryVersion)
	require.NoError(t, err)

	vm := NewPluginVM(loader, state, createRegistry(state), &fakeEventHandler{}, log.Default, nil, nil, nil)

	// Deploy contracts
	owner := loom.RootAddress("chain")
	goContractAddr1, err := deployGoContract(vm, "BluePrint:0.0.1", 0, owner)
	require.NoError(t, err)

	vmAddr1 := loom.MustParseAddress("chain:0xb16a379ec18d4093666f8f38b11a3071c920207d")
	s := "abc"
	payload := &types.BluePrintCreateAccountTx{
		Version: 1,
		Owner:   s,
		Data:    []byte(s),
	}
	require.NoError(t, callGoContractMethod(vm, vmAddr1, goContractAddr1, "CreateAccount", payload))
	require.EqualError(t, callGoContractMethod(vm, vmAddr1, goContractAddr1, "CreateAccount", payload), "Owner already exists")

	// test SaveState
	msgData := struct {
		Value int
	}{
		Value:10,
	}
	data, err := json.Marshal(msgData)
	require.NoError(t, err)

	msg := &types.BluePrintStateTx{
		Version: 1,
		Owner:   s,
		Data:    data,
	}
	require.NoError(t, callGoContractMethod(vm, vmAddr1, goContractAddr1, "SaveState", msg))
	require.EqualError(t, callGoContractMethod(vm, vmAddr1, goContractAddr1, "SaveStateNotExists", msg), "method not found")
}

func deployGoContract(vm *PluginVM, contractID string, contractNum uint64, owner loom.Address) (loom.Address, error) {
	init, err := proto.Marshal(&Request{
		ContentType: loom_plugin.EncodingType_PROTOBUF3,
		Accept:      loom_plugin.EncodingType_PROTOBUF3,
		Body:        nil,
	})
	if err != nil {
		return loom.Address{}, err
	}
	pc := &PluginCode{
		Name:  contractID,
		Input: init,
	}
	code, err := proto.Marshal(pc)
	if err != nil {
		return loom.Address{}, err
	}
	callerAddr := CreateAddress(owner, contractNum)
	_, contractAddr, err := vm.Create(callerAddr, code, loom.NewBigUIntFromInt(0))
	if err != nil {
		return loom.Address{}, err
	}
	return contractAddr, nil
}

func encodeGoCallInput(method string, inpb proto.Message) ([]byte, error) {
	args, err := proto.Marshal(inpb)
	if err != nil {
		return nil, err
	}

	body, err := proto.Marshal(&loom_plugin.ContractMethodCall{
		Method: method,
		Args:   args,
	})
	if err != nil {
		return nil, err
	}

	input, err := proto.Marshal(&loom_plugin.Request{
		ContentType: loom_plugin.EncodingType_PROTOBUF3,
		Accept:      loom_plugin.EncodingType_PROTOBUF3,
		Body:        body,
	})
	if err != nil {
		return nil, err
	}

	return input, nil
}

func callGoContractMethod(vm *PluginVM, callerAddr, contractAddr loom.Address, method string, inpb proto.Message) error {
	input, err := encodeGoCallInput(method, inpb)
	if err != nil {
		return err
	}
	_, err = vm.Call(callerAddr, contractAddr, input, loom.NewBigUIntFromInt(0))
	return err
}

func staticCallGoContractMethod(vm *PluginVM, callerAddr, contractAddr loom.Address, method string, inpb proto.Message) error {
	input, err := encodeGoCallInput(method, inpb)
	if err != nil {
		return err
	}
	_, err = vm.StaticCall(callerAddr, contractAddr, input)
	return err
}
