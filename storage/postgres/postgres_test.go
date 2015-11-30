package postgres

import (
	"database/sql"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/RangelReale/osin"
	_ "github.com/lib/pq"
	"github.com/ory-am/dockertest"
	"github.com/ory-am/osin-storage/storage"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var db *sql.DB
var store *Storage
var userDataMock = "bar"

func TestMain(m *testing.M) {
	var err error
	var c dockertest.ContainerID
	c, db, err = dockertest.OpenPostgreSQLContainerConnection(15, time.Millisecond*500)
	if err != nil {
		log.Fatalf("Could not set up PostgreSQL container: %v", err)
	}
	defer c.KillRemove()

	store = New(db)
	if err = store.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	}

	os.Exit(m.Run())
}

func TestClientOperations(t *testing.T) {
	create := &osin.DefaultClient{Id: "1", Secret: "secret", RedirectUri: "http://localhost/", UserData: ""}
	createClient(t, store, create)
	getClient(t, store, create)

	update := &osin.DefaultClient{Id: "1", Secret: "secret", RedirectUri: "http://www.google.com/", UserData: ""}
	updateClient(t, store, update)
	getClient(t, store, update)

	assert.NotNil(t, store.CreateClient(&osin.DefaultClient{Id: "1", Secret: "secret", RedirectUri: "http://www.google.com/", UserData: struct{}{}}))
}

func TestAuthorizeOperations(t *testing.T) {
	client := &osin.DefaultClient{Id: "2", Secret: "secret", RedirectUri: "http://localhost/", UserData: ""}
	createClient(t, store, client)

	for k, authorize := range []*osin.AuthorizeData{
		{
			Client:      client,
			Code:        uuid.New(),
			ExpiresIn:   int32(60),
			Scope:       "scope",
			RedirectUri: "http://localhost/",
			State:       "state",
			// FIXME this should be time.Now(), but an upstream ( https://github.com/lib/pq/issues/329 ) issue prevents this.
			CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
			UserData:  userDataMock,
		},
	} {
		// Test save
		require.Nil(t, store.SaveAuthorize(authorize))

		// Test fetch
		result, err := store.LoadAuthorize(authorize.Code)
		require.Nil(t, err)
		require.True(t, reflect.DeepEqual(authorize, result), "Case: %d\n%v\n\n%v", k, authorize, result)

		// Test remove
		require.Nil(t, store.RemoveAuthorize(authorize.Code))
		_, err = store.LoadAuthorize(authorize.Code)
		require.NotNil(t, err)
	}

	removeClient(t, store, client)
}

func TestStoreFailsOnInvalidUserData(t *testing.T) {
	client := &osin.DefaultClient{Id: "3", Secret: "secret", RedirectUri: "http://localhost/", UserData: ""}
	authorize := &osin.AuthorizeData{
		Client:      client,
		Code:        uuid.New(),
		ExpiresIn:   int32(60),
		Scope:       "scope",
		RedirectUri: "http://localhost/",
		State:       "state",
		CreatedAt:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		UserData:    struct{ foo string }{"bar"},
	}
	access := &osin.AccessData{
		Client:        client,
		AuthorizeData: authorize,
		AccessData:    nil,
		AccessToken:   uuid.New(),
		RefreshToken:  uuid.New(),
		ExpiresIn:     int32(60),
		Scope:         "scope",
		RedirectUri:   "https://localhost/",
		CreatedAt:     time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		UserData:      struct{ foo string }{"bar"},
	}
	assert.NotNil(t, store.SaveAuthorize(authorize))
	assert.NotNil(t, store.SaveAccess(access))
}

