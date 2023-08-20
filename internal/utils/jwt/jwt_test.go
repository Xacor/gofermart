package jwt

import (
	"testing"

	"github.com/Xacor/gophermart/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestBuildtoken(t *testing.T) {
	type args struct {
		user entity.User
		key  string
	}
	tests := []struct {
		name       string
		args       args
		wantHeader string
		assertion  assert.ErrorAssertionFunc
	}{
		{
			name: "TestOK",
			args: args{
				user: entity.User{
					ID:       1,
					Login:    "TestUser",
					Password: "$2a$10$qv/Omul7TF2rzhX6PZGZt.Ucg41V/88ew6jSm1oF70REYzvT0KcPm",
				},
				key: "secret",
			},
			wantHeader: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			assertion:  assert.NoError,
		},
		{
			name: "TestEmptyKey",
			args: args{
				user: entity.User{
					ID:       1,
					Login:    "TestUser",
					Password: "$2a$10$qv/Omul7TF2rzhX6PZGZt.Ucg41V/88ew6jSm1oF70REYzvT0KcPm",
				},
				key: "",
			},
			wantHeader: "",
			assertion:  assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildToken(tt.args.user, tt.args.key)
			tt.assertion(t, err)
			assert.Contains(t, got, tt.wantHeader)
		})
	}
}
