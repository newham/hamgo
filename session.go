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

// muti sessionStorage
var sessionStorages = make(map[string]sessionStorage)

// Register makes a session provider available by the provided name.
// If a Register is called twice with the same name or if the driver is nil,
// it panics.
func registsessionStorage(name string, ss sessionStorage) {
	if ss == nil {
		panic("session: Register provider is nil")
	}
	if _, dup := sessionStorages[name]; dup {
		panic("session: Register called twice for provider " + name)
	}
	sessionStorages[name] = ss
}

//Session : Session to store server's key-values (based on cookies)
type Session interface {
	Set(key, value interface{}) error //set session value
	Get(key interface{}) interface{}  //get session value
	Delete(key interface{}) error     //delete session value
	SessionID() string                //back current sessionID
	LeftTime() int64                  //get this session's timeout
	refresh()                         //update created time
}

type session struct {
	sid          string                      //session id
	timeAccessed time.Time                   //last access time
	value        map[interface{}]interface{} //session data value
	timeout      int64                       //timeout second
}

func (s *session) Set(key, value interface{}) error {
	s.value[key] = value
	return nil
}

func (s *session) Get(key interface{}) interface{} {
	if v, ok := s.value[key]; ok {
		return v
	}
	return nil
}

func (s *session) Delete(key interface{}) error {
	delete(s.value, key)
	return nil
}

func (s *session) SessionID() string {
	return s.sid
}

func (s *session) LeftTime() int64 {
	return s.timeAccessed.Unix() + s.timeout - time.Now().Unix()
}

func (s *session) refresh() {
	s.timeAccessed = time.Now()
}

type sessionStorage interface {
	InitSession(sid string, timeout int64) (Session, error)
	ReadSession(sid string) (Session, error)
	DestroySession(sid string) error
	RefreshSession(oldSid, newSid string) (Session, error)
}

type memorySessionStorage struct {
	lock     sync.Mutex         //used for lock
	sessions map[string]Session //used for inner store
	timeout  int64
}

func (msp *memorySessionStorage) InitSession(sid string, timeout int64) (Session, error) {
	msp.lock.Lock()
	defer msp.lock.Unlock()
	newsess := &session{sid: sid, timeAccessed: time.Now(), timeout: timeout, value: map[interface{}]interface{}{}}
	msp.sessions[sid] = newsess
	return newsess, nil
}

func (msp *memorySessionStorage) ReadSession(sid string) (Session, error) {
	if session, ok := msp.sessions[sid]; ok {
		if session.LeftTime() > 0 {
			return session, nil
		}
		msp.DestroySession(sid)
		return nil, errors.New("timeout")
	}
	return msp.InitSession(sid, msp.timeout)
}

func (msp *memorySessionStorage) DestroySession(sid string) error {
	if _, ok := msp.sessions[sid]; ok {
		delete(msp.sessions, sid)
		return nil
	}
	return errors.New("no sid")
}

func (pder *memorySessionStorage) RefreshSession(oldSid, newSid string) (Session, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if session, ok := pder.sessions[oldSid]; ok {
		delete(pder.sessions, oldSid)
		session.refresh()
		pder.sessions[newSid] = session
		return session, nil
	}
	return nil, errors.New("no sid")
}

type SessionManager interface {
	SessionStart(w http.ResponseWriter, r *http.Request) Session
	SessionDestroy(w http.ResponseWriter, r *http.Request)
	SessionRefresh(w http.ResponseWriter, r *http.Request) Session
}

type sessionManager struct {
	cookieName  string     //private cookiename
	lock        sync.Mutex // protects session
	provider    sessionStorage
	maxlifetime int64
}

func newSessionManager(storageName, cookieName string, maxlifetime int64) (SessionManager, error) {
	provider, ok := sessionStorages[storageName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", storageName)
	}
	return &sessionManager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}

func (sm *sessionManager) sessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (sm *sessionManager) SessionRefresh(w http.ResponseWriter, r *http.Request) (session Session) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	cookie, err := r.Cookie(sm.cookieName)
	sid := sm.sessionID()
	//read from store
	if err == nil && cookie.Value != "" {
		oldSid, _ := url.QueryUnescape(cookie.Value)
		session, _ = sm.provider.RefreshSession(oldSid, sid)
	} else {
		//init a new
		session, _ = sm.provider.InitSession(sid, sm.maxlifetime)
	}
	cookie = &http.Cookie{Name: sm.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(sm.maxlifetime)}
	http.SetCookie(w, cookie)

	return session
}

func (sm *sessionManager) SessionStart(w http.ResponseWriter, r *http.Request) Session {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	cookie, err := r.Cookie(sm.cookieName)
	//read from store
	if err == nil && cookie.Value != "" {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, err := sm.provider.ReadSession(sid)
		if err == nil {
			return session
		}
	}
	//init a new
	sid := sm.sessionID()
	session, _ := sm.provider.InitSession(sid, sm.maxlifetime)
	cookie = &http.Cookie{Name: sm.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(sm.maxlifetime)}
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
	manager.provider.DestroySession(cookie.Value)
	expiration := time.Now()
	cookie = &http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
	http.SetCookie(w, cookie)
}

var sessions SessionManager

func setSession(maxlifetime int64) {
	registsessionStorage("memory", &memorySessionStorage{sessions: map[string]Session{}})
	sessions, _ = newSessionManager("memory", "gosessionid", maxlifetime)
}
