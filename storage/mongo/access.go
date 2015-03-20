package mongo

import (
    "github.com/RangelReale/osin"
    "gopkg.in/mgo.v2"
    "time"
)

type MongoAccessData struct {
    Client              *MongoClient
    AuthorizeData       *MongoAuthorizeData
    PreviousAccessData  *MongoAccessData
    AccessToken         string
    RefreshToken        string
    ExpiresIn           int32
    Scope               string
    RedirectUri         string
    CreatedAt           time.Time
    UserData            interface{}
}

var MongoAccessCollectionIndex = mgo.Index{
    Key:        []string{"accessToken"},
    DropDups:   true,
    Background: true,
    Sparse:     true,
}

func NewMongoAccessData(a *osin.AccessData) *MongoAccessData {
    if a == nil {
        return new(MongoAccessData)
    }

    return &MongoAccessData{
        Client: NewMongoClient(a.Client),
        AuthorizeData: NewMongoAuthorizeData(a.AuthorizeData),
        PreviousAccessData: NewMongoAccessData(a.AccessData),
        AccessToken: a.AccessToken,
        RefreshToken: a.RefreshToken,
        ExpiresIn: a.ExpiresIn,
        Scope: a.Scope,
        RedirectUri: a.RedirectUri,
        CreatedAt: a.CreatedAt,
        UserData: a.UserData,
    }
}

func (a *MongoAccessData) AccessData() *osin.AccessData {
    if a == nil {
        return new(osin.AccessData)
    }

    return &osin.AccessData{
        Client: a.Client,
        AuthorizeData: a.AuthorizeData.AuthorizeData(),
        AccessData: a.PreviousAccessData.AccessData(),
        AccessToken: a.AccessToken,
        RefreshToken: a.RefreshToken,
        ExpiresIn: a.ExpiresIn,
        Scope: a.Scope,
        RedirectUri: a.RedirectUri,
        CreatedAt: a.CreatedAt,
        UserData: a.UserData,
    }
}
