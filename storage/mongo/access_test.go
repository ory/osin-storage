package mongo

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestNewMongoAccessData(t *testing.T) {
    a := osinAccessDataMock
    m := NewMongoAccessData(&a)

    assert.Equal(t, a.AccessToken, m.AccessToken)
    assert.Equal(t, a.RefreshToken, m.RefreshToken)
    assert.Equal(t, a.ExpiresIn, m.ExpiresIn)
    assert.Equal(t, a.Scope, m.Scope)
    assert.Equal(t, a.RedirectUri, m.RedirectUri)
    assert.Equal(t, a.CreatedAt, m.CreatedAt)
    assert.Equal(t, a.UserData, m.UserData)

    assert.NotNil(t, m.AccessData)
    assert.NotNil(t, m.AuthorizeData)
    assert.NotNil(t, m.AccessData)
    assert.NotNil(t, m.Client)
}

func TestToAccessData(t *testing.T) {
    o := osinAccessDataMock
    m := NewMongoAccessData(&o)
    a := m.AccessData()

    assert.Equal(t, a.AccessToken, m.AccessToken)
    assert.Equal(t, a.RefreshToken, m.RefreshToken)
    assert.Equal(t, a.ExpiresIn, m.ExpiresIn)
    assert.Equal(t, a.Scope, m.Scope)
    assert.Equal(t, a.RedirectUri, m.RedirectUri)
    assert.Equal(t, a.CreatedAt, m.CreatedAt)
    assert.Equal(t, a.UserData, m.UserData)

    assert.NotNil(t, m.AccessData)
    assert.NotNil(t, m.AuthorizeData)
    assert.NotNil(t, m.AccessData)
    assert.NotNil(t, m.Client)
}
