package mongo

import (
	"github.com/RangelReale/osin"
	"gopkg.in/mgo.v2"
	"time"
)

// Authorization data
type MongoAuthorizeData struct {
	Client      *MongoClient
	Code        string
	ExpiresIn   int32
	Scope       string
	RedirectUri string
	State       string
	CreatedAt   time.Time
	UserData    interface{}
}

var MongoAuthorizeCollectionIndex = mgo.Index{
	Key:        []string{"code"},
	DropDups:   true,
	Background: true,
	Sparse:     true,
}

// IsExpired is true if authorization expired
func (d *MongoAuthorizeData) IsExpired() bool {
	return d.IsExpiredAt(time.Now())
}

// IsExpired is true if authorization expires at time 't'
func (d *MongoAuthorizeData) IsExpiredAt(t time.Time) bool {
	return d.ExpireAt().Before(t)
}

// ExpireAt returns the expiration date
func (d *MongoAuthorizeData) ExpireAt() time.Time {
	return d.CreatedAt.Add(time.Duration(d.ExpiresIn) * time.Second)
}

func NewMongoAuthorizeData(a *osin.AuthorizeData) *MongoAuthorizeData {
	if a == nil {
		return new(MongoAuthorizeData)
	}

	return &MongoAuthorizeData{
		Client:      NewMongoClient(a.Client),
		Code:        a.Code,
		ExpiresIn:   a.ExpiresIn,
		Scope:       a.Scope,
		RedirectUri: a.RedirectUri,
		State:       a.State,
		CreatedAt:   a.CreatedAt,
		UserData:    a.UserData,
	}
}

func (a *MongoAuthorizeData) AuthorizeData() *osin.AuthorizeData {
	if a == nil {
		return new(osin.AuthorizeData)
	}

	return &osin.AuthorizeData{
		Client:      a.Client,
		Code:        a.Code,
		ExpiresIn:   a.ExpiresIn,
		Scope:       a.Scope,
		RedirectUri: a.RedirectUri,
		State:       a.State,
		CreatedAt:   a.CreatedAt,
		UserData:    a.UserData,
	}
}
