package Utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// MakeURL encode to full url request
func MakeURL(scheme, host, path string, args url.Values) string {
	u := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}
	if args != nil {
		u.RawQuery = args.Encode()
	}
	return u.String()
}

// HTTPPost do http post with timeout
func HTTPPost(client *http.Client, uri string, data url.Values) (resp []byte, err error) {
	if data == nil || client == nil {
		return nil, nil
	}

	defer func() {
		if er := recover(); er != nil {
			err = errors.New("Client.Timeout")
			resp = nil
		}
	}()

	res, err := client.PostForm(uri, data)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(res.Body)
}

// DoHTTP do http (GET, PUT, POST ...)
func DoHTTP(client *http.Client, req *http.Request) (resp []byte, isTimeout bool, statusCode int, err error) {
	err, isTimeout = nil, false

	if client == nil {
		return
	}

	defer func() {
		if er := recover(); er != nil {
			resp, isTimeout, err = nil, true, fmt.Errorf("%v", er)
		}
	}()

	res, err := client.Do(req)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()

	if err != nil {
		return nil, false, 0, err
	}

	resp, err = ioutil.ReadAll(res.Body)
	statusCode = res.StatusCode
	return
}
