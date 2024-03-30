package auth

import (
	"go-web-template/pkg/serializer"
	"net/url"
	"time"
)

var (
	ErrAuthFailed        = serializer.NewError(serializer.CodeInvalidSign, "invalid sign", nil)
	ErrAuthHeaderMissing = serializer.NewError(serializer.CodeNoPermission, "authorization header is missing", nil)
	ErrExpiresMissing    = serializer.NewError(serializer.CodeNoPermission, "expire timestamp is missing", nil)
	ErrExpired           = serializer.NewError(serializer.CodeSignTimeout, "signature expired", nil)
)
var General Auth

type Auth interface {
	Sign(body string, expires int64) string
	Check(body string, sign string) error
}

func SignURL(instance Auth, url *url.URL, duration time.Duration) *url.URL {
	expires := int64(0)
	if duration != 0 {
		expires = time.Now().Add(duration).Unix()
	}
	sign := instance.Sign(url.Path, expires)
	quires := url.Query()
	quires.Set("sign", sign)
	url.RawQuery = quires.Encode()
	return url
}

func CheckURL(instance Auth, url *url.URL) error {
	sign := url.Query().Get("sign")
	if sign == "" {
		return ErrAuthHeaderMissing
	}
	quires := url.Query()
	quires.Del("sign")
	url.RawQuery = quires.Encode()
	return instance.Check(url.Path, sign)
}
