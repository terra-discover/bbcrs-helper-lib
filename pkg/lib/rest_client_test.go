package lib

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2/utils"
)

func TestBasicRestClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	dummyRest := RestClient{
		URL:     server.URL,
		Request: `{"type" : "OK"}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Method:  "GET",
		Timeout: 5,
	}

	_, httpCode := dummyRest.Execute()
	utils.AssertEqual(t, true, httpCode >= 200, "Call Rest API with complete param")
}

func TestDefaultParamRestClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	dummyRest := RestClient{
		URL: server.URL,
	}

	_, httpCode := dummyRest.Execute()
	utils.AssertEqual(t, true, httpCode >= 200, "Call Rest API with fallback method & timeout")
}

func TestInvalidURLRestClient(t *testing.T) {
	dummyRest := RestClient{
		URL: "lorem-ipsum",
	}

	_, httpCode := dummyRest.Execute()
	utils.AssertEqual(t, true, httpCode >= 200 || httpCode == 0, "Call Rest API with fallback method & timeout")
}

func TestMethodsRestClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	dummyRest := RestClient{}
	_, httpCode := dummyRest.SetURL(server.URL).
		SetMethod("UNDEFINED").
		SetTimeout(1).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).
		SetRequest(`{"type":"OK"}`).
		Execute()

	utils.AssertEqual(t, true, httpCode >= 200 || httpCode == 0, "Call Rest API with defined method")
}

func TestFindSlice(t *testing.T) {
	index, hasSlice := FindSlice([]string{"A", "B", "C"}, "C")
	utils.AssertEqual(t, 2, index, "FindSlice index result")
	utils.AssertEqual(t, true, hasSlice, "FindSlice hasSlice result")

	emptyIndex, emptyHasSlice := FindSlice([]string{"A", "B", "C"}, "Z")
	utils.AssertEqual(t, -1, emptyIndex, "FindSlice index empty result")
	utils.AssertEqual(t, false, emptyHasSlice, "FindSlice hasSlice empty result")
}
