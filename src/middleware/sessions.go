package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"hiking_trails/src/models"
	"log"
	"net/http"
	"sync"
)

type Session struct {
	Id     string
	values map[string]interface{}
	store  *SessionStore
}

func NewSession(store *SessionStore) Session {
	id := mustGenerateSessionId()
	return Session{id, make(map[string]interface{}), store}
}

func (session *Session) Save() error {
	return session.store.Save(*session)
}

func (session *Session) Create() {
	session.store.Create(*session)
}

func (session *Session) Delete() {
	session.store.Delete(session.Id)
}

func (session *Session) Set(key string, value interface{}) {
	session.values[key] = value
}

func (session *Session) Get(key string) interface{} {
	value, exist := session.values[key]
	if !exist {
		return nil
	}

	return value
}

func mustGenerateSessionId() string {
	id := make([]byte, 32)

	_, err := rand.Read(id)
	if err != nil {
		panic(err)
	}

	idAsHex := make([]byte, hex.EncodedLen(len(id)))
	hex.Encode(idAsHex, id)

	return string(idAsHex)
}

// FIXME Implement session expiration.
type SessionStore struct {
	lock     *sync.Mutex
	sessions map[string]Session
}

func NewSessionStore() *SessionStore {
	return &SessionStore{&sync.Mutex{}, make(map[string]Session, 0)}
}

func (store *SessionStore) Get(id string) Session {
	store.lock.Lock()
	defer store.lock.Unlock()

	// Returns empty session if not in map.
	session, _ := store.sessions[id]

	// Do not return as pointer since we want a copy. Else it would not be thread safe.
	return session
}

func (store *SessionStore) Delete(id string) {
	store.lock.Lock()
	defer store.lock.Unlock()

	delete(store.sessions, id)
}

func (store *SessionStore) Create(session Session) {
	store.lock.Lock()
	defer store.lock.Unlock()

	store.sessions[session.Id] = session
}

func (store *SessionStore) Save(session Session) error {
	store.lock.Lock()
	defer store.lock.Unlock()

	_, exist := store.sessions[session.Id]
	if !exist {
		return fmt.Errorf("Session %d does not exist.", session.Id)
	}

	store.sessions[session.Id] = session
	return nil
}

func Sessions(name string, store *SessionStore) martini.Handler {
	return func(res http.ResponseWriter, request *http.Request, c martini.Context, logger *log.Logger) {
		session := NewSession(store)

		cookie, err := request.Cookie("SessionId")
		if err == nil {
			tmpSession := store.Get(cookie.Value)
			if tmpSession.Id != "" {
				session = tmpSession
			}
		}

		c.Map(&session)
	}
}

func AdministratorRequired(session *Session, render render.Render, logger *log.Logger) {
	isAdministrator := session.Get("isAdministrator")

	if isAdministrator == nil || !isAdministrator.(bool) {
		err := models.NewAPIError(401, "Unauthorized", nil)
		renderErrorAsJson(err, render, logger)
	}
}

func renderErrorAsJson(err error, render render.Render, logger *log.Logger) {
	apiError, isApiError := err.(*models.APIError)

	if !isApiError {
		apiError = models.NewAPIError(500, "Internal Server Error", err)
	}

	apiError.RenderAsJson(render, logger)
}
