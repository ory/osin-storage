package mongo

import (
    "gopkg.in/mgo.v2"
)

var MongoRefreshTokenIndex = mgo.Index{
    Key:        []string{"token"},
    DropDups:   true,
    Background: true,
    Sparse:     true,
}

type MongoRefreshToken struct {
    Token   string
}
