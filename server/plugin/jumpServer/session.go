package jumpServer

import (
	"errors"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/jumpServerMdl"
	"go.uber.org/zap"
)

type Store interface {
	Get(secret string) (*jumpServerMdl.SessionRecord, error)
	MarkUsed(secret string) error
	Save(sess *jumpServerMdl.SessionRecord) error
}

var ErrNotFound = errors.New("session not found")

type MemoryStore struct {
	mu   sync.Mutex
	data map[string]*jumpServerMdl.SessionRecord
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]*jumpServerMdl.SessionRecord),
	}
}

func (s *MemoryStore) Save(key string, sess *jumpServerMdl.SessionRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = sess
	return nil
}

func (s *MemoryStore) Get(secret string) (*jumpServerMdl.SessionRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, ok := s.data[secret]
	if !ok {
		return nil, ErrNotFound
	}
	return sess, nil
}

func (s *MemoryStore) MarkUsed(secret string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, ok := s.data[secret]
	if !ok {
		return ErrNotFound
	}
	sess.Used = true
	return nil
}

var (
	ErrInvalidSession = errors.New("invalid or expired session")
)

func (s *MemoryStore) Authenticate(username string) (*jumpServerMdl.SessionRecord, error) {
	// username == secret
	sess, err := s.Get(username)
	if err != nil {
		global.GVA_LOG.Error("session get failed", zap.String("username", username), zap.Error(err))
		return nil, ErrInvalidSession
	}

	//if sess.Used {
	//	global.GVA_LOG.Error("session used", zap.String("username", username), zap.Error(err))
	//	return nil, ErrInvalidSession
	//}

	if time.Now().After(sess.ExpiresAt) {
		global.GVA_LOG.Error("session expired", zap.String("username", username), zap.Error(err))
		return nil, ErrInvalidSession
	}

	// ⚠️ 这里立刻标记 Used，防止并发重放
	//if err := s.MarkUsed(username); err != nil {
	//	global.GVA_LOG.Error("session mark used", zap.String("username", username), zap.Error(err))
	//	return nil, ErrInvalidSession
	//}

	return sess, nil
}
