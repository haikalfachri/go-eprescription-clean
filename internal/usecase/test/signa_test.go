package usecase_test

import (
	"context"
	"testing"

	"go-eprescription-clean/internal/entity"
	"go-eprescription-clean/internal/usecase/signa"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func signaUseCase(t *testing.T) (*signa.UseCase, *MockSignaRepo) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	repo := NewMockSignaRepo(mockCtl)

	useCase := signa.New(repo)

	return useCase, repo
}

func TestCreateSigna(t *testing.T) {
	t.Parallel()
	usecase, repo := signaUseCase(t)

	input := entity.Signa{Signa: "a.c", Description: "Before eat"}
	expected := &entity.Signa{ID: "123", Signa: "a.c", Description: "Before eat"}

	tests := []TestCase{
		{
			Name: "success",
			Mock: func() {
				repo.EXPECT().Create(context.Background(), input).Return(expected, nil)
			},
			Res: expected,
			Err: nil,
		},
		{
			Name: "repo error",
			Mock: func() {
				repo.EXPECT().Create(context.Background(), input).Return(nil, ErrInternalServErr)
			},
			Res: (*entity.Signa)(nil),
			Err: ErrInternalServErr,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			tc.Mock()
			Res, Err := usecase.Create(context.Background(), input)
			require.Equal(t, tc.Res, Res)
			require.ErrorIs(t, Err, tc.Err)
		})
	}
}

func TestGetAllSignas(t *testing.T) {
	t.Parallel()
	usecase, repo := signaUseCase(t)

	expected := []entity.Signa{
		{ID: "1", Signa: "a.c", Description: "Before eat"},
	}

	tests := []TestCase{
		{
			Name: "success",
			Mock: func() {
				repo.EXPECT().GetAll(context.Background()).Return(expected, nil)
			},
			Res: expected,
			Err: nil,
		},
		{
			Name: "repo error",
			Mock: func() {
				repo.EXPECT().GetAll(context.Background()).Return(nil, ErrInternalServErr)
			},
			Res: ([]entity.Signa)(nil),
			Err: ErrInternalServErr,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			tc.Mock()
			Res, Err := usecase.GetAll(context.Background())
			require.Equal(t, tc.Res, Res)
			require.ErrorIs(t, Err, tc.Err)
		})
	}
}

func TestGetByID(t *testing.T) {
	t.Parallel()
	usecase, repo := signaUseCase(t)

	id := "123"
	expected := &entity.Signa{ID: id, Signa: "a.c", Description: "Before eat"}

	tests := []TestCase{
		{
			Name: "success",
			Mock: func() {
				repo.EXPECT().GetByID(context.Background(), id).Return(expected, nil)
			},
			Res: expected,
			Err: nil,
		},
		{
			Name: "not found / repo error",
			Mock: func() {
				repo.EXPECT().GetByID(context.Background(), id).Return(nil, ErrInternalServErr)
			},
			Res: (*entity.Signa)(nil),
			Err: ErrInternalServErr,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			tc.Mock()
			Res, Err := usecase.GetByID(context.Background(), id)
			require.Equal(t, tc.Res, Res)
			require.ErrorIs(t, Err, tc.Err)
		})
	}
}

func TestUpdateSigna(t *testing.T) {
	t.Parallel()
	usecase, repo := signaUseCase(t)

	id := "123"
	input := entity.Signa{Signa: "a.c", Description: "After eat"}
	expected := &entity.Signa{ID: id, Signa: "a.c", Description: "After eat"}

	tests := []TestCase{
		{
			Name: "success",
			Mock: func() {
				repo.EXPECT().Update(context.Background(), id, input).Return(expected, nil)
			},
			Res: expected,
			Err: nil,
		},
		{
			Name: "repo error",
			Mock: func() {
				repo.EXPECT().Update(context.Background(), id, input).Return(nil, ErrInternalServErr)
			},
			Res: (*entity.Signa)(nil),
			Err: ErrInternalServErr,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			tc.Mock()
			Res, Err := usecase.Update(context.Background(), id, input)
			require.Equal(t, tc.Res, Res)
			require.ErrorIs(t, Err, tc.Err)
		})
	}
}

func TestDeleteSigna(t *testing.T) {
	t.Parallel()
	usecase, repo := signaUseCase(t)

	id := "123"

	tests := []TestCase{
		{
			Name: "success",
			Mock: func() {
				repo.EXPECT().Delete(context.Background(), id).Return(nil)
			},
			Res: nil,
			Err: nil,
		},
		{
			Name: "repo error",
			Mock: func() {
				repo.EXPECT().Delete(context.Background(), id).Return(ErrInternalServErr)
			},
			Res: nil,
			Err: ErrInternalServErr,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			tc.Mock()
			Err := usecase.Delete(context.Background(), id)
			require.ErrorIs(t, Err, tc.Err)
		})
	}
}