func TestAccessOperations(t *testing.T) {
	client := &osin.DefaultClient{Id: "3", Secret: "secret", RedirectUri: "http://localhost/", UserData: ""}
	authorize := &osin.AuthorizeData{
		Client:      client,
		Code:        uuid.New(),
		ExpiresIn:   int32(60),
		Scope:       "scope",
		RedirectUri: "http://localhost/",
		State:       "state",
		// FIXME this should be time.Now(), but an upstream ( https://github.com/lib/pq/issues/329 ) issue prevents this.
		CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		UserData:  userDataMock,
	}
	nestedAccess := &osin.AccessData{
		Client:        client,
		AuthorizeData: authorize,
		AccessData:    nil,
		AccessToken:   uuid.New(),
		RefreshToken:  uuid.New(),
		ExpiresIn:     int32(60),
		Scope:         "scope",
		RedirectUri:   "https://localhost/",
		// FIXME this should be time.Now(), but an upstream ( https://github.com/lib/pq/issues/329 ) issue prevents this.
		CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		UserData:  userDataMock,
	}
	access := &osin.AccessData{
		Client:        client,
		AuthorizeData: authorize,
		AccessData:    nestedAccess,
		AccessToken:   uuid.New(),
		RefreshToken:  uuid.New(),
		ExpiresIn:     int32(60),
		Scope:         "scope",
		RedirectUri:   "https://localhost/",
		// FIXME this should be time.Now(), but an upstream ( https://github.com/lib/pq/issues/329 ) issue prevents this.
		CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		UserData:  userDataMock,
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
	require.Nil(t, store.RemoveAuthorize(authorize.Code))

	removeClient(t, store, client)
}

func TestRefreshOperations(t *testing.T) {
	client := &osin.DefaultClient{Id: "4", Secret: "secret", RedirectUri: "http://localhost/", UserData: ""}
	type test struct {
		access *osin.AccessData
	}

	for k, c := range []*test{
		{
			access: &osin.AccessData{
				Client: client,
				AuthorizeData: &osin.AuthorizeData{
					Client:      client,
					Code:        uuid.New(),
					ExpiresIn:   int32(60),
					Scope:       "scope",
					RedirectUri: "http://localhost/",
					State:       "state",
					CreatedAt:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
					UserData:    userDataMock,
				},
				AccessData:   nil,
				AccessToken:  uuid.New(),
				RefreshToken: uuid.New(),
				ExpiresIn:    int32(60),
				Scope:        "scope",
				RedirectUri:  "https://localhost/",
				CreatedAt:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				UserData:     userDataMock,
			},
		},
	} {
		createClient(t, store, client)
		require.Nil(t, store.SaveAuthorize(c.access.AuthorizeData), "Case %d", k)
		require.Nil(t, store.SaveAccess(c.access), "Case %d", k)

		result, err := store.LoadRefresh(c.access.RefreshToken)
		require.Nil(t, err)
		require.True(t, reflect.DeepEqual(c.access, result), "Case %d", k)

		require.Nil(t, store.RemoveRefresh(c.access.RefreshToken))
		_, err = store.LoadRefresh(c.access.RefreshToken)

		require.NotNil(t, err, "Case %d", k)
		require.Nil(t, store.RemoveAccess(c.access.AccessToken), "Case %d", k)
		require.Nil(t, store.SaveAccess(c.access), "Case %d", k)

		_, err = store.LoadRefresh(c.access.RefreshToken)
		require.Nil(t, err, "Case %d", k)

		require.Nil(t, store.RemoveAccess(c.access.AccessToken), "Case %d", k)
		_, err = store.LoadRefresh(c.access.RefreshToken)
		require.NotNil(t, err, "Case %d", k)

	}
	removeClient(t, store, client)
}

func TestAssertToString(t *testing.T) {
	res, err := assertToString(struct{}{})
	assert.NotNil(t, err)

	res, err = assertToString("foo")
	assert.Nil(t, err)
	assert.Equal(t, "foo", res)

	res, err = assertToString(nil)
	assert.Nil(t, err)
	assert.Equal(t, "", res)
}

func getClient(t *testing.T, store storage.Storage, set osin.Client) {
	client, err := store.GetClient(set.GetId())
	require.Nil(t, err)
	require.EqualValues(t, set, client)
}

func createClient(t *testing.T, store storage.Storage, set osin.Client) {
	require.Nil(t, store.CreateClient(set))
}

func updateClient(t *testing.T, store storage.Storage, set osin.Client) {
	require.Nil(t, store.UpdateClient(set))
}

func removeClient(t *testing.T, store storage.Storage, set osin.Client) {
	require.Nil(t, store.RemoveClient(set.GetId()))
}
