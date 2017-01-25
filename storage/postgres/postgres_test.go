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
	"github.com/ory-am/common/pkg"
	"github.com/ory/osin-storage/storage"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/ory-am/dockertest.v2"
)

var db *sql.DB
var store *Storage
var userDataMock = "bar"

func TestMain(m *testing.M) {
	c, err := dockertest.ConnectToPostgreSQL(15, time.Second, func(url string) bool {
		var err error
		db, err = sql.Open("postgres", url)
		if err != nil {
			return false
		}
		return db.Ping() == nil
	})

	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	store = New(db)
	if err = store.CreateSchemas(); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	retCode := m.Run()

	// force teardown
	tearDown(c)

	os.Exit(retCode)
}

func tearDown(c dockertest.ContainerID) {
	c.KillRemove()
}

func TestClientOperations(t *testing.T) {
	create := &osin.DefaultClient{Id: "1", Secret: "secret", RedirectUri: "http://localhost/", UserData: ""}
	createClient(t, store, create)
	getClient(t, store, create)

	update := &osin.DefaultClient{Id: "1", Secret: "secret123", RedirectUri: "http://www.google.com/", UserData: "{}"}
	updateClient(t, store, update)
	getClient(t, store, update)
}

func TestAuthorizeOperations(t *testing.T) {
	client := &osin.DefaultClient{Id: "2", Secret: "secret", RedirectUri: "http://localhost/", UserData: ""}
	createClient(t, store, client)

	for k, authorize := range []*osin.AuthorizeData{
		{
			Client:      client,
			Code:        uuid.New(),
			ExpiresIn:   int32(600),
			Scope:       "scope",
			RedirectUri: "http://localhost/",
			State:       "state",
			CreatedAt: time.Now().Round(time.Second),
			UserData:  userDataMock,
		},
	} {
		// Test save
		require.Nil(t, store.SaveAuthorize(authorize))

		// Test fetch
		result, err := store.LoadAuthorize(authorize.Code)
		require.Nil(t, err)
		require.Equal(t, authorize.CreatedAt.Unix(), authorize.CreatedAt.Unix())
		authorize.CreatedAt = result.CreatedAt
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
		CreatedAt: time.Now().Round(time.Second),
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
		CreatedAt: time.Now().Round(time.Second),
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
		CreatedAt: time.Now().Round(time.Second),
		UserData:  userDataMock,
	}

	createClient(t, store, client)
	require.Nil(t, store.SaveAuthorize(authorize))
	require.Nil(t, store.SaveAccess(nestedAccess))
	require.Nil(t, store.SaveAccess(access))

	result, err := store.LoadAccess(access.AccessToken)
	require.Nil(t, err)
	require.Equal(t, access.CreatedAt.Unix(), result.CreatedAt.Unix())
	require.Equal(t, access.AccessData.CreatedAt.Unix(), result.AccessData.CreatedAt.Unix())
	require.Equal(t, access.AuthorizeData.CreatedAt.Unix(), result.AuthorizeData.CreatedAt.Unix())
	access.CreatedAt = result.CreatedAt
	access.AccessData.CreatedAt = result.AccessData.CreatedAt
	access.AuthorizeData.CreatedAt = result.AuthorizeData.CreatedAt
	require.Equal(t, access, result)

	require.Nil(t, store.RemoveAuthorize(authorize.Code))
	_, err = store.LoadAccess(access.AccessToken)
	require.Nil(t, err)

	require.Nil(t, store.RemoveAccess(nestedAccess.AccessToken))
	_, err = store.LoadAccess(access.AccessToken)
	require.Nil(t, err)

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
					CreatedAt:    time.Now().Round(time.Second),
					UserData:    userDataMock,
				},
				AccessData:   nil,
				AccessToken:  uuid.New(),
				RefreshToken: uuid.New(),
				ExpiresIn:    int32(60),
				Scope:        "scope",
				RedirectUri:  "https://localhost/",
				CreatedAt:     time.Now().Round(time.Second),
				UserData:     userDataMock,
			},
		},
	} {
		createClient(t, store, client)
		require.Nil(t, store.SaveAuthorize(c.access.AuthorizeData), "Case %d", k)
		require.Nil(t, store.SaveAccess(c.access), "Case %d", k)

		result, err := store.LoadRefresh(c.access.RefreshToken)
		require.Nil(t, err)
		require.Equal(t, c.access.CreatedAt.Unix(), result.CreatedAt.Unix())
		require.Equal(t, c.access.AuthorizeData.CreatedAt.Unix(), result.AuthorizeData.CreatedAt.Unix())
		c.access.CreatedAt = result.CreatedAt
		c.access.AuthorizeData.CreatedAt = result.AuthorizeData.CreatedAt
		require.Equal(t, c.access, result, "Case %d", k)

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

func TestErrors(t *testing.T) {
	assert.Nil(t, store.CreateClient(&osin.DefaultClient{Id: "dupe"}))
	assert.NotNil(t, store.CreateClient(&osin.DefaultClient{Id: "dupe"}))
	assert.NotNil(t, store.CreateClient(&osin.DefaultClient{Id: "foo", UserData: struct{}{}}))
	assert.NotNil(t, store.SaveAccess(&osin.AccessData{AccessToken: "", AccessData: &osin.AccessData{}, AuthorizeData: &osin.AuthorizeData{}}))
	assert.Nil(t, store.SaveAuthorize(&osin.AuthorizeData{Code: "a", Client: &osin.DefaultClient{}}))
	assert.NotNil(t, store.SaveAuthorize(&osin.AuthorizeData{Code: "a", Client: &osin.DefaultClient{}}))
	assert.NotNil(t, store.SaveAuthorize(&osin.AuthorizeData{Code: "b", Client: &osin.DefaultClient{}, UserData: struct{}{}}))
	_, err := store.LoadAccess("")
	assert.Equal(t, pkg.ErrNotFound, err)
	_, err = store.LoadAuthorize("")
	assert.Equal(t, pkg.ErrNotFound, err)
	_, err = store.LoadRefresh("")
	assert.Equal(t, pkg.ErrNotFound, err)
	_, err = store.GetClient("")
	assert.Equal(t, pkg.ErrNotFound, err)
}

type ts struct{}

func (s *ts) String() string {
	return "foo"
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

	res, err = assertToString(new(ts))
	assert.Nil(t, err)
	assert.Equal(t, "foo", res)
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
