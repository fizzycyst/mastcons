package main

import (
	"amcrpc/amcpb"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"

	"google.golang.org/grpc"
)

type server struct{}

// Person ready for use within this server
type Person struct {
	username    string
	status      string
	communityID int32
}

func (*server) User(ctx context.Context, req *amcpb.UserRequest) (*amcpb.UserResponse, error) {
	fmt.Printf("User function run with %v", req)
	var nullPerson Person

	searchname := req.GetUsername()

	ourUser, err := searchDb(searchname)

	if err != nil {
		return nil, err
	}

	if ourUser == nullPerson {
		return nil, nil
	}

	res := &amcpb.User{
		Username:    ourUser.username,
		Status:      ourUser.status,
		CommunityId: ourUser.communityID,
	}
	return &amcpb.UserResponse{
		User: res,
	}, nil
}

func searchDb(user string) (Person, error) {
	var dummyPerson Person
	db, err := sql.Open("mysql",
		"root:goroutine@tcp(127.0.0.1:3306)/amc")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("We have a successful database connection yeehaa!\n")
	}
	defer db.Close()

	var username string
	var status string
	var communityID int32

	err = db.QueryRow("select username, status, communityId from users where username = ?", user).Scan(&username, &status, &communityID)

	if err != nil {
		if err == sql.ErrNoRows {
			//log.Fatalf("There were no rows returned %v", err)
			return dummyPerson, nil
		} else {
			log.Fatal(err)
		}
	}

	//defer rows.Close()

	var newPerson Person
	newPerson.username = username
	newPerson.status = status
	newPerson.communityID = communityID

	return newPerson, nil

}

func main() {

	fmt.Println("Master Console Server Running...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Fatal to listen: %v", err)
	}

	s := grpc.NewServer()

	amcpb.RegisterUserServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}

}
