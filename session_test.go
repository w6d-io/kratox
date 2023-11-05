/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 05/11/2023
*/

package kratox

import (
	"context"
	client "github.com/ory/kratos-client-go"
	"k8s.io/utils/pointer"
	"reflect"
	"testing"
	"time"
)

func TestGetAddressFromCtx(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "with no context",
			args:    args{ctx: context.Background()},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAddressFromCtx(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddressFromCtx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAddressFromCtx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSessionFromCtx(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *client.Session
		wantErr bool
	}{
		{
			name:    "raise an error on empty context",
			args:    args{ctx: context.Background()},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSessionFromCtx(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSessionFromCtx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSessionFromCtx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAddressInCtx(t *testing.T) {
	type args struct {
		ctx     context.Context
		address string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "set an address into the context",
			args: args{ctx: context.Background(), address: "http://localhost"},
			want: context.WithValue(context.Background(), AddressKey, "http://localhost"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetAddressInCtx(tt.args.ctx, tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetAddressInCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetSessionInCtx(t *testing.T) {
	type args struct {
		ctx     context.Context
		session *client.Session
	}
	var now = time.Now()
	var session = &client.Session{
		Active:                      pointer.Bool(true),
		AuthenticatedAt:             &now,
		AuthenticationMethods:       nil,
		AuthenticatorAssuranceLevel: nil,
		Devices:                     nil,
		ExpiresAt:                   nil,
		Id:                          "none",
		Identity:                    client.Identity{},
		IssuedAt:                    nil,
		AdditionalProperties:        nil,
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "set session into context",
			args: args{
				ctx:     context.Background(),
				session: session,
			},
			want: context.WithValue(context.Background(), SessionKey, session),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetSessionInCtx(tt.args.ctx, tt.args.session); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetSessionInCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}
