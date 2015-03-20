package mongo

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/ory-platform/dockertest"
    "github.com/ory-platform/common/mgopath"
    "github.com/ory-platform/common/rand/sequence"
    "fmt"
    "gopkg.in/mgo.v2"
    "github.com/RangelReale/osin")

func TestNewOAuthMongoStorage(t *testing.T) {
    config := osin.NewServerConfig()
    container, db := connect(t)
    defer container.KillRemove(t)

    storage, err := NewOAuthMongoStorage(db, config)
    assert.Nil(t, err)
    assert.NotNil(t, storage)
}

func connect(t *testing.T) (*dockertest.ContainerID, *mgo.Database) {
    containerID, ip, port := dockertest.SetupMongoContainer(t)
    dbName, err := sequence.RuneSequence(22, []rune("abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTUVXYZ1234567890"))
    assert.Nil(t, err)
    path := fmt.Sprintf("mongodb://%s:%d/%s", ip, port, string(dbName))
    db, err := mgopath.Connect(path)
    assert.Nil(t, err)
    return &containerID, db
}
