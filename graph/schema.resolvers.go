package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal"
	"github.com/carlqt/ezsplit/internal/auth"
	"github.com/carlqt/ezsplit/internal/repository"
	"github.com/go-jet/jet/v2/qrm"
)

// SharedBy is the resolver for the sharedBy field.
func (r *itemResolver) SharedBy(ctx context.Context, obj *model.Item) ([]*model.User, error) {
	users, err := r.Repositories.UserOrdersRepository.SelectAllUsersFromItem(obj.ID)
	if err != nil {
		return nil, err
	}

	var modelUsers []*model.User
	for _, user := range users {
		modelUser := newModelUser(user.ID, user.Name, user.IsVerified())
		modelUsers = append(modelUsers, modelUser)
	}

	return modelUsers, nil
}

// AddItemToReceipt is the resolver for the addItemToReceipt field.
func (r *mutationResolver) AddItemToReceipt(ctx context.Context, input *model.AddItemToReceiptInput) (*model.Item, error) {
	price := toPriceCents(*input.Price)

	item := repository.Item{}
	item.Name = &input.Name
	item.ReceiptID = repository.BigInt(input.ReceiptID)
	item.Price = price

	err := r.Repositories.ItemRepository.Create(&item)
	if err != nil {
		return nil, errors.New("failed to add the item")
	}

	return newModelItem(item), nil
}

// AssignUserToItem is the resolver for the assignUserToItem field.
func (r *mutationResolver) AssignUserToItem(ctx context.Context, input *model.AssignUserToItemInput) (*model.Item, error) {
	panic(fmt.Errorf("not implemented: AssignUserToItem - assignUserToItem"))
}

// AssignMeToItem is the resolver for the assignMeToItem field.
func (r *mutationResolver) AssignMeToItem(ctx context.Context, input *model.AssignOrDeleteMeToItemInput) (*model.Item, error) {
	// TODO: Need to refactor since it's really huge

	userClaim := ctx.Value(auth.UserClaimKey).(auth.UserClaim)

	// Create a New UserOrder using ItemID and userID
	// Need to return an item or user object?
	err := r.Repositories.UserOrdersRepository.Create(userClaim.ID, input.ItemID)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	// fetch item by ID
	item, err := r.Repositories.ItemRepository.FindByID(input.ItemID)
	if err != nil {
		slog.Error(err.Error())
		return nil, errors.New("failed to you to item")
	}

	return newModelItem(item), nil
}

// RemoveMeFromItem is the resolver for the removeMeFromItem field.
func (r *mutationResolver) RemoveMeFromItem(ctx context.Context, input *model.AssignOrDeleteMeToItemInput) (*model.DeleteItemPayload, error) {
	userClaim := ctx.Value(auth.UserClaimKey).(auth.UserClaim)

	err := r.Repositories.UserOrdersRepository.Delete(userClaim.ID, input.ItemID)
	if err != nil {
		return nil, errors.New("failed to unassign user")
	}

	// fetch item by ID
	return &model.DeleteItemPayload{
		ID:  input.ItemID,
		Msg: "Item removed",
	}, nil
}

