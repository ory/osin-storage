package mongo

import (
    "github.com/RangelReale/osin"
    "gopkg.in/mgo.v2"
)

var MongoClientIndex = mgo.Index{
    Key:        []string{"clientid"},
    DropDups:   true,
    Background: true,
    Sparse:     true,
}

// DefaultClient stores all data in struct variables
type MongoClient struct {
    Id          string `bson:"id"`
    Secret      string `bson:"secret"`
    RedirectUri string `bson:"redirectUri"`
    UserData    interface{}
}

func NewMongoClient(c osin.Client) *MongoClient {
    mc := &MongoClient{
        Id: c.GetId(),
        Secret: c.GetSecret(),
        RedirectUri: c.GetRedirectUri(),
        UserData: c.GetUserData(),
    }
    return mc
}

func (d *MongoClient) GetId() string {
    return d.Id
}

func (d *MongoClient) GetSecret() string {
    return d.Secret
}

func (d *MongoClient) GetRedirectUri() string {
    return d.RedirectUri
}

func (d *MongoClient) GetUserData() interface{} {
    return d.UserData
}

func (d *MongoClient) CopyFrom(client osin.Client) {
    d.Id = client.GetId()
    d.Secret = client.GetSecret()
    d.RedirectUri = client.GetRedirectUri()
    d.UserData = client.GetUserData()
}
