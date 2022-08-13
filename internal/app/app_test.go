package app

import (
	// "bytes"
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
		{long: "http://yandex.ru", short: "obrnfe4"},
		{long: "http://yandex.ru/", short: "xtklsi"},
		{long: "http://praktikum.yandex.ru", short: "5te+dbq"},
		{long: "http://maps.yandex.ru", short: "lbohctw"},
		{long: "", short: "moz4qn4"},
		{long: "//test_link", short: "pf7smo8"},
	}
)

func TestGetURL(t *testing.T) {
	for _, testCase := range tests {
		short, err := shortener(testCase.long)
		if err != nil {
			t.Errorf("can't shorten link %v", testCase.long)
		}

		if short != "http://localhost:8080/"+testCase.short {
			t.Fatalf("expected short link %v; got %v", "http://localhost:8080/"+testCase.short, short)
		}

		long, _ := getURL(testCase.short)
		if long != testCase.long {
			t.Fatalf("expected long link %v; got %v", testCase.long, long)
		}
	}
}

func TestShortener(t *testing.T) {
	for _, testCase := range tests {
		short, err := shortener(testCase.long)
		if err != nil {
			t.Errorf("can't shorten link %v", testCase.long)
		}

		if short != "http://localhost:8080/"+testCase.short {
			t.Fatalf("expected short link %v; got %v", "http://localhost:8080/"+testCase.short, short)
		}
	}
}

func TestNewLink(t *testing.T) {
	{
		e := echo.New()
		request := httptest.NewRequest(http.MethodPost, "localhost:8080", strings.NewReader(strings.Repeat("A", 2049)))

		recorder := httptest.NewRecorder()
		c := e.NewContext(request, recorder)
		newLink(c)

		result := recorder.Result()
		defer result.Body.Close()

		if result.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status %v; got %v", http.StatusBadRequest, result.StatusCode)
		}
	}

	// {
	//     e := echo.New()
	//     request := httptest.NewRequest(http.MethodPost, "localhost:8080", strings.NewReader("123456"))

	//     recorder := httptest.NewRecorder()
	//     c := e.NewContext(request, recorder)
	//     newLink(c)

	//     result := recorder.Result()
	//     defer result.Body.Close()

	//     if result.StatusCode != http.StatusBadRequest {
	//         t.Errorf("expected status %v; got %v", http.StatusBadRequest, result.StatusCode)
	//     }
	// }

	for _, testCase := range tests {
		if testCase.long == "" {
			continue
		}

		e := echo.New()
		request := httptest.NewRequest(http.MethodPost, "localhost:8080", strings.NewReader(testCase.long))

		recorder := httptest.NewRecorder()
		c := e.NewContext(request, recorder)
		newLink(c)

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
		if short != "http://localhost:8080/"+testCase.short {
			t.Fatalf("expected answer to be %v; got %v", "http://localhost:8080/"+testCase.short, short)
		}
	}
}

func TestGetLink(t *testing.T) {
	{
		e := echo.New()
		request := httptest.NewRequest(http.MethodGet, "localhost:8080", nil)

		recorder := httptest.NewRecorder()
		c := e.NewContext(request, recorder)
		getLink(c)

		result := recorder.Result()
		defer result.Body.Close()

		if result.StatusCode != http.StatusNotFound {
			t.Errorf("expected status %v; got %v", http.StatusNotFound, result.StatusCode)
		}
	}
	// for _, testCase := range tests {
	//     e := echo.New()
	//     request := httptest.NewRequest(http.MethodGet, "localhost:8080/" + testCase.short, nil)

	//     recorder := httptest.NewRecorder()
	//     c := e.NewContext(request, recorder)
	//     getLink(c)

	//     result := recorder.Result()
	//     defer result.Body.Close()

	//     if result.StatusCode != http.StatusTemporaryRedirect {
	//         t.Errorf("expected status %v; got %v", http.StatusTemporaryRedirect, result.StatusCode)
	//     }

	//     body, err := ioutil.ReadAll(result.Body)
	//     if err != nil {
	//         t.Fatalf("could not read response: %v", err)
	//     }

	//     long := string(body)
	//     if long != testCase.long {
	//         t.Fatalf("expected answer to be %v; got %v", testCase.long, long)
	//     }
	// }
}
