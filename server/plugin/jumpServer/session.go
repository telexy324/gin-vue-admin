package jumpServer

import (
	"errors"
	"sync"
	"time"
)

type SessionRecord struct {
	Secret     string
	UserID     int64
	TargetHost string
	TargetPort int

	ExpiresAt time.Time
	Used      bool
}

type Store interface {
	Get(secret string) (*SessionRecord, error)
	MarkUsed(secret string) error
	Save(sess *SessionRecord) error
}

var ErrNotFound = errors.New("session not found")

type MemoryStore struct {
	mu   sync.Mutex
	data map[string]*SessionRecord
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]*SessionRecord),
	}
}

func (s *MemoryStore) Save(sess *SessionRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[sess.Secret] = sess
	return nil
}

func (s *MemoryStore) Get(secret string) (*SessionRecord, error) {
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
