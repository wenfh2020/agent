package db

import (
	"errors"
	"github.com/jinzhu/gorm"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/thinkboy/log4go"
)

var (
	// find database failed
	ErrCantFindDatabase = errors.New("find db failed")

	// init database mgr failed.
	ErrNotInitDatabaseMgr = errors.New("not init database mgr")

	// add db info duplicate.
	ErrDuplicate = errors.New("duplicate")

	// add db info duplicate.
	ErrInvalidDbParams = errors.New("invalid db's params")
)

// database info.
type DbConfig struct {
	name        string
	url         string
	maxOpenConn int
	maxIdleConn int
}

// database info.
type DbInfo struct {
	config *DbConfig
	db     *gorm.DB
}

// database manager
type DbMgr struct {
	dbMap map[string]*DbInfo
	mutex *sync.Mutex
}

var dbMgr *DbMgr

// init database manager.
func InitDbMgr() {
	dbMgr = newDbMgr()
}

// uninit database manager.
func UninitSqlMgr() {
	if dbMgr != nil {
		dbMgr.Close()
		dbMgr = nil
	}
}

func AddDbInfo(name, url string, maxIdleConn, maxOpenConn int) error {
	if len(name) == 0 || len(url) == 0 || maxIdleConn == 0 || maxOpenConn == 0 {
		return ErrInvalidDbParams
	}

	log.Info("insert db info, name: %v, url: %v, max-idle-conn: %v, max-open-conn: %v",
		name, url, maxIdleConn, maxOpenConn)

	dbi := &DbInfo{
		db: nil,
		config: &DbConfig{
			name:        name,
			url:         url,
			maxIdleConn: maxIdleConn,
			maxOpenConn: maxOpenConn,
		},
	}
	return dbMgr.addDbInfo(dbi)
}

func GetDB(name string) (*gorm.DB, error) {
	if dbMgr == nil {
		panic(ErrNotInitDatabaseMgr)
	}
	return dbMgr.getDB(name)
}

func newDbMgr() *DbMgr {
	mgr := &DbMgr{
		dbMap: make(map[string]*DbInfo),
		mutex: &sync.Mutex{},
	}
	return mgr
}

func (mgr *DbMgr) addDbInfo(dbi *DbInfo) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	_, ok := mgr.dbMap[dbi.config.name]
	if ok {
		log.Error("add db info duplicate! name: %v", dbi.config.name)
		return ErrDuplicate
	} else {
		mgr.dbMap[dbi.config.name] = dbi
		return nil
	}
}

func (mgr *DbMgr) getDB(name string) (*gorm.DB, error) {
	mgr.mutex.Lock()
	dbi, ok := mgr.dbMap[name]
	if !ok {
		mgr.mutex.Unlock()
		return nil, ErrCantFindDatabase
	}
	mgr.mutex.Unlock()

	if dbi.db != nil {
		return dbi.db, nil
	} else {
		db, err := mgr.initDB(dbi)
		if err != nil {
			return nil, err
		}

		if err = mgr.setDB(name, db); err != nil {
			return nil, err
		}

		return db, nil
	}
}

func (mgr *DbMgr) setDB(name string, db *gorm.DB) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	dbi, ok := mgr.dbMap[name]
	if !ok {
		return ErrCantFindDatabase
	}
	dbi.db = db
	return nil
}

func (mgr *DbMgr) initDB(dbi *DbInfo) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", dbi.config.url)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)
	db.DB().SetConnMaxLifetime(2 * time.Hour)

	maxIdleConn := dbi.config.maxIdleConn
	if maxIdleConn != 0 {
		db.DB().SetMaxIdleConns(maxIdleConn)
	}
	maxOpenConn := dbi.config.maxOpenConn
	if maxOpenConn != 0 {
		db.DB().SetMaxOpenConns(maxOpenConn)
	}

	if err := db.DB().Ping(); err != nil {
		return db, err
	}

	log.Info("create mysql connection, name: %s, url: %s, maxIdleConn:%d, maxOpenConn: %d",
		dbi.config.name, dbi.config.url, maxIdleConn, maxOpenConn)

	dbi.db = db
	return dbi.db, nil
}

func (mgr *DbMgr) Close() {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	for _, dbi := range mgr.dbMap {
		dbi.db.Close()
	}

	/* clear map. */
	mgr.dbMap = make(map[string]*DbInfo)
}
