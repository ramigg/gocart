package handlers

import (
	"net/http"
	"testing"

	"net/http/httptest"

	"github.com/gorilla/mux"
)

func TestQparam(t *testing.T) {
	req, err := http.NewRequest("GET", "/?first=a&2nd=b&5nd=e&last=z", nil)
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		param    string
		expected string
	}{
		{"first", "a"},
		{"2nd", "b"},
		{"5nd", "e"},
		{"last", "z"},
		{"notexists", ""},
	}

	for _, test := range tests {
		qp := qParam(test.param, req)
		if qp != test.expected {
			t.Errorf("Expected qParam(%q) to be %v got %v", test.param, test.expected, qp)
		}
	}

}

func TestMuxVarMustInt(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mr := mux.NewRouter()
	mr.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := muxVarMustInt("id", r)
		if id != 1 {
			t.Errorf("Expected %d got %v", 1, id)
		}
	}).Methods("GET")

	mr.ServeHTTP(rr, req)
}
