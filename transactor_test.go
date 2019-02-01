package main

import (
	"testing"
	"encoding/json"
	rpcclient "github.com/tendermint/tendermint/rpc/lib/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmrpc "github.com/tendermint/tendermint/rpc/client"
	"fmt"
	"os"
	"bytes"
	"io/ioutil"
	types "github.com/tendermint/tendermint/rpc/lib/types"
)

func TestMarshall(t *testing.T) {
	paramsJSON, _ := json.Marshal(map[string]interface{}{"tx": "asdasdasd"})
	t.Log(string(paramsJSON))

	c := rpcclient.NewJSONRPCClient("localhost:26657")




	request, err := types.MapToRequest(c.cdc, types.JSONRPCStringID("jsonrpc-client"), method, params)
	if err != nil {
		return nil, err
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	// log.Info(string(requestBytes))
	requestBuf := bytes.NewBuffer(requestBytes)
	// log.Info(Fmt("RPC request to %v (%v): %v", c.remote, method, string(requestBytes)))
	httpResponse, err := c.client.Post(c.address, "text/json", requestBuf)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close() // nolint: errcheck

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)

	result := &ctypes.ResultStatus{}
	x, err := rpccli.Call("status", map[string]interface{}{}, result)
	t.Log("err is : ", err)
	t.Log(x)
}

func TestTransacter_Start(t *testing.T) {
	startTransacters(
		[]string{"localhost:26657"},
		1,
		40,
		250,
		"broadcast_tx_"+"async",
	)
}

func latestBlockHeight2(client tmrpc.Client) int64 {
	status, err := client.Status()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return status.SyncInfo.LatestBlockHeight
}
