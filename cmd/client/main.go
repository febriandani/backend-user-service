// ./cmd/client/main.go

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/febriandani/backend-user-service/protogen/golang/users"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("could not connect to user service: %v", err)
	}
	userServiceAddr := fmt.Sprintf("0.0.0.0:%s", viper.GetString("APP.PORT"))

	// Set up a connection to the user server.
	fmt.Println("Connecting to user service via", userServiceAddr)
	conn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to user service: %v", err)
	}
	defer conn.Close()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	if err = users.RegisterUsersHandler(context.Background(), mux, conn); err != nil {
		log.Fatalf("failed to register the user server: %v", err)
	}

	addr := fmt.Sprintf("0.0.0.0:%s", viper.GetString("APP.PORT_CLIENT"))
	// start listening to requests from the gateway server
	fmt.Println("API gateway server is running on " + addr)
	if err = http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("gateway server closed abruptly: ", err)
	}
}
