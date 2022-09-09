package handler

import (
	"crypto"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/combodga/Project/internal/storage"

	"github.com/btcsuite/btcutil/base58"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	ServerAddr string
	BaseURL    string
	Storage    *storage.Storage
}

func New(serverAddr, baseURL, dbFile string) *Handler {
	return &Handler{
		ServerAddr: serverAddr,
		BaseURL:    "http://" + baseURL,
		Storage:    storage.New(dbFile),
	}
}

type Link struct {
	Result string `json:"result"`
}

func (h *Handler) CreateURL(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	link := string(body)

	id, err := h.fetchID(c, link)
	if err != nil {
		return err
	}

	return c.String(http.StatusCreated, h.BaseURL+"/"+id)
}

func (h *Handler) CreateURLInJSON(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	data := make(map[string]string)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	link, ok := data["url"]
	if !ok {
		return errors.New("error reading json")
	}

	id, err := h.fetchID(c, link)
	if err != nil {
		return err
	}

	l := &Link{
		Result: h.BaseURL + "/" + id,
	}
	return c.JSON(http.StatusCreated, l)
}

func (h *Handler) RetrieveURL(c echo.Context) error {
	id := c.Param("id")

	url, ok := h.Storage.GetURL(id)
	if !ok {
		return c.String(http.StatusNotFound, "error, there is no such link")
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) fetchID(c echo.Context, link string) (string, error) {
	if len(link) > 2048 {
		return "", c.String(http.StatusBadRequest, "error, the link cannot be longer than 2048 characters")
	}

	_, err := url.ParseRequestURI(link)
	if err != nil {
		return "", c.String(http.StatusBadRequest, "error, the link is invalid")
	}

	id, ok := h.Storage.GetURL(link)
	if !ok {
		id, err = shortener(link)
		if err != nil {
			return "", c.String(http.StatusBadRequest, "error, failed to create a shortened URL")
		}
	}

	err = h.Storage.SetURL(id, link)
	if err != nil {
		return "", c.String(http.StatusInternalServerError, "error, failed to store a shortened URL")
	}

	return id, nil
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
