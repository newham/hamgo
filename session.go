package hamgo

import (
	"container/list"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	second                = 1
	minute                = 60
	hour                  = 60 * 60
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

func (manager *sessionManager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.sessionID()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}

func (manager *sessionManager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxlifetime)
	time.AfterFunc(time.Duration(manager.maxlifetime), func() { manager.GC() })
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
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
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
}

var pder = &provider{list: list.New()}

type sessionStore struct {
	sid          string                      //session id唯一标示
	timeAccessed time.Time                   //最后访问时间
	value        map[interface{}]interface{} //session里面存储的值
}

func (st *sessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	pder.SessionUpdate(st.sid)
	return nil
}

func (st *sessionStore) Get(key interface{}) interface{} {
	pder.SessionUpdate(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	}
	return nil
}

func (st *sessionStore) Delete(key interface{}) error {
	delete(st.value, key)
	pder.SessionUpdate(st.sid)
	return nil
}

func (st *sessionStore) SessionID() string {
	return st.sid
}

type provider struct {
	lock     sync.Mutex               //用来锁
	sessions map[string]*list.Element //用来存储在内存
	list     *list.List               //用来做gc
}

func (pder *provider) SessionInit(sid string) (Session, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newsess := &sessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	element := pder.list.PushBack(newsess)
	pder.sessions[sid] = element
	return newsess, nil
}

func (pder *provider) SessionRead(sid string) (Session, error) {
	if element, ok := pder.sessions[sid]; ok {
		return element.Value.(*sessionStore), nil
	}
	sess, err := pder.SessionInit(sid)
	return sess, err
}

func (pder *provider) SessionDestroy(sid string) error {
	if element, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		pder.list.Remove(element)
		return nil
	}
	return nil
}

func (pder *provider) SessionGC(maxlifetime int64) {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	for {
		element := pder.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*sessionStore).timeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*sessionStore).sid)
		} else {
			break
		}
	}
}

func (pder *provider) SessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*sessionStore).timeAccessed = time.Now()
		pder.list.MoveToFront(element)
		return nil
	}
	return nil
}

var sessions *sessionManager

func newSession(maxlifetime int64) {
	pder.sessions = make(map[string]*list.Element, 0)
	registSessionProvider("memory", pder)
	sessions, _ = newSessionManager("memory", "gosessionid", maxlifetime)
	go sessions.GC()
}
