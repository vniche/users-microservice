package graphql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	"github.com/vniche/users-microservice/entities"
)

type Resolver struct{}

func (r *mutationResolver) Signup(ctx context.Context, input NewUser) (string, error) {
	return entities.SignUp(&entities.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
	})
}

func (r *queryResolver) List(ctx context.Context) ([]*User, error) {
	users, err := entities.List()
	if err != nil {
		return nil, err
	}

	parsed := make([]*User, len(users))
	for index, user := range users {
		parsed[index] = &User{
			UID:       user.UID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: user.CreatedAt.String(),
		}
	}

	return parsed, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
