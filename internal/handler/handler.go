package handler

import (
	"crypto"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/combodga/Project/internal/storage"

	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	ServerAddr    string
	BaseURL       string
	Storage       *storage.Storage
	DBCredentials string
	Key           string
}

type Link struct {
	Result string `json:"result"`
}

type LinkJSON struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchLink struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type Element struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func New(serverAddr, baseURL, dbFile, dbCredentials string) (*Handler, error) {
	s, err := storage.New(dbFile, dbCredentials)
	if err != nil {
		err = fmt.Errorf("new storage: %w", err)
	}
	return &Handler{
		ServerAddr:    serverAddr,
		BaseURL:       baseURL,
		Storage:       s,
		DBCredentials: dbCredentials,
		Key:           "b8ffa0f4-3f11-44b1-b0bf-9109f47e468b",
	}, err
}

func (h *Handler) CreateURL(c echo.Context) error {
	user := getUser(c, h.Key)
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return fmt.Errorf("read request body: %w", err)
	}

	link := string(body)

	id, err := h.fetchID(c, user, link)
	if errors.Is(err, h.Storage.ErrDupKey) {
		return c.String(http.StatusConflict, h.BaseURL+"/"+id)
	}
	if err != nil {
		return fmt.Errorf("fetch id: %w", err)
	}

	return c.String(http.StatusCreated, h.BaseURL+"/"+id)
}

func (h *Handler) CreateURLInJSON(c echo.Context) error {
	user := getUser(c, h.Key)
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return fmt.Errorf("read request body: %w", err)
	}

	data := make(map[string]string)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	link, ok := data["url"]
	if !ok {
		return errors.New("error reading json")
	}

	uniqueErr := false
	id, err := h.fetchID(c, user, link)
	if errors.Is(err, h.Storage.ErrDupKey) {
		err = nil
		uniqueErr = true
	}
	if err != nil {
		return fmt.Errorf("fetch id: %w", err)
	}

	l := &Link{
		Result: h.BaseURL + "/" + id,
	}

	if uniqueErr {
		return c.JSON(http.StatusConflict, l)
	}

	return c.JSON(http.StatusCreated, l)
}

func (h *Handler) CreateBatchURL(c echo.Context) error {
	user := getUser(c, h.Key)
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return fmt.Errorf("read request body: %w", err)
	}

	var l []LinkJSON
	err = json.Unmarshal(body, &l)
	if err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	var bl []BatchLink
	var uniqueErr bool
	for result := range l {
		link := l[result]

		var id string
		id, err = h.fetchID(c, user, link.OriginalURL)
		if errors.Is(err, h.Storage.ErrDupKey) {
			err = nil
			uniqueErr = true
		}
		if err != nil {
			return fmt.Errorf("fetch id: %w", err)
		}

		bl = append(bl, BatchLink{
			CorrelationID: link.CorrelationID,
			ShortURL:      h.BaseURL + "/" + id,
		})
	}

	if uniqueErr {
		return c.JSON(http.StatusConflict, bl)
	}

	return c.JSON(http.StatusCreated, bl)
}

func (h *Handler) RetrieveURL(c echo.Context) error {
	user := getUser(c, h.Key)
	id := c.Param("id")

	url, status := h.Storage.GetURL(user, id)
	if status == 0 {
		return c.String(http.StatusNotFound, "error, there is no such link")
	} else if status == 2 {
		return c.String(http.StatusGone, "error, link was deleted")
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) ListURL(c echo.Context) error {
	user := getUser(c, h.Key)
	list, ok := h.Storage.ListURL(user)
	if !ok {
		return c.String(http.StatusNoContent, "error, you haven't any saved links")
	}

	var arr []*Element
	for shortURL, originalURL := range list {
		arr = append(arr, &Element{
			ShortURL:    h.BaseURL + "/" + shortURL,
			OriginalURL: originalURL,
		})
	}

	return c.JSON(http.StatusOK, arr)
}

func (h *Handler) DeleteURL(c echo.Context) error {
	user := getUser(c, h.Key)
	list, ok := h.Storage.ListURL(user)
	if !ok {
		return c.String(http.StatusBadRequest, "error, you haven't any saved links")
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return fmt.Errorf("read request body: %w", err)
	}

	var l []string
	err = json.Unmarshal(body, &l)
	if err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	for _, linkToDelete := range l {
		for id, savedLink := range list {
			if savedLink == linkToDelete {
				h.Storage.UpdateURL(user, id, true)
				break
			}
		}
	}

	return c.String(http.StatusAccepted, "URLs deleted")
}

func (h *Handler) Ping(c echo.Context) error {
	ok := h.Storage.Ping()
	if !ok {
		return c.String(http.StatusInternalServerError, "error, no connection to db")
	}

	return c.String(http.StatusOK, "db connected")
}

func (h *Handler) fetchID(c echo.Context, user, link string) (string, error) {
	if len(link) > 2048 {
		return "", c.String(http.StatusBadRequest, "error, the link cannot be longer than 2048 characters")
	}

	_, err := url.ParseRequestURI(link)
	if err != nil {
		return "", c.String(http.StatusBadRequest, "error, the link is invalid")
	}

	id, status := h.Storage.GetURL(user, link)
	if status == 0 {
		id, err = shortener(link)
		if err != nil {
			return "", c.String(http.StatusBadRequest, "error, failed to create a shortened URL")
		}
	}

	err = h.Storage.SetURL(user, id, link)
	if errors.Is(err, h.Storage.ErrDupKey) {
		return id, err
	}

	if err != nil {
		return "", c.String(http.StatusInternalServerError, "error, failed to store a shortened URL")
	}

	return id, err
}

func shortener(s string) (string, error) {
	h := crypto.MD5.New()
	if _, err := h.Write([]byte(s)); err != nil {
		return "", fmt.Errorf("abbreviation error URL: %w", err)
	}

	hash := string(h.Sum([]byte{}))
	hash = hash[len(hash)-5:]
	id := base58.Encode([]byte(hash))

	return id, nil
}

func getUser(c echo.Context, key string) string {
	user, err1 := readCookie(c, "user")
	sign, err2 := readCookie(c, "sign")
	if err1 == nil && err2 == nil && sign == getSign(user, key) {
		return user
	}

	user = randUser()
	writeCookie(c, "user", user)
	writeCookie(c, "sign", getSign(user, key))
	return user
}

func randUser() string {
	uuidWithHyphen := uuid.New()
	return uuidWithHyphen.String()
}

func getSign(user, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(user))
	dst := h.Sum(nil)
	return hex.EncodeToString(dst)[:32]
}

func writeCookie(c echo.Context, name, value string) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	c.SetCookie(cookie)
}

func readCookie(c echo.Context, name string) (string, error) {
	cookie, err := c.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
