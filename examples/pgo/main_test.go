package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.com/golang-commonmark/markdown"
)

func BenchmarkReadBody(b *testing.B) {
	req, _ := http.NewRequest("POST", "/render", strings.NewReader("Hello, world!"))
	for i := 0; i < b.N; i++ {
		_, err := io.ReadAll(req.Body)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRenderMarkdown(b *testing.B) {
	md := []byte("# Hello, world!")
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		renderer := markdown.New(markdown.XHTMLOutput(true), markdown.Typographer(true), markdown.Linkify(true), markdown.Tables(true))
		err := renderer.Render(&buf, md)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteResponse(b *testing.B) {
	resp := httptest.NewRecorder()
	buf := bytes.NewBufferString("Hello, world!")
	for i := 0; i < b.N; i++ {
		_, err := io.Copy(resp.Body, buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRender(b *testing.B) {
	req, _ := http.NewRequest("POST", "/render", strings.NewReader("# Hello, world!"))
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(render)
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(resp, req)
		if resp.Code != http.StatusOK {
			b.Fatalf("expected status %d; got %d", http.StatusOK, resp.Code)
		}
	}
}
