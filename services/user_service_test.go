package service_test

import (
	"testing"

	"github.com/davidyannick/repository-pattern/domain"
	mock_repository "github.com/davidyannick/repository-pattern/mocks"
	service "github.com/davidyannick/repository-pattern/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
)

func TestGetAll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(mockCtrl)
	mockRepo.EXPECT().GetAllUsers(gomock.Any()).Return([]domain.User{
		{ID: uuid.New(), Name: "John Doe", Email: "john.doe@example.com"},
		{ID: uuid.New(), Name: "Jane Smith", Email: "jane.smith@example.com"},
	}, nil)

	userService := service.NewUserService(mockRepo)

	users, err := userService.GetAllUsers(t.Context())

	require.NoError(t, err)
	require.Len(t, users, 2)
}

func TestAddUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(mockCtrl)
	user := domain.User{ID: uuid.New(), Name: "John Doe", Email: "user@email.com"}
	mockRepo.EXPECT().AddUser(gomock.Any(), user).Return(&user, nil)

	userService := service.NewUserService(mockRepo)

	userAdded, err := userService.AddUser(t.Context(), user)
	require.NoError(t, err)
	require.Equal(t, userAdded.Name, user.Name)
	require.Equal(t, userAdded.Email, user.Email)
	require.Equal(t, userAdded.ID, user.ID)
}
