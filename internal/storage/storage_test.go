package storage

import (
	"testing"
)

var (
	tests = []struct {
		key   string
		value string
	}{
		{key: "key", value: "value"},
		{key: "a", value: "b"},
	}
)

func TestSetURL(t *testing.T) {
	for _, testCase := range tests {
		err := SetURL(testCase.key, testCase.value)
		if err != nil {
			t.Fatalf("can't save value %v for key %v", testCase.value, testCase.key)
		}
	}
}

func TestGetURL(t *testing.T) {
	_, ok := GetURL("non-existant-key")
	if ok {
		t.Fatal("got value for non existant key")
	}

	_, ok = GetURL("")
	if ok {
		t.Fatal("got value for empty key")
	}

	for _, testCase := range tests {
		val, ok := GetURL(testCase.key)
		if !ok {
			t.Fatalf("can't get value for key %v", testCase.key)
		}
		if val != testCase.value {
			t.Errorf("expected value %v for key %v; got %v", testCase.value, testCase.key, val)
		}
	}
}
