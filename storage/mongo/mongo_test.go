package mongo

import (
    "github.com/RangelReale/osin"
    "time"
)

type UserDataMock struct {
    Username string
    Password string
}

var (
    userDataMock = UserDataMock{
        Username: "foo",
        Password: "bar",
    }
    osinDefaultClientMock = osin.DefaultClient{
        Id:          "1234",
        Secret:      "foo",
        RedirectUri: "bar",
        UserData:     &userDataMock,
    }
    osinAuthorizeDataMock = osin.AuthorizeData{
        Client: &osinDefaultClientMock,
        Code: "code",
        ExpiresIn: 60,
        Scope: "scope",
        RedirectUri: "redirect",
        State: "state",
        CreatedAt: time.Time{},
        UserData: &userDataMock,
    }
    osinAccessDataMock = osin.AccessData{
        Client: &osinDefaultClientMock,
        AuthorizeData: &osinAuthorizeDataMock,
        AccessData: nil,
        AccessToken: "access",
        RefreshToken: "refresh",
        ExpiresIn: 60,
        Scope: "scope",
        RedirectUri: "redirect",
        CreatedAt: time.Time{},
        UserData: &userDataMock,
    }
)
