package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestEmptyAndDuplicateSelections(t *testing.T) {
	d, e := NewData(map[string]interface{}{"results": []interface{}{map[string]interface{}{"a": "yes", "b": "yes"}}})
	if e != nil {
		t.Fatal(e)
	}
	if !d.FilterValue("yes") || len(d.currentResultSeq) != 1 {
		t.Fatal("duplicate row")
	}
	if d.FilterKeyValue("a", "absent") || len(d.currentResultSeq) != 0 {
		t.Fatal("stale results retained")
	}
}
func TestHTTPUploadFilterDownload(t *testing.T) {
	r := newRouter()
	for _, tt := range []struct {
		input string
		want  int
	}{{`{"results":[{"name":"alpha"}]}`, 200}, {`{"other":[]}`, 400}, {`{"results":[3]}`, 400}} {
		var b bytes.Buffer
		m := multipart.NewWriter(&b)
		f, _ := m.CreateFormFile("file", "sample.json")
		f.Write([]byte(tt.input))
		m.Close()
		req := httptest.NewRequest("POST", "/upload", &b)
		req.Header.Set("Content-Type", m.FormDataContentType())
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tt.want {
			t.Fatalf("upload %s: %d", tt.input, w.Code)
		}
	}
	form := url.Values{"key": {"name"}, "value": {"missing"}}
	req := httptest.NewRequest("POST", "/filter_key_value", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/download", nil))
	var rows []interface{}
	if e := json.Unmarshal(w.Body.Bytes(), &rows); e != nil || len(rows) != 0 {
		t.Fatalf("download: %s", w.Body.String())
	}
}
