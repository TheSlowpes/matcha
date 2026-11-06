package view

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/floatpane/matcha/clib"
)

func TestPrefetchRemoteImagesRunsInParallel(t *testing.T) {
	var pngData bytes.Buffer
	if err := png.Encode(&pngData, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	const delay = 200 * time.Millisecond
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(pngData.Bytes())
	}))
	defer srv.Close()

	const n = 6
	elements := make([]clib.HTMLElement, 0, n+1)
	for i := range n {
		elements = append(elements, clib.HTMLElement{Type: clib.HElemImage, Attr1: fmt.Sprintf("%s/img%d.png", srv.URL, i)})
	}
	elements = append(elements, clib.HTMLElement{Type: clib.HElemImage, Attr1: srv.URL + "/img0.png"})

	start := time.Now()
	got := prefetchRemoteImages(elements)
	elapsed := time.Since(start)

	if len(got) != n {
		t.Fatalf("got %d images, want %d", len(got), n)
	}
	for url, payload := range got {
		if payload == "" {
			t.Errorf("empty payload for %s", url)
		}
	}
	if elapsed >= n*delay/2 {
		t.Errorf("prefetch took %v; downloads do not appear to run in parallel", elapsed)
	}
}
