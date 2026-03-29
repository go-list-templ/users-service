package entity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	type args struct {
		name     *string
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success - create user",
			args: args{
				name:     stringPtr("test"),
				email:    "test@example.com",
				password: "password",
			},
			wantErr: false,
		},
		{
			name: "success - russian domain email",
			args: args{
				name:     stringPtr("test"),
				email:    "пользователь@компания.рф",
				password: "password",
			},
			wantErr: false,
		},
		{
			name: "success - email with one domain",
			args: args{
				name:     stringPtr("test"),
				email:    "invalid@email",
				password: "password",
			},
			wantErr: false,
		},
		{
			name: "success - empty name",
			args: args{
				name:     nil,
				email:    "example@example.com",
				password: "password",
			},
			wantErr: false,
		},
		{
			name: "fail - min length name",
			args: args{
				name:     stringPtr("t"),
				email:    "test@example.com",
				password: "password",
			},
			wantErr: true,
		},
		{
			name: "fail - max length name",
			args: args{
				name:     stringPtr("test1test1test1test1test1test123"),
				email:    "test@example.com",
				password: "password",
			},
			wantErr: true,
		},
		{
			name: "fail - empty email",
			args: args{
				name:     stringPtr("test"),
				email:    "",
				password: "password",
			},
			wantErr: true,
		},
		{
			name: "fail - empty password",
			args: args{
				name:     stringPtr("test"),
				email:    "test@gmail.com",
				password: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.name, tt.args.email, tt.args.password)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				var username *string
				if name, ok := got.Name.Get(); ok {
					str := name.Value()
					username = &str
				}

				require.Equal(t, tt.args.name, username)
				require.Equal(t, tt.args.email, got.Email.Value())
				require.True(t, got.Password.Compare(tt.args.password))
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
