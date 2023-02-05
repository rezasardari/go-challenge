package main

import (
	"estimation-service/external/pg"
	cache "estimation-service/external/redis"
	"estimation-service/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

func main() {
	setUpViper()
	runDbMigrations()
	db := getGormDb()
	redis := getRedisClient()
	repository := pg.RepositoryImpl{
		DB: db,
	}
	redisClient := cache.ClientImpl{
		Rdb:      redis,
		Duration: time.Hour,
	}
	service := pkg.ServiceImpl{
		Repository: &repository,
		Redis:      &redisClient,
	}

	controller := pkg.Controller{
		Service: &service,
	}

	router := gin.New()
	controller.SetRoutes(router)
	router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/health"),
		gin.Recovery(),
	)

	router.Run(viper.GetString("serverPort"))

}

func setUpViper() {
	viper.SetConfigName(getEnv("CONFIG_NAME", "dev-conf"))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func runDbMigrations() {
	m, err := migrate.New(
		"file://db/migrations",
		viper.GetString("pgMigrationSource"))
	if err != nil {
		log.Error(errors.Wrap(err, "failed to find migration files"))
		panic("failed to find migration files")
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error(errors.Wrap(err, "failed to apply migrations"))
		panic("failed to apply migrations")
	}

}

func getGormDb() *gorm.DB {
	db, err := gorm.Open(postgres.Open(viper.GetString("postgresSource")), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Error(errors.Wrap(err, "failed to initial gorm DB"))
		panic("failed to initial gorm DB")
	}
	return db
}

func getRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}
