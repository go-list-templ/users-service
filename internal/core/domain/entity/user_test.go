package entity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	type args struct {
		name  string
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success - create user",
			args: args{
				name:  "test",
				email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "success - russian domain email",
			args: args{
				name:  "test",
				email: "пользователь@компания.рф",
			},
			wantErr: false,
		},
		{
			name: "success - email with one domain",
			args: args{
				name:  "test",
				email: "invalid@email",
			},
			wantErr: false,
		},
		{
			name: "fail - min length name",
			args: args{
				name:  "t",
				email: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "fail - max length name",
			args: args{
				name:  "test1test1test1test1test1test123",
				email: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "fail - empty name",
			args: args{
				name:  "",
				email: "example@example.com",
			},
			wantErr: true,
		},
		{
			name: "fail - empty email",
			args: args{
				name:  "test",
				email: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.name, tt.args.email)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				require.Equal(t, tt.args.name, got.Name.Value())
				require.Equal(t, tt.args.email, got.Email.Value())
			}
		})
	}
}
