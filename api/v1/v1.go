package api

import (
	"context"
)

//go:generate ../../scripts/protoc.sh ./v1.proto

// Service is the registrar server interface
//
// This interface is implemented by the server and the rpc client
type Service interface {
	CreateMedia(ctx context.Context, r *CreateMediaRequest) (*Media, error)
	GetMedia(ctx context.Context, r *GetMediaRequest) (*Media, error)
}
