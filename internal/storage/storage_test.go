package storage

import (
	"testing"
)

var (
	tests = []struct {
		user  string
		key   string
		value string
	}{
		{user: "test", key: "key", value: "value"},
		{user: "test", key: "a", value: "b"},
	}
	S *Storage
)

func TestInit(t *testing.T) {
	var err error
	S, err = New("")
	if err != nil {
		t.Fatal("can't start test")
	}
}

func TestSetURL(t *testing.T) {
	for _, testCase := range tests {
		err := S.SetURL(testCase.user, testCase.key, testCase.value)
		if err != nil {
			t.Fatalf("can't save value %v for key %v", testCase.value, testCase.key)
		}
	}
}

func TestGetURL(t *testing.T) {
	_, ok := S.GetURL("non-existant-key")
	if ok {
		t.Fatal("got value for non existant key")
	}

	_, ok = S.GetURL("")
	if ok {
		t.Fatal("got value for empty key")
	}

	for _, testCase := range tests {
		val, ok := S.GetURL(testCase.key)
		if !ok {
			t.Fatalf("can't get value for key %v", testCase.key)
		}
		if val != testCase.value {
			t.Errorf("expected value %v for key %v; got %v", testCase.value, testCase.key, val)
		}
	}
}
