package transporttest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewClient(decorators ...Decorator) *http.Client {
	return &http.Client{
		Transport: Decorate(http.DefaultTransport, decorators...),
	}
}

type Interceptor func(req *http.Request) (*http.Response, error)

func (f Interceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type Decorator func(http.RoundTripper) http.RoundTripper

func Decorate(transport http.RoundTripper, decorators ...Decorator) http.RoundTripper {
	decorated := transport
	for i := len(decorators) - 1; i >= 0; i-- {
		decorated = decorators[i](decorated)
	}
	return decorated
}

func AssertHost(t *testing.T, want string) Decorator {
	return func(transport http.RoundTripper) http.RoundTripper {
		return Interceptor(func(req *http.Request) (*http.Response, error) {
			if req.URL.Host != want {
				t.Errorf("unexpected host: got %s, want %s", req.URL.Host, want)
			}
			return transport.RoundTrip(req)
		})
	}

}

func AssertMethod(t *testing.T, want string) Decorator {
	return func(transport http.RoundTripper) http.RoundTripper {
		return Interceptor(func(req *http.Request) (*http.Response, error) {
			if req.Method != want {
				t.Errorf("unexpected method: got %s, want %s", req.Method, want)
			}
			return transport.RoundTrip(req)
		})
	}
}

func AssertPath(t *testing.T, want string) Decorator {
	return func(transport http.RoundTripper) http.RoundTripper {
		return Interceptor(func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != want {
				t.Errorf("unexpected path: got %s, want %s", req.URL.Path, want)
			}
			return transport.RoundTrip(req)
		})
	}
}

func RespondJSON(data any, code int) Decorator {
	return func(transport http.RoundTripper) http.RoundTripper {
		return Interceptor(func(req *http.Request) (*http.Response, error) {
			w := httptest.NewRecorder()
			w.WriteHeader(code)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(data)
			return w.Result(), nil
		})
	}
}
