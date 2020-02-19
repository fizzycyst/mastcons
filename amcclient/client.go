package main

import (
	"amcrpc/amcpb"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Master console client in operation... please wait...")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to %v", err)
	}

	defer cc.Close()

	c := amcpb.NewUserServiceClient(cc)

	doSingleUser(c, "Hypatia11102")
	doSingleUser(c, "Russell Ealing")
	doSingleUser(c, "Owen")
	doSingleUser(c, "Russell Ealing")

}

func doSingleUser(c amcpb.UserServiceClient, user string) {

	req := &amcpb.UserRequest{
		Username: user,
	}

	res, err := c.User(context.Background(), req)

	if err != nil {

		e, _ := status.FromError(err)

		if e.Message() == "sql: no rows in result set" {
			fmt.Printf("%v not found in database\n", user)
			return
		}

		log.Printf("%v", err)
		return
	}

	fmt.Printf("%v\n", res.User)
}
