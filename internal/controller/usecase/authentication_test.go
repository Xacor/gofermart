package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Xacor/gophermart/internal/entity"
	gomock "github.com/golang/mock/gomock"
)

func TestNewAuthUseCase(t *testing.T) {
	type args struct {
		repo      UserRepo
		secretKey string
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockUserRepo(ctrl)

	tests := []struct {
		name string
		args args
		want *AuthUseCase
	}{
		{
			name: "TestOK",
			args: args{
				repo:      m,
				secretKey: "verysecretkey",
			},
			want: &AuthUseCase{
				repo:      m,
				secretKey: "verysecretkey",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthUseCase(tt.args.repo, tt.args.secretKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthUseCase_Register(t *testing.T) {
	type fields struct {
		repo      *MockUserRepo
		secretKey string
	}
	type args struct {
		ctx  context.Context
		user entity.User
	}

	tests := []struct {
		name    string
		prepare func(f *fields, a *args)
		args    args
		wantErr bool
	}{
		{
			name: "TestNewUser",
			prepare: func(f *fields, a *args) {
				gomock.InOrder(
					f.repo.EXPECT().GetByLogin(a.ctx, a.user.Login).Return(entity.User{}, errors.New("any error")),
					f.repo.EXPECT().CreateUser(a.ctx, gomock.Any()).Return(nil),
				)
			},
			args: args{
				context.Background(),
				entity.User{
					ID:       1,
					Login:    "TestUser1",
					Password: "password",
				}},

			wantErr: false,
		},
		{
			name: "TestUserExists",
			prepare: func(f *fields, a *args) {
				f.repo.EXPECT().GetByLogin(a.ctx, a.user.Login).Return(a.user, nil)
			},
			args: args{
				context.Background(),
				entity.User{
					ID:       2,
					Login:    "TestUser2",
					Password: "password",
				}},
			wantErr: true,
		},
		{
			name: "TestLongPassword",
			prepare: func(f *fields, a *args) {
				f.repo.EXPECT().GetByLogin(a.ctx, a.user.Login).Return(entity.User{}, errors.New("any error"))
			},
			args: args{
				context.Background(),
				entity.User{
					ID:       3,
					Login:    "TestUser3",
					Password: string(make([]byte, 100)),
				}},
			wantErr: true,
		},
		{
			name: "TestFailedToCreate",
			prepare: func(f *fields, a *args) {
				gomock.InOrder(
					f.repo.EXPECT().GetByLogin(a.ctx, a.user.Login).Return(entity.User{}, errors.New("any error")),
					f.repo.EXPECT().CreateUser(a.ctx, gomock.Any()).Return(errors.New("any error")),
				)
			},
			args: args{
				context.Background(),
				entity.User{
					ID:       4,
					Login:    "TestUser4",
					Password: "password",
				}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{NewMockUserRepo(ctrl), "verysecretkey"}

			if tt.prepare != nil {
				tt.prepare(&f, &tt.args)
			}

			a := &AuthUseCase{
				repo:      f.repo,
				secretKey: f.secretKey,
			}
			if err := a.Register(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("AuthUseCase.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthUseCase_Authenticate(t *testing.T) {
	type fields struct {
		repo      UserRepo
		secretKey string
	}
	type args struct {
		ctx  context.Context
		user entity.User
	}
	tests := []struct {
		name    string
		prepare func(f *fields, a *args)
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{NewMockUserRepo(ctrl), "verysecretkey"}

			a := &AuthUseCase{
				repo:      f.repo,
				secretKey: f.secretKey,
			}
			got, err := a.Authenticate(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthUseCase.Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthUseCase.Authenticate() = %v, want %v", got, tt.want)
			}
		})
	}
}
