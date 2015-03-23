package mongo

import (
	"github.com/RangelReale/osin"
	"github.com/stretchr/testify/assert"
	"testing"
)

var copyClientMock = &osin.DefaultClient{
	Id:          "4321",
	Secret:      "foo",
	RedirectUri: "bar",
	UserData:    &userDataMock,
}

func TestNewMongoClient(t *testing.T) {
	u := &osinDefaultClientMock
	m := NewMongoClient(u)

	assert.Equal(t, u.GetId(), m.GetId())
	assert.Equal(t, u.GetRedirectUri(), m.GetRedirectUri())
	assert.Equal(t, u.GetSecret(), m.GetSecret())
	assert.Equal(t, u.GetUserData(), m.GetUserData())
}

func TestCopyClient(t *testing.T) {
	u := &osinDefaultClientMock
	u.CopyFrom(copyClientMock)

	assert.Equal(t, u.GetId(), copyClientMock.GetId())
	assert.Equal(t, u.GetRedirectUri(), copyClientMock.GetRedirectUri())
	assert.Equal(t, u.GetSecret(), copyClientMock.GetSecret())
	assert.Equal(t, u.GetUserData(), copyClientMock.GetUserData())
}
