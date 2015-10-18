package storage

import "github.com/RangelReale/osin"

type Storage interface {
	osin.Storage
	CreateClient(id, secret, redirectURI string) (osin.Client, error)
	UpdateClient(id, secret, redirectURI string) (osin.Client, error)
}
