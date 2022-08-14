package middlewares

import (
	"errors"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var (
	LoginSessionKey    = "login_session"
	LoginSessionMaxAge = 12 * 60 * 60
)

type LoginSession struct {
	ID int
}

func SetLoginSession(c echo.Context, s LoginSession) error {
	sess, err := session.Get(LoginSessionKey, c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   LoginSessionMaxAge,
		HttpOnly: true,
	}
	sess.Values["ID"] = s.ID

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	return nil
}

func GetLoginSession(c echo.Context) (LoginSession, error) {
	sess, err := session.Get(LoginSessionKey, c)
	if err != nil {
		return LoginSession{}, err
	}
	id := sess.Values["ID"]
	if v, ok := id.(int); ok {
		return LoginSession{ID: v}, nil
	}

	return LoginSession{}, errors.New("unauthorized request")
}

func IsLoggedIn(c echo.Context) bool {
	_, err := GetLoginSession(c)
	if err == nil {
		return true
	} else {
		return false
	}
}

func LogOut(c echo.Context) error {
	sess, err := session.Get(LoginSessionKey, c)
	if err != nil {
		return err
	}
	sess.Values["ID"] = nil
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	return nil
}

type User struct {
	ID   int
	Pass string
	Name string
}

var AutoIncrement = 4

type InMemoryStore struct {
	store1 map[string]int
	store2 map[int]User
	mu     sync.RWMutex
}

func (h *InMemoryStore) Save(email, pass string, name *string) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	newID := AutoIncrement
	h.store1[email] = newID
	xname := ""
	if name == nil {
		xname = ""
	} else {
		xname = *name
	}
	h.store2[newID] = User{
		ID:   newID,
		Pass: pass,
		Name: xname,
	}
	AutoIncrement++
}

func (s *InMemoryStore) FindByEmail(email string) *User {
	if v1, ok := s.store1[email]; ok {
		if v2, ok := s.store2[v1]; ok {
			return &v2
		}
	}

	return nil
}

func (s *InMemoryStore) FindByID(id int) *User {
	if v, ok := s.store2[id]; ok {
		return &v
	}

	return nil
}

var DataStore = InMemoryStore{
	store1: map[string]int{
		"test1@example.com": 1,
		"test2@example.com": 2,
		"test3@example.com": 3,
	},
	store2: map[int]User{
		1: {ID: 1, Pass: "pass1", Name: "山田"},
		2: {ID: 2, Pass: "pass2", Name: "鈴木"},
		3: {ID: 3, Pass: "pass3", Name: "佐藤"},
	},
	mu: sync.RWMutex{},
}
