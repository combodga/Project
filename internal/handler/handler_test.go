package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

var (
	tests = []struct {
		long  string
		short string
	}{
		{long: "http://yandex.ru", short: "K4fnNYm"},
		{long: "http://yandex.ru/", short: "PFPxFeq"},
		{long: "http://praktikum.yandex.ru", short: "T3Z75wy"},
		{long: "http://maps.yandex.ru", short: "HtkzjkX"},
		{long: "", short: "JFh7dgh"},
		{long: "//test_link", short: "7p3d7Tt"},
	}
	H *Handler
)

func TestInit(t *testing.T) {
	H = New("localhost:8080", "http://localhost:8080", "")
}

func TestShortener(t *testing.T) {
	for _, testCase := range tests {
		short, err := shortener(testCase.long)
		if err != nil {
			t.Errorf("can't shorten link %v", testCase.long)
		}

		if short != testCase.short {
			t.Fatalf("expected short link %v; got %v", testCase.short, short)
		}
	}
}

func TestCreateURL(t *testing.T) {
	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, H.ServerAddr, strings.NewReader(strings.Repeat("A", 2049)))

	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	H.CreateURL(c)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %v; got %v", http.StatusBadRequest, result.StatusCode)
	}

	for _, testCase := range tests {
		if testCase.long == "" {
			continue
		}

		e := echo.New()
		request := httptest.NewRequest(http.MethodPost, H.ServerAddr, strings.NewReader(testCase.long))

		recorder := httptest.NewRecorder()
		c := e.NewContext(request, recorder)
		H.CreateURL(c)

		result := recorder.Result()
		defer result.Body.Close()

		if result.StatusCode != http.StatusCreated {
			t.Errorf("expected status %v; got %v", http.StatusCreated, result.StatusCode)
		}

		body, err := io.ReadAll(result.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}

		short := string(body)
		expected := H.BaseURL + "/" + testCase.short
		if short != expected {
			t.Fatalf("expected answer to be %v; got %v", expected, short)
		}
	}
}

func TestCreateURLInJSON(t *testing.T) {
	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "http://" + H.ServerAddr + "/api/shorten", strings.NewReader("{\"url\":\""+strings.Repeat("A", 2049)+"\"}"))

	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	H.CreateURLInJSON(c)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %v; got %v", http.StatusBadRequest, result.StatusCode)
	}

	for _, testCase := range tests {
		if testCase.long == "" {
			continue
		}

		e := echo.New()
		request := httptest.NewRequest(http.MethodPost, H.ServerAddr + "/api/shorten", strings.NewReader("{\"url\":\"" + testCase.long + "\"}"))

		recorder := httptest.NewRecorder()
		c := e.NewContext(request, recorder)
		H.CreateURLInJSON(c)

		result := recorder.Result()
		defer result.Body.Close()

		if result.StatusCode != http.StatusCreated {
			t.Errorf("expected status %v; got %v for link %v", http.StatusCreated, result.StatusCode, testCase.long)
		}
	}
}

func TestRetrieveURL(t *testing.T) {
	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "http://" + H.ServerAddr + "/", nil)

	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	H.RetrieveURL(c)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %v; got %v", http.StatusNotFound, result.StatusCode)
	}
}
