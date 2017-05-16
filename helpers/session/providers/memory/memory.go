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

func (provider *Provider) SessionInit(sid string) (session.Session, error) {}
func (provider *Provider) SessionRead(sid string) (session.Session, error) {}
func (provider *Provider) SessionDestory(sid string) error                 {}
func (provider *Provider) SessionGC(maxLiftTime int64)                     {}
func (provider *Provider) SessionUpdate(sid string) error                  {}

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
	rerurn s.sid
}

func init() {
	provider.session = make(map[string]*list.Element, 0)
	session.Register("memory", provider)
}
