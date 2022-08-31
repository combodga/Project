package handler

import (
	"crypto"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/combodga/Project/internal/storage"

	"github.com/btcsuite/btcutil/base58"
	"github.com/labstack/echo/v4"
)

var (
	Host = "localhost"
	Port = "8080"
)

// func SetHost(s string) {
// 	host = s
// }

// func SetPort(s string) {
// 	port = s
// }

func CreateURL(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	link := string(body)

	if len(link) > 2048 {
		return c.String(http.StatusBadRequest, "error, the link cannot be longer than 2048 characters")
	}

	_, err = url.ParseRequestURI(link)
	if err != nil {
		return c.String(http.StatusBadRequest, "error, the link is invalid")
	}

	id, ok := storage.GetURL(link)
	if !ok {
		id, err = shortener(link)
		if err != nil {
			return c.String(http.StatusBadRequest, "error, failed to create a shortened URL")
		}
	}

	err = storage.SetURL(id, link)
	if err != nil {
		return c.String(http.StatusInternalServerError, "error, failed to store a shortened URL")
	}

	return c.String(http.StatusCreated, "http://"+Host+":"+Port+"/"+id)
}

func RetrieveURL(c echo.Context) error {
	id := c.Param("id")

	url, ok := storage.GetURL(id)
	if !ok {
		return c.String(http.StatusNotFound, "error, there is no such link")
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func shortener(s string) (string, error) {
	h := crypto.MD5.New()
	if _, err := h.Write([]byte(s)); err != nil {
		return "", fmt.Errorf("abbreviation error URL: %v", err)
	}

	hash := string(h.Sum([]byte{}))
	hash = hash[len(hash)-5:]
	id := base58.Encode([]byte(hash))

	return id, nil
}
