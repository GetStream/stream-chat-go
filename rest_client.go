package stream_chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/getstream/easyjson"
)

type RestClient interface {
	Get(path string, urlParams url.Values, result easyjson.Unmarshaler) error
	Post(path string, urlParams url.Values, body interface{}, result easyjson.Unmarshaler) error
	Delete(path string, urlParams url.Values, result easyjson.Unmarshaler) error
}

type HTTPDo interface {
	Do(r *http.Request) (*http.Response, error)
}

var MakeRequest = func(client HTTPDo, header http.Header, method, path string, data interface{}, result easyjson.Unmarshaler) (err error) {
	req, err := newRequest(method, path, data)
	if err != nil {
		return err
	}

	req.Header = header

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	return parseResponse(resp, result)
}

func parseResponse(resp *http.Response, result easyjson.Unmarshaler) error {
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode >= 399 {
		msg, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("chat-client: HTTP %s %s status %s: %s", resp.Request.Method, resp.Request.URL, resp.Status, string(msg))
	}

	if result != nil {
		return easyjson.UnmarshalFromReader(resp.Body, result)
	}

	return nil
}

func newRequest(method, _url string, data interface{}) (req *http.Request, err error) {
	var body []byte

	if m, ok := data.(easyjson.Marshaler); ok {
		body, err = easyjson.Marshal(m)
	} else {
		body, err = json.Marshal(data)
	}

	if err != nil {
		return nil, err
	}

	return http.NewRequest(method, _url, bytes.NewReader(body))
}
