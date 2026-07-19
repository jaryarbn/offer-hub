package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"offer-hub/backend/src/config"
)

const connectionTimeout = 10 * time.Second

var globalData *Data

type Data struct {
	MySQL   *gorm.DB
	MongoDB *mongo.Database

	mongoClient *mongo.Client
}

func NewData(conf *config.TomlConfig) (*Data, error) {
	if conf == nil {
		return nil, errors.New("config is not initialized")
	}

	mysqlDB, err := gorm.Open(mysql.Open(conf.MySQL.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to MySQL: %w", err)
	}

	mongoOptions := options.Client().
		ApplyURI(conf.MongoDB.URL).
		SetMaxPoolSize(conf.MongoDB.MaxPool).
		SetMinPoolSize(conf.MongoDB.MinPool)
	if conf.MongoDB.User != "" {
		mongoOptions.SetAuth(options.Credential{
			Username: conf.MongoDB.User,
			Password: conf.MongoDB.Pwd,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, mongoOptions)
	if err != nil {
		closeMySQL(mysqlDB)
		return nil, fmt.Errorf("connect to MongoDB: %w", err)
	}
	if err := mongoClient.Ping(ctx, nil); err != nil {
		_ = mongoClient.Disconnect(ctx)
		closeMySQL(mysqlDB)
		return nil, fmt.Errorf("ping MongoDB: %w", err)
	}

	initialized := &Data{
		MySQL:       mysqlDB,
		MongoDB:     mongoClient.Database(conf.MongoDB.Database),
		mongoClient: mongoClient,
	}
	globalData = initialized
	return initialized, nil
}

func GetData() *Data {
	return globalData
}

func (data *Data) Close() error {
	if data == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	var mysqlErr error
	if data.MySQL != nil {
		sqlDB, err := data.MySQL.DB()
		if err != nil {
			mysqlErr = err
		} else {
			mysqlErr = sqlDB.Close()
		}
	}

	var mongoErr error
	if data.mongoClient != nil {
		mongoErr = data.mongoClient.Disconnect(ctx)
	}
	if globalData == data {
		globalData = nil
	}

	return errors.Join(mysqlErr, mongoErr)
}

func (data *Data) Ping(ctx context.Context) map[string]string {
	status := map[string]string{
		"mysql":   "up",
		"mongodb": "up",
	}

	sqlDB, err := data.MySQL.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		status["mysql"] = "down"
	}
	if err := data.mongoClient.Ping(ctx, nil); err != nil {
		status["mongodb"] = "down"
	}

	return status
}

func closeMySQL(database *gorm.DB) {
	sqlDB, err := database.DB()
	if err == nil {
		_ = sqlDB.Close()
	}
}
