// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysql

import (
	"fmt"
	"sync"

	"github.com/marmotedu/iam/pkg/db"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/marmotedu/goserver/internal/goserver/store"
	"github.com/marmotedu/goserver/internal/pkg/logger"
	v1 "github.com/marmotedu/goserver/internal/pkg/model/goserver/v1"
)

type datastore struct {
	db *gorm.DB

	// can include two database instance if needed
	// docker *grom.DB
	// db *gorm.DB
}

func (ds *datastore) Users() store.UserStore {
	return newUsers(ds)
}

func (ds *datastore) Posts() store.PostStore {
	return newPosts(ds)
}

func (ds *datastore) Close() error {
	if ds.db == nil {
		return nil
	}

	db, err := ds.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

var (
	mysqlFactory store.Factory
	once         sync.Once
)

// GetMySQLFactoryOr create mysql factory with the given config.
func GetMySQLFactoryOr() (store.Factory, error) {
	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		options := &db.Options{
			Host:                  viper.GetString("db.host"),
			Username:              viper.GetString("db.username"),
			Password:              viper.GetString("db.password"),
			Database:              viper.GetString("db.database"),
			MaxIdleConnections:    viper.GetInt("db.max-idle-connections"),
			MaxOpenConnections:    viper.GetInt("db.max-open-connections"),
			MaxConnectionLifeTime: viper.GetDuration("db.max-connection-life-time"),
			LogLevel:              viper.GetInt("db.log-level"),
			Logger:                logger.New(viper.GetInt("db.log-level")),
		}
		dbIns, err = db.New(options)

		// uncomment the following line if you need auto migration the given models
		// not suggested in production environment.
		// migrateDatabase(dbIns)

		mysqlFactory = &datastore{dbIns}
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}

	return mysqlFactory, nil
}

// cleanDatabase tear downs the database tables.
// nolint:unused // may be reused in the feature, or just show a migrate usage.
func cleanDatabase(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&v1.User{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&v1.Post{}); err != nil {
		return err
	}

	return nil
}

// migrateDatabase run auto migration for given models, will only add missing fields,
// won't delete/change current data.
// nolint:unused // may be reused in the feature, or just show a migrate usage.
func migrateDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(&v1.User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&v1.Post{}); err != nil {
		return err
	}

	return nil
}

// resetDatabase resets the database tables.
// nolint:unused,deadcode // may be reused in the feature, or just show a migrate usage.
func resetDatabase(db *gorm.DB) error {
	if err := cleanDatabase(db); err != nil {
		return err
	}
	if err := migrateDatabase(db); err != nil {
		return err
	}

	return nil
}
