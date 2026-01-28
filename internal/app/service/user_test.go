package service

import (
	"context"
	"github.com/go-list-templ/grpc/internal/core/domain/entity"
	"github.com/go-list-templ/grpc/internal/core/port"
	"reflect"
	"testing"
)

func TestNewUser(t *testing.T) {
	type args struct {
		u port.UserRepo
		o port.OutboxRepo
		t port.TransactionManager
	}
	tests := []struct {
		name string
		args args
		want *User
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUser(tt.args.u, tt.args.o, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_All(t *testing.T) {
	type fields struct {
		userRepo   port.UserRepo
		outboxRepo port.OutboxRepo
		trm        port.TransactionManager
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entity.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &User{
				userRepo:   tt.fields.userRepo,
				outboxRepo: tt.fields.outboxRepo,
				trm:        tt.fields.trm,
			}
			got, err := s.All(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("All() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_Create(t *testing.T) {
	type fields struct {
		userRepo   port.UserRepo
		outboxRepo port.OutboxRepo
		trm        port.TransactionManager
	}
	type args struct {
		ctx  context.Context
		user entity.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &User{
				userRepo:   tt.fields.userRepo,
				outboxRepo: tt.fields.outboxRepo,
				trm:        tt.fields.trm,
			}
			got, err := s.Create(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}
