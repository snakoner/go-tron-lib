package tron

import (
	"context"
	"encoding/json"
)

type Raw = json.RawMessage

func (c *Client) GetNowBlock(ctx context.Context) (Raw, error) {
	var out Raw
	err := c.Call(ctx, "getnowblock", nil, &out)
	return out, err
}

type GetBlockByNumReq struct {
	Num int64 `json:"num"`
}

func (c *Client) GetBlockByNum(ctx context.Context, num int64) (Raw, error) {
	var out Raw
	err := c.Call(ctx, "getblockbynum", GetBlockByNumReq{Num: num}, &out)
	return out, err
}

type GetBlockByIDReq struct {
	Value string `json:"value"`
}

func (c *Client) GetBlockByID(ctx context.Context, blockID string) (Raw, error) {
	var out Raw
	err := c.Call(ctx, "getblockbyid", GetBlockByIDReq{Value: blockID}, &out)
	return out, err
}

type GetTransactionByIDReq struct {
	Value string `json:"value"`
}

func (c *Client) GetTransactionByID(ctx context.Context, txID string) (Raw, error) {
	var out Raw
	err := c.Call(ctx, "gettransactionbyid", GetTransactionByIDReq{Value: txID}, &out)
	return out, err
}

type GetTransactionInfoByIDReq struct {
	Value string `json:"value"`
}

func (c *Client) GetTransactionInfoByID(ctx context.Context, txID string) (Raw, error) {
	var out Raw
	err := c.Call(ctx, "gettransactioninfobyid", GetTransactionInfoByIDReq{Value: txID}, &out)
	return out, err
}

type GetAccountReq struct {
	Address string `json:"address"`
	Visible bool   `json:"visible,omitempty"`
}

func (c *Client) GetAccount(ctx context.Context, address string, visible bool) (Raw, error) {
	var out Raw
	err := c.Call(ctx, "getaccount", GetAccountReq{
		Address: address,
		Visible: visible,
	}, &out)
	return out, err
}

type TriggerConstantContractReq struct {
	OwnerAddress    string `json:"owner_address"`
	ContractAddress string `json:"contract_address"`
	Function        string `json:"function_selector"`
	Parameter       string `json:"parameter,omitempty"`
	CallValue       int64  `json:"call_value,omitempty"`
	FeeLimit        int64  `json:"fee_limit,omitempty"`
	Visible         bool   `json:"visible,omitempty"`
}

func (c *Client) TriggerConstantContract(ctx context.Context, req TriggerConstantContractReq) (Raw, error) {
	var out Raw
	err := c.Call(ctx, "triggerconstantcontract", req, &out)
	return out, err
}

type TriggerSmartContractReq struct {
	OwnerAddress    string `json:"owner_address"`
	ContractAddress string `json:"contract_address"`
	Function        string `json:"function_selector"`
	Parameter       string `json:"parameter,omitempty"`
	CallValue       int64  `json:"call_value,omitempty"`
	FeeLimit        int64  `json:"fee_limit,omitempty"`
	Visible         bool   `json:"visible,omitempty"`
}

func (c *Client) TriggerSmartContract(ctx context.Context, req TriggerSmartContractReq) (Raw, error) {
	var out Raw
	err := c.Call(ctx, "triggersmartcontract", req, &out)
	return out, err
}

func (c *Client) GetNodeInfo(ctx context.Context) (Raw, error) {
	var out Raw
	err := c.Call(ctx, "getnodeinfo", nil, &out)
	return out, err
}

type BroadcastResp struct {
	Result  bool   `json:"result"`
	TxID    string `json:"txid,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (c *Client) BroadcastTransaction(ctx context.Context, signedTx []byte) (*BroadcastResp, error) {
	var out BroadcastResp
	if err := c.Call(ctx, "broadcasttransaction", json.RawMessage(signedTx), &out); err != nil {
		return nil, err
	}

	return &out, nil
}
