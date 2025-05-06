package keycloak

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	UserNotFoundErr = errors.New("keycloak: user not found")
)

type BadResponseErr struct {
	Code     int
	Response string
}

func NewBadResponseErr(resp *http.Response) *BadResponseErr {
	buf := new(strings.Builder)
	io.Copy(buf, resp.Body)

	return &BadResponseErr{
		Code:     resp.StatusCode,
		Response: buf.String(),
	}
}

func (err *BadResponseErr) Error() string {
	return fmt.Sprintf("keycloak: bad response: code: %d; response: %s", err.Code, err.Response)
}
