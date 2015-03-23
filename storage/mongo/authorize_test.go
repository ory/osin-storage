package mongo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMongoAuthorizeData(t *testing.T) {
	a := osinAuthorizeDataMock
	m := NewMongoAuthorizeData(&a)

	assert.Equal(t, a.ExpireAt(), m.ExpireAt())
	assert.Equal(t, a.IsExpired(), m.IsExpired())

	assert.Equal(t, a.Code, m.Code)
	assert.Equal(t, a.CreatedAt, m.CreatedAt)
	assert.Equal(t, a.ExpiresIn, m.ExpiresIn)
	assert.Equal(t, a.RedirectUri, m.RedirectUri)
	assert.Equal(t, a.Scope, m.Scope)
	assert.Equal(t, a.State, m.State)
	assert.Equal(t, a.RedirectUri, m.RedirectUri)

	assert.NotNil(t, m.AuthorizeData)
	assert.NotNil(t, m.UserData)
	assert.NotNil(t, m.Client)
}

func TestToAuthorizeData(t *testing.T) {
	o := osinAuthorizeDataMock
	m := NewMongoAuthorizeData(&o)
	a := m.AuthorizeData()

	assert.Equal(t, a.ExpireAt(), m.ExpireAt())
	assert.Equal(t, a.IsExpired(), m.IsExpired())

	assert.Equal(t, a.Code, m.Code)
	assert.Equal(t, a.CreatedAt, m.CreatedAt)
	assert.Equal(t, a.ExpiresIn, m.ExpiresIn)
	assert.Equal(t, a.RedirectUri, m.RedirectUri)
	assert.Equal(t, a.Scope, m.Scope)
	assert.Equal(t, a.State, m.State)
	assert.Equal(t, a.RedirectUri, m.RedirectUri)

	assert.NotNil(t, m.AuthorizeData)
	assert.NotNil(t, m.UserData)
	assert.NotNil(t, m.Client)
}
