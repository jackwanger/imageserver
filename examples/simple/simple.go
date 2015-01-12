package main

import (
	"net/http"

	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_nfntresize "github.com/pierrre/imageserver/http/nfntresize"
	imageserver_native "github.com/pierrre/imageserver/native"
	_ "github.com/pierrre/imageserver/native/bmp"
	_ "github.com/pierrre/imageserver/native/gif"
	_ "github.com/pierrre/imageserver/native/jpeg"
	imageserver_native_nfntresize "github.com/pierrre/imageserver/native/nfntresize"
	_ "github.com/pierrre/imageserver/native/png"
	_ "github.com/pierrre/imageserver/native/tiff"
	imageserver_testdata "github.com/pierrre/imageserver/testdata"
)

func main() {
	server := imageserver_testdata.Server
	server = &imageserver_native.Server{
		Server:    server,
		Processor: &imageserver_native_nfntresize.Processor{},
	}
	handler := &imageserver_http.Handler{
		Parser: &imageserver_http.ListParser{
			&imageserver_http.SourceParser{},
			&imageserver_http_nfntresize.Parser{},
			&imageserver_http.FormatParser{},
			&imageserver_http.QualityParser{},
		},
		Server: server,
		ErrorFunc: func(err error, request *http.Request) {
			println(err.Error())
		},
	}
	http.Handle("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
