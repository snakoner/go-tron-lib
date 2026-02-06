package tron

import (
	"context"
	"encoding/json"
	"math/big"
	"time"
)

const (
	isVisible = true
)

const (
	TxStatusSuccess = "SUCCESS"
	TxStatusFailed  = "FAILED"
	TxStatusPending = "PENDING"
)

const (
	newBlockGenerationTime = 3 * time.Second
	maxSolidBlockWaitTime  = 80 * time.Second
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

func (c *Client) GetTransactionInfoByIDSolid(ctx context.Context, txID string) (Raw, error) {
	var out Raw
	err := c.Call(ctx, "walletsolidity/gettransactioninfobyid", GetTransactionInfoByIDReq{Value: txID}, &out)
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

type CreateTransactionReq struct {
	OwnerAddress string `json:"owner_address"`
	ToAddress    string `json:"to_address"`
	Amount       int64  `json:"amount"`
	Visible      bool   `json:"visible,omitempty"`
}

func (c *Client) CreateTransaction(ctx context.Context, req CreateTransactionReq) (Raw, error) {
	var out Raw
	err := c.Call(ctx, "createtransaction", req, &out)
	return out, err
}

func (c *Client) BuildTransferTRXTx(ctx context.Context, from string, to string, amount *big.Int) (Raw, error) {
	return c.CreateTransaction(ctx, CreateTransactionReq{
		OwnerAddress: from,
		ToAddress:    to,
		Amount:       amount.Int64(),
		Visible:      isVisible,
	})
}

type GetTransactionInfoStatusResult struct {
	BlockNumber int64 `json:"blockNumber"`
	Receipt     struct {
		Result string `json:"result"`
	} `json:"receipt"`
}

func (c *Client) GetTransactionStatusSolid(ctx context.Context, txID string) (string, error) {
	tx, err := c.GetTransactionInfoByIDSolid(ctx, txID)
	if err != nil {
		return "", err
	}

	var out GetTransactionInfoStatusResult
	if err := json.Unmarshal(tx, &out); err != nil {
		return "", err
	}

	return convertTransactionStatus(out), nil
}

func (c *Client) GetTransactionStatus(ctx context.Context, txID string) (string, error) {
	tx, err := c.GetTransactionInfoByID(ctx, txID)
	if err != nil {
		return "", err
	}

	var out GetTransactionInfoStatusResult
	if err := json.Unmarshal(tx, &out); err != nil {
		return "", err
	}

	return convertTransactionStatus(out), nil
}

func convertTransactionStatus(result GetTransactionInfoStatusResult) string {
	if result.BlockNumber == 0 {
		return TxStatusPending
	}

	resultStr := result.Receipt.Result
	if resultStr == "" || resultStr == TxStatusSuccess {
		return TxStatusSuccess
	}

	return TxStatusFailed
}

func (c *Client) WaitForStatusSuccessSolid(ctx context.Context, txID string) (string, error) {
	return waitForStatus(ctx, txID, c.GetTransactionStatusSolid, maxSolidBlockWaitTime)
}

func (c *Client) WaitForStatusSuccess(ctx context.Context, txID string, maxWaitTime time.Duration) (string, error) {
	return waitForStatus(ctx, txID, c.GetTransactionStatus, maxWaitTime)
}

func waitForStatus(ctx context.Context, txID string, statusFunc func(ctx context.Context, txID string) (string, error), maxWaitTime time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, maxWaitTime)
	defer cancel()

	ticker := time.NewTicker(newBlockGenerationTime)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return TxStatusFailed, ctx.Err()
		case <-ticker.C:
			status, err := statusFunc(ctx, txID)
			if err != nil {
				return "", err
			}
			switch status {
			case TxStatusSuccess,
				TxStatusFailed:
				return status, nil
			case TxStatusPending:
				continue
			}
		}
	}
}
