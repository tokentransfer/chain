package store

import (
	"encoding/hex"
	"errors"
	"sync"

	"github.com/tokentransfer/interfaces/core"
)

type MemoryService struct {
	Name string

	db *sync.Map
}

func (service *MemoryService) Close() error {
	return nil
}

func (service *MemoryService) Init(c core.Config) error {
	service.db = &sync.Map{}
	return nil
}

func (service *MemoryService) Start() error {
	return nil
}

func (service *MemoryService) PutData(key []byte, value []byte) error {
	db := service.db

	s := hex.EncodeToString(key)
	db.Store(s, value)
	return nil
}

func (service *MemoryService) PutDatas(keys [][]byte, values [][]byte) error {
	db := service.db

	lk := len(keys)
	lv := len(values)
	if lk != lv {
		return errors.New("length error")
	}
	for i := 0; i < lk; i++ {
		s := hex.EncodeToString(keys[i])
		db.Store(s, values[i])
	}
	return nil
}

func (service *MemoryService) Flush() error {
	service.db = &sync.Map{}
	return nil
}

func (service *MemoryService) GetData(key []byte) ([]byte, error) {
	db := service.db

	s := hex.EncodeToString(key)
	value, ok := db.Load(s)
	if ok {
		return value.([]byte), nil
	}
	return nil, nil
}

func (service *MemoryService) GetDatas(keys [][]byte) ([][]byte, error) {
	db := service.db

	l := len(keys)
	bytes := make([][]byte, l)
	for i := 0; i < l; i++ {
		s := hex.EncodeToString(keys[i])
		value, ok := db.Load(s)
		if ok {
			bytes[i] = value.([]byte)
		} else {
			bytes[i] = nil
		}
	}
	return bytes, nil
}

func (service *MemoryService) HasData(key []byte) bool {
	db := service.db

	s := hex.EncodeToString(key)
	value, ok := db.Load(s)
	if !ok {
		return false
	}
	if value == nil {
		return false
	}
	data := value.([]byte)
	if len(data) == 0 {
		return false
	}

	return true
}

func (service *MemoryService) RemoveData(key []byte) error {
	db := service.db

	s := hex.EncodeToString(key)
	db.Delete(s)
	return nil
}

func (service *MemoryService) ListData(each func(key []byte, value []byte) error) error {
	db := service.db

	db.Range(func(k, v interface{}) bool {
		var key, value []byte

		if k != nil {
			key = k.([]byte)
		} else {
			key = nil
		}

		if v != nil {
			value = v.([]byte)
		} else {
			value = nil
		}

		err := each(key, value)
		if err != nil {
			return false
		}
		return true
	})
	return nil
}