// AssignOrRemoveMeFromItem is the resolver for the assignOrRemoveMeFromItem field.
func (r *mutationResolver) AssignOrRemoveMeFromItem(ctx context.Context, itemID string) (*model.UserOrderRef, error) {
	userClaim := ctx.Value(auth.UserClaimKey).(auth.UserClaim)

	_, err := r.Repositories.UserOrdersRepository.FindByUserIDAndItemID(userClaim.ID, itemID)

	// IF assigned, remove user from item (Deleting)
	// nil error means user order is found
	if err == nil {
		e := r.Repositories.UserOrdersRepository.Delete(userClaim.ID, itemID)

		if e != nil {
			slog.Error("failed to unselect self to item", "error", e.Error(), "itemID", itemID, "userID", userClaim.ID)
			return nil, fmt.Errorf("failed to unassign user from item")
		}

		return newUserOrderRef(userClaim.ID, itemID), nil
	}

	// If error return is ErrNoRows, create/assign a user order
	if errors.Is(err, qrm.ErrNoRows) {
		e := r.Repositories.UserOrdersRepository.Create(userClaim.ID, itemID)
		if e != nil {
			slog.Error("failed to assign self to item", "error", e.Error(), "itemID", itemID, "userID", userClaim.ID)
			return nil, fmt.Errorf("failed to assign self to item")
		}

		return newUserOrderRef(userClaim.ID, itemID), nil
	}

	// IF not assigned, assign user
	return nil, err
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input *model.UserInput) (*model.Me, error) {
	if input.ConfirmPassword != input.Password {
		return nil, errors.New("password doesn't match confirm password")
	}

	password, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user, err := r.Repositories.UserRepository.CreateWithAccount(input.Username, password)
	if err != nil {
		slog.Error(err.Error())
		return nil, errors.New("user already exists")
	}

	userClaim := auth.NewUserClaim(user.ID, user.Name, user.IsVerified())
	signedToken, err := auth.CreateAndSignToken(userClaim, r.Config.JWTSecret)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error signing token")
	}

	setCookieFn, ok := ctx.Value(internal.ContextKeySetCookie).(func(*http.Cookie))
	if !ok {
		return nil, errors.New("error setting cookie")
	}

	setCookieFn(&http.Cookie{
		Name:     string(internal.JWTCookie),
		Value:    signedToken,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	return newModelMe(
		user.ID,
		user.Name,
		user.IsVerified(),
	), nil
}

// CreateGuestUser is the resolver for the createGuestUser field.
func (r *mutationResolver) CreateGuestUser(ctx context.Context, input *model.CreateGuestUserInput) (*model.User, error) {
	if input == nil || input.Username == "" {
		return nil, errors.New("name required")
	}

	user, err := r.Repositories.UserRepository.CreateGuest(input.Username)
	if err != nil {
		return nil, err
	}

	userClaim := auth.NewUserClaim(user.ID, user.Name, user.IsVerified())
	signedToken, err := auth.CreateAndSignToken(userClaim, r.Config.JWTSecret)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error signing token")
	}

	setCookieFn, ok := ctx.Value(internal.ContextKeySetCookie).(func(*http.Cookie))
	if !ok {
		return nil, errors.New("error setting cookie")
	}

	setCookieFn(&http.Cookie{
		Name:     string(internal.JWTCookie),
		Value:    signedToken,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	return newModelUser(
		user.ID,
		user.Name,
		user.IsVerified(),
	), nil
}

// LoginUser is the resolver for the loginUser field.
func (r *mutationResolver) LoginUser(ctx context.Context, input *model.LoginUserInput) (*model.Me, error) {
	if input == nil {
		slog.Warn("input is nil")
		return nil, errors.New("incorrect username or password")
	}

	user, err := r.Repositories.UserRepository.FindVerifiedByUsername(input.Username)
	if err != nil {
		slog.Warn(err.Error())
	}

	ok := auth.ComparePassword(user.Account.Password, input.Password)
	if !ok {
		return nil, errors.New("incorrect username or password")
	}

	userClaim := auth.NewUserClaim(user.ID, user.Name, user.IsVerified())
	signedToken, err := auth.CreateAndSignToken(userClaim, r.Config.JWTSecret)
	if err != nil {
		slog.Error(err.Error())
		return nil, errors.New("something went wrong")
	}

	setCookieFn, ok := ctx.Value(internal.ContextKeySetCookie).(func(*http.Cookie))
	if !ok {
		slog.Error("cookie type assertion failed")
		return nil, errors.New("something went wrong")
	}

	setCookieFn(&http.Cookie{
		Name:     string(internal.JWTCookie),
		Value:    signedToken,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	return newModelMe(
		user.ID,
		user.Name,
		user.IsVerified(),
	), nil
}

// LogoutUser is the resolver for the logoutUser field.
func (r *mutationResolver) LogoutUser(ctx context.Context) (string, error) {
	setCookieFn, ok := ctx.Value(internal.ContextKeySetCookie).(func(*http.Cookie))
	if !ok {
		return "", errors.New("error setting cookie")
	}

	setCookieFn(&http.Cookie{
		Name:     string(internal.JWTCookie),
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   0,
		Expires:  time.Unix(0, 0),
	})

	return "ok", nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	users, err := r.Repositories.UserRepository.GetAllUsers()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	var modelUsers []*model.User
	for _, user := range users {
		modelUser := newModelUser(
			user.ID,
			user.Name,
			user.IsVerified(),
		)
		modelUsers = append(modelUsers, modelUser)
	}

	return modelUsers, nil
}

// Item returns ItemResolver implementation.
func (r *Resolver) Item() ItemResolver { return &itemResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type itemResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
