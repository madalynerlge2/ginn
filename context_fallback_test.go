package gin

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShouldBindJSON_FallbackRead(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(`{"foo": `))
	c, _ := CreateTestContext(w)
	c.Request = req

	var myStruct struct {
		Foo string `json:"foo"`
	}
	err := c.ShouldBindJSON(&myStruct)
	if err == nil {
		t.Fatal("expected error binding invalid JSON, got nil")
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		t.Fatalf("failed to read request body: %v", err)
	}

	expected := `{"foo": `
	if string(bodyBytes) != expected {
		t.Errorf("expected body %q, got %q", expected, string(bodyBytes))
	}
}
