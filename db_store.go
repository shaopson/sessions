package sessions

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type DBStore struct {
	db     *sql.DB
	driver string
	table  string
}

func NewDBStore(name string, db *sql.DB) *DBStore {
	table := name + "_session"
	driver := db.Driver()
	driverType := reflect.TypeOf(driver)
	if driverType.Kind() == reflect.Ptr {
		driverType = driverType.Elem()
	}
	driverName := driverType.String()
	lowerName := strings.ToLower(driverName)
	if strings.Contains(lowerName, "mysql") {
		if err := createByMysql(db, table); err != nil {
			panic(fmt.Errorf("NewDBStore: createTable error:%s", err))
		}
	} else if strings.Contains(lowerName, "sqlite") {
		if err := createBySqlite(db, table); err != nil {
			panic(fmt.Errorf("NewDBStore: createTable error:%s", err))
		}
	} else {
		panic(fmt.Errorf("NewDBStore: Not support database driver:%s", driverName))
	}
	return &DBStore{
		db:     db,
		driver: driverName,
		table:  table,
	}
}

func (store *DBStore) Get(key string) (StoreValue, error) {
	now := time.Now()
	query := fmt.Sprintf("SELECT session_data FROM %s WHERE session_key = ? AND expire_date >= ?", store.table)
	row := store.db.QueryRow(query, key, now)
	s := ""
	if err := row.Scan(&s); err != nil {
		if err == sql.ErrNoRows {
			err = DoesNotExists
		}
		return nil, err
	}
	value := StoreValue{}
	if err := DecodeString(s, &value); err != nil {
		return nil, InvalidData
	}
	return value, nil
}

func (store *DBStore) Set(key string, value StoreValue, expire int32) error {
	expireTime := time.Now().Add(time.Second * time.Duration(expire))
	update := fmt.Sprintf("UPDATE %s SET session_data = ?, expire_date = ? WHERE session_key = ?", store.table)
	data, err := EncodeToString(value)
	if err != nil {
		return err
	}
	if result, err := store.db.Exec(update, data, expireTime, key); err != nil {
		return err
	} else if n, err := result.RowsAffected(); err != nil || n > 0 {
		//update success or has error
		return err
	}
	//update affected 0 row, execute insert
	insert := fmt.Sprintf("INSERT INTO %s (session_key, session_data, expire_date) VALUES (?, ?, ?)", store.table)
	if _, err := store.db.Exec(insert, key, data, expireTime); err != nil {
		return err
	}
	return nil
}

func (store *DBStore) Delete(key string) {
	del := fmt.Sprintf("DELETE FROM %s WHERE session_key = ?", store.table)
	if _, err := store.db.Exec(del, key); err != nil {
		panic(err)
	}
}

func (store *DBStore) Exists(key string) bool {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE session_key = ?", store.table)
	row := store.db.QueryRow(query, key)
	temp := ""
	if err := row.Scan(&temp); err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return false
	}
	return true
}

func (store *DBStore) GetExpireTime(key string) (t time.Time) {
	query := fmt.Sprintf("SELECT expire_date FROM %s WHERE session_key = ?", store.table)
	row := store.db.QueryRow(query, key)
	if err := row.Scan(&t); err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	return
}

func (store *DBStore) SetExpireTime(key string, t time.Time) {
	update := fmt.Sprintf("UPDATE %s SET expire_date = ? WHERE session_key = ?", store.table)
	if _, err := store.db.Exec(update, t, key); err != nil {
		panic(err)
	}
}

func (store *DBStore) ClearExpired() {
	now := time.Now()
	go func() {
		del := fmt.Sprintf("DELETE FROM %s WHERE expire_date < ?", store.table)
		if _, err := store.db.Exec(del, now); err != nil {
			panic(err)
		}
	}()
}

func createBySqlite(db *sql.DB, table string) error {
	table = "\"" + table + "\""
	tableSQL := "CREATE TABLE IF NOT EXISTS " + table + `(
		"session_key" VARCHAR(40) NOT NULL,
        "session_data" TEXT NOT NULL,
		"expire_date" DATETIME NOT NULL,
		PRIMARY KEY ("session_key"));`

	indexSQL := `CREATE INDEX "expire_date_index" on ` + table + ` ("expire_date");`
	if result, err := db.Exec(tableSQL); err != nil {
		return err
	} else if n, _ := result.RowsAffected(); n <= 0 {
		return nil
	}
	if _, err := db.Exec(indexSQL); err != nil {
		return err
	}
	return nil
}

func createByMysql(db *sql.DB, table string) error {
	table = "`" + table + "`"
	tableSQL := "CREATE TABLE IF NOT EXISTS " + table
	tableSQL += "(`session_key` VARCHAR(40) NOT NULL, `session_data` TEXT NOT NULL, `expire_date` DATETIME NOT NULL, PRIMARY KEY(`session_key`), KEY `expire_date_index` (`expire_date`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	if _, err := db.Exec(tableSQL); err != nil {
		return err
	}
	return nil
}
