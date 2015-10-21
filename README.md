# osin-storage

[![Build Status](https://travis-ci.org/ory-am/osin-storage.svg)](https://travis-ci.org/ory-am/osin-storage)

A postgres storage backend for [osin oauth2](https://github.com/RangelReale/osin).

Additional to implementing the `osin.Storage` interface, the `github.com/ory-am/osin-storage/storage.Storage` interface adds
```
CreateClient(id, secret, redirectURI string) (osin.Client, error)
UpdateClient(id, secret, redirectURI string) (osin.Client, error)
```
to the signature.

## Usage

First, install this library with `go get "github.com/ory-am/osin-storage/storage/postgres"`.

```go
import (
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/ory-am/osin-storage/storage/postgres"
	"github.com/RangelReale/osin"
)

func main() {
    // url := "postgres://my-postgres-url/database"
	db, err = sql.Open("postgres", url)
    if err != nil {
        return nil, err
    }

	store := postgres.New(db)
    server := osin.NewServer(osin.NewServerConfig(), store)

    // See the osin documentation for more information
    // e.g.: server.HandleAuthorizeRequest(resp, r)
}
```