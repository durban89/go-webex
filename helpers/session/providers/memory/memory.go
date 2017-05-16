package memory

import (
	"container/list"
	"sync"
	"time"

	"github.com/durban.zhang/webex/helpers/session"
)

type Provider struct {
	lock     sync.Mutex               // 锁操作
	sessions map[string]*list.Element // 内存存储
	list     *list.List               // gc 操作时使用
}

var provider = &Provider{list: list.New()}

func (provider *Provider) SessionInit(sid string) (session.Session, error) {
	provider.lock.Lock()
	defer provider.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newsess := &Session{sid: sid, lastAccessTime: time.Now(), value: v}
	element := provider.list.PushBack(newsess)
	provider.sessions[sid] = element
	return newsess, nil
}

func (provider *Provider) SessionRead(sid string) (session.Session, error) {
	if element, ok := provider.sessions[sid]; ok {
		return element.Value.(*Session), nil
	} else {
		sess, err := provider.SessionInit(sid)
		return sess, err
	}

	return nil, nil
}

func (provider *Provider) SessionDestory(sid string) error {
	if element, ok := provider.sessions[sid]; ok {
		delete(provider.sessions, sid)
		provider.list.Remove(element)
		return nil
	}

	return nil
}

func (provider *Provider) SessionGC(maxLiftTime int64) {
	provider.lock.Lock()
	defer provider.lock.Unlock()

	for {
		element := provider.list.Back()
		if element == nil {
			break
		}

		if (element.Value.(*Session).lastAccessTime.Unix() + maxLiftTime) < time.Now().Unix() {
			provider.list.Remove(element)
			delete(provider.sessions, element.Value.(*Session).sid)
		} else {
			break
		}
	}
}

func (provider *Provider) SessionUpdate(sid string) error {
	provider.lock.Lock()
	defer provider.lock.Lock()
	if element, ok := provider.sessions[sid]; ok {
		element.Value.(*Session).lastAccessTime = time.Now()
		provider.list.MoveToFront(element)
		return nil
	}

	return nil
}

type Session struct {
	sid            string                      // session id 唯一标识
	lastAccessTime time.Time                   // 最后访问时间
	value          map[interface{}]interface{} //session里面存储的值
}

func (s *Session) Set(key, value interface{}) error {
	s.value[key] = value
	provider.SessionUpdate(s.sid)
	return nil
}

func (s *Session) Get(key interface{}) interface{} {
	provider.SessionUpdate(s.sid)
	if v, ok := s.value[key]; ok {
		return v
	} else {
		return nil
	}

	return nil
}

func (s *Session) Delete(key interface{}) error {
	delete(s.value, key)
	provider.SessionUpdate(s.sid)
	return nil
}

func (s *Session) SessionID() string {
	return s.sid
}

func init() {
	provider.sessions = make(map[string]*list.Element, 0)
	session.Register("memory", provider)
}
