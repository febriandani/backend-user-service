// ./cmd/server/main.go

package main

import (
	"fmt"
	"net"

	"github.com/febriandani/backend-user-service/internal/api"
	database "github.com/febriandani/backend-user-service/internal/db"
	infra "github.com/febriandani/backend-user-service/internal/infra"
	"github.com/febriandani/backend-user-service/protogen/golang/users"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	conf, err := getConfigKey()
	if err != nil {
		panic(err)
	}

	db, log, dblist, err := newDbContext(conf)
	if err != nil {
		panic(err)
	}

	// create a TCP listener on the specified port
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", viper.GetString("APP.PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server instance
	server := grpc.NewServer()

	userService := api.NewUserService(db, log, dblist, conf)

	// register the user service with the grpc server
	users.RegisterUsersServer(server, &userService)

	// start listening to requests
	log.Printf("server listening at %v", listener.Addr())
	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getConfigKey() (*infra.AppService, error) {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	conf := infra.AppService{
		App: infra.AppUser{
			Name:         viper.GetString("APP.NAME"),
			Environtment: viper.GetString("APP.ENV"),
			URL:          viper.GetString("APP.URL"),
			Port:         viper.GetString("APP.PORT"),
			SecretKey:    viper.GetString("APP.KEY"),
		},
		Route: infra.RouteUser{
			Methods: viper.GetStringSlice("ROUTE.METHODS"),
			Headers: viper.GetStringSlice("ROUTE.HEADERS"),
			Origins: viper.GetStringSlice("ROUTE.ORIGIN"),
		},
		DatabaseUser: infra.DatabaseUser{
			Read: infra.DBDetailUser{
				Username:     viper.GetString("DATABASE.READ.USERNAME"),
				Password:     viper.GetString("DATABASE.READ.PASSWORD"),
				URL:          viper.GetString("DATABASE.READ.URL"),
				Port:         viper.GetString("DATABASE.READ.PORT"),
				DBName:       viper.GetString("DATABASE.READ.DB_NAME"),
				MaxIdleConns: viper.GetInt("DATABASE.READ.MAXIDLECONNS"),
				MaxOpenConns: viper.GetInt("DATABASE.READ.MAXOPENCONNS"),
				MaxLifeTime:  viper.GetInt("DATABASE.READ.MAXLIFETIME"),
				Timeout:      viper.GetString("DATABASE.READ.TIMEOUT"),
				SSLMode:      viper.GetString("DATABASE.READ.SSL_MODE"),
			},
			Write: infra.DBDetailUser{
				Username:     viper.GetString("DATABASE.WRITE.USERNAME"),
				Password:     viper.GetString("DATABASE.WRITE.PASSWORD"),
				URL:          viper.GetString("DATABASE.WRITE.URL"),
				Port:         viper.GetString("DATABASE.WRITE.PORT"),
				DBName:       viper.GetString("DATABASE.WRITE.DB_NAME"),
				MaxIdleConns: viper.GetInt("DATABASE.WRITE.MAXIDLECONNS"),
				MaxOpenConns: viper.GetInt("DATABASE.WRITE.MAXOPENCONNS"),
				MaxLifeTime:  viper.GetInt("DATABASE.WRITE.MAXLIFETIME"),
				Timeout:      viper.GetString("DATABASE.WRITE.TIMEOUT"),
				SSLMode:      viper.GetString("DATABASE.WRITE.SSL_MODE"),
			},
		},
		Redis: infra.RedisUser{
			Username:     viper.GetString("REDIS.USERNAME"),
			Password:     viper.GetString("REDIS.PASSWORD"),
			URL:          viper.GetString("REDIS.URL"),
			Port:         viper.GetInt("REDIS.PORT"),
			MinIdleConns: viper.GetInt("REDIS.MINIDLECONNS"),
			Timeout:      viper.GetString("REDIS.TIMEOUT"),
		},
		Authorization: infra.AuthUser{
			JWT: infra.JWTCredential{
				IsActive:              viper.GetBool("AUTHORIZATION.JWT.IS_ACTIVE"),
				AccessTokenSecretKey:  viper.GetString("AUTHORIZATION.JWT.ACCESS_TOKEN_SECRET_KEY"),
				AccessTokenDuration:   viper.GetInt("AUTHORIZATION.JWT.ACCESS_TOKEN_DURATION"),
				RefreshTokenSecretKey: viper.GetString("AUTHORIZATION.JWT.REFRESH_TOKEN_SECRET_KEY"),
				RefreshTokenDuration:  viper.GetInt("AUTHORIZATION.JWT.REFRESH_TOKEN_DURATION"),
			},
			Public: infra.PublicCredential{
				SecretKey: viper.GetString("AUTHORIZATION.PUBLIC.SECRET_KEY"),
			},
		},
		KeyData: infra.KeyUser{
			User: viper.GetString("KEY.USER"),
		},
		Minio: infra.MinioSecret{
			BucketName: viper.GetString("MINIO.BUCKET_NAME"),
			Endpoint:   viper.GetString("MINIO.ENDPOINT"),
			Key:        viper.GetString("MINIO.KEY"),
			Secret:     viper.GetString("MINIO.SECRET"),
			Region:     viper.GetString("MINIO.REGION"),
			TempFolder: viper.GetString("MINIO.TEMP_FOLDER"),
			BaseURL:    viper.GetString("MINIO.BASE_URL"),
		},
	}

	return &conf, nil
}

func newDbContext(conf *infra.AppService) (*database.DB, *logrus.Logger, *infra.DatabaseList, error) {
	// Init Log
	logger := infra.NewLogger(conf)

	// Init DB Read Connection.
	dbRead := infra.NewDB(logger)
	dbRead.ConnectDB(&conf.DatabaseUser.Read)
	if dbRead.Err != nil {
		return &database.DB{}, logger, nil, dbRead.Err
	}

	// Init DB Write Connection.
	dbWrite := infra.NewDB(logger)
	dbWrite.ConnectDB(&conf.DatabaseUser.Write)
	if dbWrite.Err != nil {
		return &database.DB{}, logger, nil, dbWrite.Err
	}

	dbList := &infra.DatabaseList{
		Backend: infra.DatabaseType{
			Read:  &dbRead,
			Write: &dbWrite,
		},
	}

	// Init Minio Config.
	// media, err := infra.NewMinio(conf.Minio)
	// if err != nil {
	// 	return handler, logger, err
	// }

	// redis, err := infra.NewRedisClient(conf.Redis)
	// if err != nil {
	// 	return handler, logger, err
	// }

	db := database.NewDB(dbList, logger)

	return db, logger, dbList, dbRead.Err
}
