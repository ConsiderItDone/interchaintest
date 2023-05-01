package net

import (
	"context"
	"errors"
	"io"
	"strconv"

	"github.com/ethereum/go-ethereum/rpc"
)

type (
	InfoIsBootsrappedRequest struct {
		Chain string `json:"chain"`
	}
	InfoIsBootsrappedResponse struct {
		IsBootstrapped bool `json:"isBootstrapped"`
	}
	PlatformGetHeightResponse struct {
		Height string `json:"height"`
	}
	PlatformCreateSubnetRequest struct {
		Username    string   `json:"string"`
		Password    string   `json:"password"`
		ControlKeys []string `json:"controlKeys"`
		Threshold   int32    `json:"threshold"`
		From        []string `json:"from"`
		ChangeAddr  string   `json:"changeAddr"`
	}
	PlatformCreateSubnetResponse struct {
		TxID string `json:"txID"`
	}
)

func get[RES any](ctx context.Context, addr, method string) (*RES, error) {
	client, err := rpc.Dial(addr)
	if err != nil {
		return nil, err
	}

	var result RES
	return &result, client.CallContext(ctx, &result, method)
}

func call[REQ any, RES any](ctx context.Context, addr, method string, input REQ) (*RES, error) {
	client, err := rpc.Dial(addr)
	if err != nil {
		return nil, err
	}

	var result RES
	return &result, client.CallContext(ctx, &result, method, input)
}

func InfoIsBootstrapped(ctx context.Context, addr, chain string) (bool, error) {
	data, err := call[InfoIsBootsrappedRequest, InfoIsBootsrappedResponse](
		ctx,
		addr,
		"info.isBootstrapped",
		InfoIsBootsrappedRequest{Chain: chain},
	)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return false, nil
		}
		return false, err
	}
	return data.IsBootstrapped, nil
}

func PlatformGetHeight(ctx context.Context, addr string) (uint64, error) {
	data, err := get[PlatformGetHeightResponse](ctx, addr, "platform.getHeight")
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(data.Height, 10, 64)
}

func PlatformCreateSubnet(ctx context.Context, addr string, input *PlatformCreateSubnetRequest) (string, error) {
	output, err := call[PlatformCreateSubnetRequest, PlatformCreateSubnetResponse](
		ctx,
		addr,
		"platform.createSubnet",
		*input,
	)
	if err != nil {
		return "", err
	}
	return output.TxID, nil
}
