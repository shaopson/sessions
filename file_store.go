package sessions

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path"
	"time"
)

var MaxFileSize = 1024 * 16

type FileStore struct {
	Path string
}

func NewFileStore(path string) *FileStore {
	return &FileStore{
		Path: path,
	}
}

func (store *FileStore) Get(key string) (StoreValue, error) {
	now := time.Now()
	fileName := path.Join(store.Path, key)
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, DoesNotExists
		}
		return nil, err
	}
	defer file.Close()
	value := StoreValue{}
	buf := make([]byte, MaxFileSize)
	if n, err := file.Read(buf); err == nil {
		buf = buf[:n]
	} else if err == io.EOF {
		return value, nil
	} else {
		return nil, err
	}
	//在文件中过期时间和内容数据用"|"分隔   前部分为过期时间 后面为内容数据
	list := bytes.SplitN(buf, []byte("|"), 2)
	if len(list) != 2 {
		return nil, InvalidData
	}
	expire := time.Time{}
	//decode出错，跳过检查过期时间
	if err := expire.UnmarshalText(list[0]); err == nil {
		//内容已过期
		if expire.Before(now) {
			return value, nil
		}
	}

	if err := Decode(list[1], &value); err != nil {
		return nil, InvalidData
	}
	return value, nil
}

func (store *FileStore) Set(key string, value StoreValue, expire int32) (err error) {
	expireTime := time.Now().Add(time.Second * time.Duration(expire))
	fileName := path.Join(store.Path, key)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	// 先写入Text格式的过期时间
	data, _ := expireTime.MarshalText()
	data = append(data, '|') //过期时间和数据的分隔符
	if _, err = writer.Write(data); err != nil {
		return err
	}
	//写入数据部分
	data, err = Encode(value)
	if err != nil {
		return err
	}
	if _, err = writer.Write(data); err != nil {
		return err
	}
	return writer.Flush()
}

func (store *FileStore) Delete(key string) {
	fileName := path.Join(store.Path, key)
	os.Remove(fileName)
}

func (store *FileStore) Exists(key string) bool {
	fileName := path.Join(store.Path, key)
	if _, err := os.Stat(fileName); err != nil {
		return os.IsExist(err)
	}
	return true
}

func (store *FileStore) GetExpireTime(key string) (expire time.Time) {
	fileName := path.Join(store.Path, key)
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	if data, err := reader.ReadBytes('|'); err == nil {
		expire.UnmarshalText(data[:len(data)-1])
	}
	return
}

func (store *FileStore) SetExpireTime(key string, expire time.Time) {
	fileName := path.Join(store.Path, key)
	if expire.Year() > 9999 {
		expire = time.Date(9999, 0, 0, 0, 0, 0, 0, expire.Location())
	}
	file, err := os.OpenFile(fileName, os.O_RDWR, 0)
	if err != nil {
		return
	}
	defer file.Close()
	buf := make([]byte, MaxFileSize)
	n, err := file.Read(buf)
	if err != nil {
		return
	}
	buf = buf[:n]
	n = bytes.IndexByte(buf, '|')
	if n <= 0 {
		return
	}
	buf = buf[n:]
	exp, _ := expire.MarshalText()
	if err = file.Truncate(0); err != nil {
		return
	}
	if _, err = file.Seek(0, 0); err != nil {
		return
	}
	writer := bufio.NewWriter(file)
	if _, err := writer.Write(exp); err != nil {
		return
	}
	if _, err = writer.Write(buf); err == nil {
		writer.Flush()
	}
}
