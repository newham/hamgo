package hamgo

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	second                = 1
	minute                = 60 * second
	hour                  = 60 * minute
	defaultSessionMaxTime = hour
	confSessionMaxTime    = "session_max_time"
)

type sessionManager struct {
	cookieName  string     //private cookiename
	lock        sync.Mutex // protects session
	provider    sessionProvider
	maxlifetime int64
}

func newSessionManager(provideName, cookieName string, maxlifetime int64) (*sessionManager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &sessionManager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}
func (manager *sessionManager) sessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *sessionManager) SessionStart(w http.ResponseWriter, r *http.Request) Session {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	//read from store
	if err == nil && cookie.Value != "" {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, err := manager.provider.SessionRead(sid, manager.maxlifetime)
		if err == nil {
			return session
		}
	}
	//init a new
	sid := manager.sessionID()
	session, _ := manager.provider.SessionInit(sid)
	cookie = &http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
	http.SetCookie(w, cookie)

	return session
}

//Destroy sessionid
func (manager *sessionManager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionDestroy(cookie.Value)
	expiration := time.Now()
	cookie = &http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
	http.SetCookie(w, cookie)
}

type sessionProvider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string, maxTime int64) (Session, error)
	SessionDestroy(sid string) error
}

var provides = make(map[string]sessionProvider)

// Register makes a session provider available by the provided name.
// If a Register is called twice with the same name or if the driver is nil,
// it panics.
func registSessionProvider(name string, provider sessionProvider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provider " + name)
	}
	provides[name] = provider
}

//Session : Session to store server's key-values (based on cookies)
type Session interface {
	Set(key, value interface{}) error //set session value
	Get(key interface{}) interface{}  //get session value
	Delete(key interface{}) error     //delete session value
	SessionID() string                //back current sessionID
	leftTime(timeout int64) int64     //get this session's timeout
	refresh()                         //update created time
}

type sessionStore struct {
	sid          string                      //session id
	timeAccessed time.Time                   //last access time
	value        map[interface{}]interface{} //session data value
}

func (st *sessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	return nil
}

func (st *sessionStore) Get(key interface{}) interface{} {
	if v, ok := st.value[key]; ok {
		return v
	}
	return nil
}

func (st *sessionStore) Delete(key interface{}) error {
	delete(st.value, key)
	return nil
}

func (st *sessionStore) SessionID() string {
	return st.sid
}

func (st *sessionStore) leftTime(timeout int64) int64 {
	return st.timeAccessed.Unix() + timeout - time.Now().Unix()
}

func (st *sessionStore) refresh() {
	st.timeAccessed = time.Now()
}

type provider struct {
	lock     sync.Mutex         //used for lock
	sessions map[string]Session //used for inner store
}

func (pder *provider) SessionInit(sid string) (Session, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	newsess := &sessionStore{sid: sid, timeAccessed: time.Now(), value: map[interface{}]interface{}{}}
	pder.sessions[sid] = newsess
	return newsess, nil
}

func (pder *provider) SessionRead(sid string, maxTime int64) (Session, error) {
	if session, ok := pder.sessions[sid]; ok {
		println(session.leftTime(maxTime))
		if session.leftTime(maxTime) > 0 {
			session.refresh()
			return session, nil
		}
		pder.SessionDestroy(sid)
		return nil, errors.New("timeout")
	}
	return pder.SessionInit(sid)
}

func (pder *provider) SessionDestroy(sid string) error {
	if _, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		return nil
	}
	return errors.New("no sid")
}

func (pder *provider) SessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if session, ok := pder.sessions[sid]; ok {
		session.refresh()
		return nil
	}
	return errors.New("no sid")
}

var sessions *sessionManager

func setSession(maxlifetime int64) {
	registSessionProvider("memory", &provider{sessions: map[string]Session{}})
	sessions, _ = newSessionManager("memory", "gosessionid", maxlifetime)
}
