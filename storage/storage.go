package storage

import "github.com/RangelReale/osin"

type OAuthStorage interface {
    osin.Storage
    SetClient(id string, client osin.Client) error
}