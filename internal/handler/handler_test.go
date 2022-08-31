package handler

import (
	"io/ioutil"
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
)

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
	Host = "localhost"
	Port = "8080"

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, Host+":"+Port, strings.NewReader(strings.Repeat("A", 2049)))

	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	CreateURL(c)

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
		request := httptest.NewRequest(http.MethodPost, Host+":"+Port, strings.NewReader(testCase.long))

		recorder := httptest.NewRecorder()
		c := e.NewContext(request, recorder)
		CreateURL(c)

		result := recorder.Result()
		defer result.Body.Close()

		if result.StatusCode != http.StatusCreated {
			t.Errorf("expected status %v; got %v", http.StatusCreated, result.StatusCode)
		}

		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}

		short := string(body)
		if short != "http://"+Host+":"+Port+"/"+testCase.short {
			t.Fatalf("expected answer to be %v; got %v", "http://"+Host+":"+Port+"/"+testCase.short, short)
		}
	}
}

func TestCreateURLInJSON(t *testing.T) {
	Host = "localhost"
	Port = "8080"

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, Host+":"+Port+"/api/shorten", strings.NewReader("{\"url\":\""+strings.Repeat("A", 2049)+"\"}"))

	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	CreateURL(c)

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
		request := httptest.NewRequest(http.MethodPost, Host+":"+Port+"/api/shorten", strings.NewReader("{\"url\":\""+testCase.long+"\"}"))

		recorder := httptest.NewRecorder()
		c := e.NewContext(request, recorder)
		CreateURL(c)

		result := recorder.Result()
		defer result.Body.Close()

		if result.StatusCode != http.StatusCreated {
			t.Errorf("expected status %v; got %v for link %v", http.StatusCreated, result.StatusCode, testCase.long)
		}

		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatalf("could not read response: %v", err)
		}

		short := string(body)
		if short != "http://"+Host+":"+Port+"/"+testCase.short {
			t.Fatalf("expected answer to be %v; got %v", "http://"+Host+":"+Port+"/"+testCase.short, short)
		}
	}
}

func TestRetrieveURL(t *testing.T) {
	Host = "localhost"
	Port = "8080"

	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "http://"+Host+":"+Port+"/", nil)

	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	RetrieveURL(c)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %v; got %v", http.StatusNotFound, result.StatusCode)
	}
}
