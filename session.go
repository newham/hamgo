package hamgo

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

//********** Session **********

const (
	sessionCookieName     = "gosessionid"
	second                = 1
	minute                = 60 * second
	hour                  = 60 * minute
	defaultSessionMaxTime = hour
	confSessionMaxTime    = "session_max_time"
)

type Session interface {
	Set(key, value interface{}) error //set session value
	Get(key interface{}) interface{}  //get session value
	Delete(key interface{}) error     //delete session value
	ID() string                       //back current sessionID
	LeftTime() int64                  //get this session's timeout
	newID()
	refreshTime()
}

type session struct {
	sid          string                      //session id
	timeAccessed time.Time                   //last access time
	value        map[interface{}]interface{} //session data value
	timeout      int                         //timeout second
}

func newSession(max int) Session {
	return &session{sid: uuid(32), value: map[interface{}]interface{}{}, timeAccessed: time.Now(), timeout: max}
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

func (s *session) ID() string {
	return s.sid
}

func (s *session) LeftTime() int64 {
	return s.timeAccessed.Unix() + int64(s.timeout) - time.Now().Unix()
}

func (s *session) newID() {
	s.sid = uuid(32)
}

func (s *session) refreshTime() {
	s.timeAccessed = time.Now()
}

//********** SessionStorage **********

type sessionStorage interface {
	Put(Session)
	Get(string) Session
}

type memorySessionStorage struct {
	lock     sync.Mutex         //used for lock
	sessions map[string]Session //used for inner store
}

func (msp *memorySessionStorage) Put(session Session) {
	msp.lock.Lock()
	msp.sessions[session.ID()] = session
	msp.lock.Unlock()
}

func (msp *memorySessionStorage) Get(ID string) Session {
	return msp.sessions[ID]
}

//********** SessionManager **********

type sessionManager struct {
	storage     sessionStorage
	maxlifetime int
}

func (sm *sessionManager) GetSession(r *http.Request, w http.ResponseWriter) Session {
	session := sm.storage.Get(getSessionID(r))
	if session == nil {
		session = newSession(sm.maxlifetime)
		sm.SetSession(w, session)
	} else {
		if session.LeftTime() < 1 {
			sm.DelSession(r, w)
			return nil
		}
		//refresh left time
		session.refreshTime()
	}
	return session
}

func (sm *sessionManager) SetSession(w http.ResponseWriter, session Session) {
	sm.storage.Put(session)
	cookie := &http.Cookie{Name: sessionCookieName, Value: session.ID(), Path: "/", HttpOnly: true}
	http.SetCookie(w, cookie)
}

func (sm *sessionManager) DelSession(r *http.Request, w http.ResponseWriter) {
	cookie := &http.Cookie{Name: sessionCookieName, Value: getSessionID(r), Path: "/", HttpOnly: true, Expires: time.Now(), MaxAge: -1}
	http.SetCookie(w, cookie)
}

func (sm *sessionManager) RefreshSession(r *http.Request, w http.ResponseWriter) {
	session := sm.GetSession(r, w)
	if session == nil {
		return
	}
	session.newID()
	sm.SetSession(w, session)
}

//********** Tool **********

func getSessionID(r *http.Request) string {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return ""
	}
	id, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return ""
	}
	return id
}
