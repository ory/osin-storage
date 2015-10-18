package postgres

import (
	"database/sql"
	"fmt"
	"github.com/RangelReale/osin"
	_ "github.com/lib/pq"
	"github.com/ory-am/common/rand/sequence"
	"github.com/ory-am/dockertest"
	"github.com/ory-am/osin-storage/storage"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func TestClientOperations(t *testing.T) {
	container, db := connect(t)
	defer func() {
		require.Nil(t, container.KillRemove())
	}()
	store := New(db)
	require.Nil(t, createAndUseTestDB(db))
	require.Nil(t, store.CreateSchemas())

	create := &osin.DefaultClient{"1", "secret", "http://localhost/", nil}
	createClient(t, store, create)
	getClient(t, store, create)

	update := &osin.DefaultClient{"1", "secret", "http://www.google.com/", nil}
	updateClient(t, store, update)
	getClient(t, store, update)
}

func TestAuthorizeOperations(t *testing.T) {
	container, db := connect(t)
	defer func() {
		require.Nil(t, container.KillRemove())
	}()
	store := New(db)
	require.Nil(t, createAndUseTestDB(db))
	require.Nil(t, store.CreateSchemas())

	client := &osin.DefaultClient{"123", "secret", "http://localhost/", nil}
	createClient(t, store, client)

	authorize := &osin.AuthorizeData{
		Client:      client,
		Code:        "code",
		ExpiresIn:   int32(60),
		Scope:       "scope",
		RedirectUri: "http://localhost/",
		State:       "state",
		// FIXME this should be time.Now(), but an upstream ( https://github.com/lib/pq/issues/329 ) issue prevents this.
		CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		UserData:  nil,
	}

	require.Nil(t, store.SaveAuthorize(authorize))

	result, err := store.LoadAuthorize(authorize.Code)
	require.Nil(t, err)
	require.True(t, reflect.DeepEqual(authorize, result))

	require.Nil(t, store.RemoveAuthorize(authorize.Code))
	_, err = store.LoadAuthorize(authorize.Code)
	require.NotNil(t, err)
}

func TestAccessOperations(t *testing.T) {
	container, db := connect(t)
	defer func() {
		require.Nil(t, container.KillRemove())
	}()
	store := New(db)
	require.Nil(t, createAndUseTestDB(db))
	require.Nil(t, store.CreateSchemas())

	client := &osin.DefaultClient{"123", "secret", "http://localhost/", nil}
	authorize := &osin.AuthorizeData{
		Client:      client,
		Code:        "code",
		ExpiresIn:   int32(60),
		Scope:       "scope",
		RedirectUri: "http://localhost/",
		State:       "state",
		// FIXME this should be time.Now(), but an upstream ( https://github.com/lib/pq/issues/329 ) issue prevents this.
		CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		UserData:  nil,
	}
	nestedAccess := &osin.AccessData{
		Client:        client,
		AuthorizeData: authorize,
		AccessData:    nil,
		AccessToken:   "previous_access",
		RefreshToken:  "previous_refresh",
		ExpiresIn:     int32(60),
		Scope:         "scope",
		RedirectUri:   "https://localhost/",
		// FIXME this should be time.Now(), but an upstream ( https://github.com/lib/pq/issues/329 ) issue prevents this.
		CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		UserData:  nil,
	}
	access := &osin.AccessData{
		Client:        client,
		AuthorizeData: authorize,
		AccessData:    nestedAccess,
		AccessToken:   "access",
		RefreshToken:  "refresh",
		ExpiresIn:     int32(60),
		Scope:         "scope",
		RedirectUri:   "https://localhost/",
		// FIXME this should be time.Now(), but an upstream ( https://github.com/lib/pq/issues/329 ) issue prevents this.
		CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		UserData:  nil,
	}

	createClient(t, store, client)
	require.Nil(t, store.SaveAuthorize(authorize))
	require.Nil(t, store.SaveAccess(nestedAccess))
	require.Nil(t, store.SaveAccess(access))

	result, err := store.LoadAccess(access.AccessToken)
	require.Nil(t, err)
	require.True(t, reflect.DeepEqual(access, result))

	require.Nil(t, store.RemoveAccess(access.AccessToken))
	_, err = store.LoadAccess(access.AccessToken)
	require.NotNil(t, err)
}

func TestRefreshOperations(t *testing.T) {
	container, db := connect(t)
	defer func() {
		require.Nil(t, container.KillRemove())
	}()
	store := New(db)
	require.Nil(t, createAndUseTestDB(db))
	require.Nil(t, store.CreateSchemas())

	client := &osin.DefaultClient{"123", "secret", "http://localhost/", nil}
	authorize := &osin.AuthorizeData{
		Client:      client,
		Code:        "code",
		ExpiresIn:   int32(60),
		Scope:       "scope",
		RedirectUri: "http://localhost/",
		State:       "state",
		CreatedAt:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		UserData:    nil,
	}
	access := &osin.AccessData{
		Client:        client,
		AuthorizeData: authorize,
		AccessData:    nil,
		AccessToken:   "access",
		RefreshToken:  "refresh",
		ExpiresIn:     int32(60),
		Scope:         "scope",
		RedirectUri:   "https://localhost/",
		CreatedAt:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		UserData:      nil,
	}

	createClient(t, store, client)
	require.Nil(t, store.SaveAuthorize(authorize))
	require.Nil(t, store.SaveAccess(access))
	result, err := store.LoadRefresh(access.RefreshToken)
	require.Nil(t, err)

	require.True(t, reflect.DeepEqual(access, result))
	require.Nil(t, store.RemoveRefresh(access.RefreshToken))
	_, err = store.LoadRefresh(access.RefreshToken)
	require.NotNil(t, err)
	require.Nil(t, store.RemoveAccess(access.AccessToken))
	require.Nil(t, store.SaveAccess(access))
	_, err = store.LoadRefresh(access.RefreshToken)
	require.Nil(t, err)
	require.Nil(t, store.RemoveAccess(access.AccessToken))
	_, err = store.LoadRefresh(access.RefreshToken)
	require.NotNil(t, err)
}

func getClient(t *testing.T, store storage.Storage, set osin.Client) {
	client, err := store.GetClient(set.GetId())
	require.Nil(t, err)
	require.EqualValues(t, set, client)
}

func createClient(t *testing.T, store storage.Storage, set osin.Client) {
	client, err := store.CreateClient(set.GetId(), set.GetSecret(), set.GetRedirectUri())
	require.Nil(t, err)
	require.EqualValues(t, set, client)
}

func updateClient(t *testing.T, store storage.Storage, set osin.Client) {
	client, err := store.UpdateClient(set.GetId(), set.GetSecret(), set.GetRedirectUri())
	require.Nil(t, err)
	require.EqualValues(t, set, client)
}

func connect(t *testing.T) (c dockertest.ContainerID, db *sql.DB) {
	c, ip, port, err := dockertest.SetupPostgreSQLContainer()
	require.Nil(t, err)
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=disable", dockertest.PostgresUsername, dockertest.PostgresPassword, ip, port)
	db, err = sql.Open("postgres", url)
	require.Nil(t, err)
	require.Nil(t, db.Ping())
	require.Nil(t, createAndUseTestDB(db))
	return
}

func createAndUseTestDB(db *sql.DB) error {
	name, err := sequence.RuneSequence(10, []rune("abcdefghijklmnopqrstuvwxyz"))
	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", string(name)))
	return err
}
