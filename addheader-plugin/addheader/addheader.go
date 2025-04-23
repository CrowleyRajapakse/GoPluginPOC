package main

import (
	"context"

	"github.com/wso2/apk/gateway/enforcer/pkg/plugins"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	extprocv3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
)

type AddHeader struct{}

func (a *AddHeader) Name() string {
	return "AddHeader"
}

func (a *AddHeader) ApplyRequestHeaders(ctx context.Context, req *extprocv3.ProcessingRequest) ([]*corev3.HeaderValueOption, error) {
	return []*corev3.HeaderValueOption{
		{Header: &corev3.HeaderValue{Key: "x-added-req", Value: "header-plugin"}},
	}, nil
}

func (a *AddHeader) ApplyResponseHeaders(ctx context.Context, req *extprocv3.ProcessingRequest) ([]*corev3.HeaderValueOption, error) {
	return []*corev3.HeaderValueOption{
		{Header: &corev3.HeaderValue{Key: "x-added-res", Value: "header-plugin"}},
	}, nil
}

func (a *AddHeader) ApplyRequestBody(ctx context.Context, req *extprocv3.ProcessingRequest) ([]byte, error) {
	return nil, nil
}

func (a *AddHeader) ApplyResponseBody(ctx context.Context, req *extprocv3.ProcessingRequest) ([]byte, error) {
	return nil, nil
}

var Plugin plugins.Policy = &AddHeader{}
