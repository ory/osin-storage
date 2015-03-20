# osin-storage

[![Build Status](https://travis-ci.org/ory-platform/osin-storage.svg)](https://travis-ci.org/ory-platform/osin-storage)

Different storage backends for [osin oauth2](https://github.com/RangelReale/osin).
Currently only supporting MongoDB.

Additional to implementing the `osin.Storage` interface, the `OAuthStorage` interface
adds the signature `SetClient(id string, client osin.Client) error` for adding clients.

## Usage

```
go get "github.com/ory-platform/osin-storage/storage/mongo"
```

```go
import (
    "github.com/ory-platform/osin-storage/storage/mongo"
    "gopkg.in/mgo.v2"
)

func main() {
    mgoSession, err := mgo.Dial("localhost")

    if err != nil {
        return nil, err
    }

    defer mgoSession.Close()

    oauthServer, oauthStorage := mongo.NewOAuthServer(mongoSession, conf, "oauthdb")

    // See the osin documentation for more information
    // e.g.: oauthServer.HandleAuthorizeRequest(resp, r)
}
```

## To be done

* Write tests for `storage/mongo/oauth.go`
* Add additional storage back ends