package default_session

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	. "github.com/infrago/base"
	"github.com/infrago/session"
)

var (
	errInvalidCacheConnection = errors.New("Invalid session connection.")
	errInvalidCacheData       = errors.New("Invalid session data.")
)

type (
	defaultDriver  struct{}
	defaultConnect struct {
		mutex sync.RWMutex

		instance *session.Instance
		setting  defaultSetting
		sessions sync.Map
	}
	defaultSetting struct {
	}
	defaultValue struct {
		Value  []byte
		Expire time.Time
	}
)

// 连接
func (driver *defaultDriver) Connect(inst *session.Instance) (session.Connect, error) {
	setting := defaultSetting{}

	return &defaultConnect{
		instance: inst, setting: setting,
		sessions: sync.Map{},
	}, nil
}

// 打开连接
func (this *defaultConnect) Open() error {
	return nil
}

// 关闭连接
func (this *defaultConnect) Close() error {
	return nil
}

// 查询会话，
func (this *defaultConnect) Read(id string) ([]byte, error) {
	if value, ok := this.sessions.Load(id); ok {
		if vv, ok := value.(defaultValue); ok {
			if vv.Expire.Unix() > time.Now().Unix() {
				return vv.Value, nil
			} else {
				//过期了就删除
				this.Delete(id)
			}
		}
	}
	return nil, errInvalidCacheData
}

// 更新会话
func (this *defaultConnect) Write(id string, data []byte, expire time.Duration) error {
	now := time.Now()

	value := defaultValue{
		Value: data, Expire: now.Add(expire),
	}

	this.sessions.Store(id, value)

	return nil
}

// 查询会话，
func (this *defaultConnect) Exists(id string) (bool, error) {
	if _, ok := this.sessions.Load(id); ok {
		return ok, nil
	}
	return false, errors.New("会话读取失败")
}

// 删除会话
func (this *defaultConnect) Delete(id string) error {
	this.sessions.Delete(id)
	return nil
}

func (this *defaultConnect) Keys(prefix string) ([]string, error) {
	ids := []string{}

	this.sessions.Range(func(k, _ Any) bool {
		id := fmt.Sprintf("%v", k)

		if strings.HasPrefix(id, prefix) {
			ids = append(ids, id)
		}
		return true
	})
	return ids, nil
}
func (this *defaultConnect) Clear(prefix string) error {
	if ids, err := this.Keys(prefix); err == nil {
		for _, id := range ids {
			this.sessions.Delete(id)
		}
		return nil
	} else {
		return err
	}
}
