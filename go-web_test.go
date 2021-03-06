package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestView(t *testing.T) {
	t.Fail()
}
func TestEdit(t *testing.T) {
	t.Fail()
}

func TestSave(t *testing.T) {
	t.Fail()
}

func TestRedirectFromRootHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rootHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusPermanentRedirect)
	}

	redirUrl, err := rr.Result().Location()
	if err != nil {
		t.Fatal(err)
	}

	if redirUrl.Path != rootRedirectURL {
		t.Errorf("handler returned unexpected redirect: got %v want %v", redirUrl.Path, rootRedirectURL)
	}
}
