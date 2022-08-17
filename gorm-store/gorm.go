package gorm_store

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"time"
)


type Session struct {
	Key string `gorm:"primaryKey"`
	Data string
	ExpireTime time.Time `gorm:"index"`
}


func New(db *gorm.DB) *GormStore {
	store := &GormStore{
		db: db,
		model: &Session{},
	}
	err := store.migrate()
	if err != nil {
		panic(err)
	}
	return store
}


type GormStore struct {
	db *gorm.DB
	model *Session
}


func (self *GormStore) Get(key string) (map[string]interface{}, error) {
	model := &Session{}
	now := time.Now()
	result := self.db.Find(model,"`key` = ? And expire_time >= ?", key, now)
	if result.Error != nil {
		return nil, result.Error
	}
	return self.decode(model.Data)
}

func (self *GormStore) Set(key string, value map[string]interface{}, expire time.Duration) error {
	data, err := self.encode(value)
	if err != nil {
		return err
	}
	model := Session{
		Key: key,
		Data: data,
		ExpireTime: self.convertTime(expire),
	}
	result := self.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&model)
	return result.Error
}

func (self *GormStore) Delete(key string) {
	self.db.Delete(self.model, key)
}

func (self *GormStore) Exists(key string) bool {
	result := self.db.Select("key").Find(self.model, "`key` = ?", key)
	if result.RowsAffected >= 1 {
		return true
	}
	return false
}

//If record not exists, return time.Time Zero Value
func (self *GormStore) GetExpireTime(key string) time.Time {
	model := &Session{}
	self.db.Select("expire_time").Find(model, "`key` = ?", key)
	return model.ExpireTime
}

func (self *GormStore) SetExpireTime(key string, duration time.Duration) {
	self.db.Model(self.model).Where("`key` = ?", key).Update("expire_time", self.convertTime(duration))
}

func (self *GormStore) ClearExpired() {
	now := time.Now()
	self.db.Where("expire_time < ?", now).Delete(self.model)
}

func (self *GormStore) migrate() error {
	return self.db.AutoMigrate(self.model)
}

func (self *GormStore) encode(data map[string]interface{}) (string, error) {
	buf := bytes.Buffer{}
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(data)
	if err != nil {
		return "", err
	}
	s := base64.StdEncoding.EncodeToString(buf.Bytes())
	return s, nil
}

func (self *GormStore) decode(src string) (map[string]interface{}, error) {
	raw, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(raw)
	decoder := gob.NewDecoder(reader)
	data := make(map[string]interface{})
	if err = decoder.Decode(&data); err != nil && err != io.EOF {
		return nil, err
	}
	return data, nil
}

func (self *GormStore) convertTime(duration time.Duration) time.Time {
	now := time.Now()
	return now.Add(duration)
}
