package entities

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/mosuka/cete/protobuf"
	"github.com/vniche/users-microservice/datastore"
	"google.golang.org/grpc"
)

// User stands for a profile
type User struct {
	UID       string    `json:"uid,omitempty"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func SignUp(newUser *User) (string, error) {
	user := &User{
		UID:       uuid.New().String(),
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		CreatedAt: time.Now(),
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	_, err = datastore.Client.KVSClient.Set(context.Background(), &protobuf.SetRequest{
		Key:   user.UID,
		Value: userJson,
	})
	if err != nil {
		return "", err
	}

	return user.UID, nil
}

func List() ([]*User, error) {
	resp, err := datastore.Client.KVSClient.Scan(context.Background(), &protobuf.ScanRequest{
		Prefix: "",
	}, &grpc.EmptyCallOption{})
	if err != nil {
		return nil, err
	}

	var users []*User
	for _, curr := range resp.Values {
		user := &User{}
		err = json.Unmarshal(curr, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
