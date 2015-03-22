package mongo

import (
    "log"

    "github.com/ory-platform/osin-storage/storage"
    "github.com/RangelReale/osin"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

// collection names for the entities
const (
    CLIENT_COLLECTION       = "client"
    AUTHORIZE_COLLECTION    = "authorization"
    ACCESS_COLLECTION       = "access"
    REFRESH_COLLECTION      = "refresh"
)

type MongoStorage struct {
    database *mgo.Database
}

func NewOAuthMongoStorage(database *mgo.Database) (storage.OAuthStorage, error) {
    err := database.C(REFRESH_COLLECTION).EnsureIndex(MongoRefreshTokenIndex)
    if err != nil {
        return nil, err
    }

    err = database.C(CLIENT_COLLECTION).EnsureIndex(MongoClientIndex)
    if err != nil {
        return nil, err
    }

    err = database.C(ACCESS_COLLECTION).EnsureIndex(MongoAccessCollectionIndex)
    if err != nil {
        return nil, err
    }

    err = database.C(AUTHORIZE_COLLECTION).EnsureIndex(MongoAuthorizeCollectionIndex)
    if err != nil {
        return nil, err
    }

    oAuthStorage := &MongoStorage{
        database: database,
    }
    return oAuthStorage, nil
}

func logErr(err error) {
    if err != nil {
        log.Panic(err)
    }
}

func (s *MongoStorage) Clone() osin.Storage {
    return s
}

func (s *MongoStorage) Close() {
    // We do not want to close the session on exception because the service will crash after that
}

func (s *MongoStorage) GetClient(id string) (osin.Client, error) {
    log.Printf("GetClient: %s", id)

    clients := s.database.C(CLIENT_COLLECTION)
    client := new(MongoClient)
    err := clients.Find(bson.M{"id": id}).One(client)
    return client, err
}

func (s *MongoStorage) SetClient(id string, client osin.Client) error {
    log.Printf("SetClient: %s", id)

    clients := s.database.C(CLIENT_COLLECTION)
    _, err := clients.Upsert(bson.M{"id": id}, client)
    return err
}

func (s *MongoStorage) SaveAuthorize(data *osin.AuthorizeData) error {
    log.Printf("SaveAuthorize: %s", data.Code)

    a := NewMongoAuthorizeData(data)
    ac := s.database.C(AUTHORIZE_COLLECTION)
    _, err := ac.Upsert(bson.M{"code": data.Code}, a)
    return err
}

func (s *MongoStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
    log.Printf("LoadAuthorize: %s", code)

    a := new(MongoAuthorizeData)
    ac := s.database.C(AUTHORIZE_COLLECTION)
    err := ac.Find(bson.M{"code": code}).One(a)
    return a.AuthorizeData(), err
}

func (s *MongoStorage) RemoveAuthorize(code string) error {
    log.Printf("RemoveAuthorize: %s", code)

    ac := s.database.C(AUTHORIZE_COLLECTION)
    _, err := ac.RemoveAll(bson.M{"code": code})
    return err
}

func (s *MongoStorage) SaveAccess(data *osin.AccessData) error {
    var err error

    log.Printf("SaveAccess: %s", data.AccessToken)

    a := NewMongoAccessData(data)
    ac := s.database.C(ACCESS_COLLECTION)
    _, err = ac.Upsert(bson.M{"accesstoken": data.AccessToken}, a)

    if err != nil {
        return err
    }

    if data.RefreshToken != "" {
        t := &MongoRefreshToken{Token: data.AccessToken}
        rc := s.database.C(REFRESH_COLLECTION)
        _, err = rc.Upsert(bson.M{"token": data.AccessToken}, t)
    }

    return err
}

func (s *MongoStorage) LoadAccess(code string) (*osin.AccessData, error) {
    log.Printf("LoadAccess: %s", code)

    a := new(MongoAccessData)
    as := s.database.C(ACCESS_COLLECTION)
    err := as.Find(bson.M{"accesstoken": code}).One(a)
    return a.AccessData(), err
}

func (s *MongoStorage) RemoveAccess(code string) error {
    log.Printf("RemoveAccess: %s", code)

    ac := s.database.C(ACCESS_COLLECTION)
    _, err := ac.RemoveAll(bson.M{"accesstoken": code})
    return err
}

func (s *MongoStorage) LoadRefresh(code string) (*osin.AccessData, error) {
    log.Printf("LoadRefresh: %s", code)

    r := new(MongoRefreshToken)
    rc := s.database.C(REFRESH_COLLECTION)
    err := rc.Find(bson.M{"code": code}).One(r)

    if err != nil {
        return &osin.AccessData{}, err
    }

    return s.LoadAccess(r.Token)
}

func (s *MongoStorage) RemoveRefresh(code string) error {
    log.Printf("RemoveRefresh: %s", code)

    rc := s.database.C(REFRESH_COLLECTION)
    _, err := rc.RemoveAll(bson.M{"code": code})

    return err
}
