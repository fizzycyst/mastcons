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

//NoRowsError Exported for fun
type NoRowsError struct {
	msg string
}

// Person ready for use within this server
type Person struct {
	username    string
	status      string
	communityID int32
}

func (error *NoRowsError) Error() string {
	return error.msg
}

func noRows() error {
	return &NoRowsError{"User Does Not Exist"}
}
func (*server) User(ctx context.Context, req *amcpb.UserRequest) (*amcpb.UserResponse, error) {
	fmt.Printf("User function run with %v\n", req)

	searchname := req.GetUsername()

	ourUser, err := searchDb(searchname)

	if err == noRows() {
		log.Printf("We finally have no rows caught")
	}

	if err != nil {
		log.Printf("Checking error %v", err)
		return nil, err
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

func searchDb(user string) (*Person, error) {
	db, err := sql.Open("mysql",
		"root:71ODp4rXjNmr0fJ0@tcp(10.10.0.9:3306)/ums2") // Not for production..
	/*
		db, err := sql.Open("mysql",
			"root:goroutine@tcp(127.0.0.1:3306)/amc")
	*/
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
	var rootCommunityID int32

	err = db.QueryRow("select username, status, rootCommunityId from Users where username = ?", user).Scan(&username, &status, &rootCommunityID)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("There were no rows returned %v", noRows())
			return nil, noRows()
		}
		log.Fatal(err)
	}

	//defer rows.Close()

	var newPerson Person
	newPerson.username = username
	newPerson.status = status
	newPerson.communityID = rootCommunityID
	return &newPerson, nil

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
