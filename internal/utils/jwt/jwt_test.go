package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildtoken(t *testing.T) {
	type args struct {
		userID int
		key    string
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
				userID: 1,
				key:    "secret",
			},
			wantHeader: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			assertion:  assert.NoError,
		},
		{
			name: "TestEmptyKey",
			args: args{
				userID: 1,
				key:    "",
			},
			wantHeader: "",
			assertion:  assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildToken(tt.args.userID, tt.args.key)
			tt.assertion(t, err)
			assert.Contains(t, got, tt.wantHeader)
		})
	}
}
