package httpsource

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver.Server = &Server{}

func TestGet(t *testing.T) {
	listener := createTestHTTPServer(t)
	defer listener.Close()
	params := imageserver.Params{imageserver.SourceParam: createTestURL(listener, testdata.MediumFileName)}
	server := &Server{}
	im, err := server.Get(params)
	if err != nil {
		t.Fatal(err)
	}
	if im == nil {
		t.Fatal("no image")
	}
	if im.Format != testdata.Medium.Format {
		t.Fatalf("unexpected image format: got \"%s\", wanted \"%s\"", im.Format, testdata.Medium.Format)
	}
	if len(im.Data) != len(testdata.Medium.Data) {
		t.Fatalf("unexpected image data length: got %d, wanted %d", len(im.Data), len(testdata.Medium.Data))
	}
}

func TestGetErrorNoSource(t *testing.T) {
	listener := createTestHTTPServer(t)
	defer listener.Close()
	params := imageserver.Params{}
	server := &Server{}
	_, err := server.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorNotFound(t *testing.T) {
	listener := createTestHTTPServer(t)
	defer listener.Close()
	source := createTestURL(listener, testdata.MediumFileName)
	source.Path += "foobar"
	params := imageserver.Params{imageserver.SourceParam: source}
	server := &Server{}
	_, err := server.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorInvalidUrl(t *testing.T) {
	params := imageserver.Params{imageserver.SourceParam: "foobar"}
	server := &Server{}
	_, err := server.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorInvalidUrlScheme(t *testing.T) {
	params := imageserver.Params{imageserver.SourceParam: "custom://foobar"}
	server := &Server{}
	_, err := server.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorRequest(t *testing.T) {
	params := imageserver.Params{imageserver.SourceParam: "http://localhost:123456"}
	server := &Server{}
	_, err := server.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
}

type errorReadCloser struct{}

func (erc *errorReadCloser) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

func (erc *errorReadCloser) Close() error {
	return errors.New("error")
}

func TestParseResponseErrorData(t *testing.T) {
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       &errorReadCloser{},
	}
	_, err := parseResponse(response)
	if err == nil {
		t.Fatal("no error")
	}
}

func createTestHTTPServer(t *testing.T) *net.TCPListener {
	addr, err := net.ResolveTCPAddr("tcp", "")
	if err != nil {
		t.Fatal(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	server := &http.Server{
		Handler: http.FileServer(http.Dir(testdata.Dir)),
	}
	go server.Serve(listener)
	return listener
}

func createTestURL(listener *net.TCPListener, filename string) *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   listener.Addr().String(),
		Path:   filename,
	}
}
