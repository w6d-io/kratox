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
	"net/http"
	"reflect"
	"testing"
	"time"

	client "github.com/ory/kratos-client-go"
	"k8s.io/utils/pointer"
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
		{
			name:    "with a no string address",
			args:    args{ctx: context.WithValue(context.Background(), AddressKey, 1)},
			want:    "",
			wantErr: true,
		},
		{
			name:    "get address",
			args:    args{ctx: context.WithValue(context.Background(), AddressKey, "http://localhost")},
			want:    "http://localhost",
			wantErr: false,
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
			args: args{ctx: nil, address: "http://localhost"},
			want: context.WithValue(context.Background(), AddressKey, "http://localhost"),
		},
		{
			name: "get an empty context",
			args: args{ctx: nil, address: ""},
			want: context.Background(),
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

func TestGetCookieFromCtx(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil context",
			args: args{ctx: nil},
			want: "",
		},
		{
			name: "no cookie from context",
			args: args{ctx: context.Background()},
			want: "",
		},
		{
			name: "get cookie from context",
			args: args{ctx: context.WithValue(context.Background(), CookieKey, "test")},
			want: "test",
		},
		{
			name: "get empty string due to wrong type",
			args: args{ctx: context.WithValue(context.Background(), CookieKey, 1)},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCookieFromCtx(tt.args.ctx); got != tt.want {
				t.Errorf("GetCookieFromCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetCookieInCtx(t *testing.T) {
	type args struct {
		ctx    context.Context
		cookie string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "get empty context",
			args: args{ctx: nil},
			want: context.Background(),
		},
		{
			name: "get context with cook",
			args: args{ctx: nil, cookie: "test"},
			want: context.WithValue(context.Background(), CookieKey, "test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetCookieInCtx(tt.args.ctx, tt.args.cookie); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetCookieInCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetCookieFromHttpToCtx(t *testing.T) {
	type args struct {
		ctx context.Context
		req *http.Request
	}
	r1, err := http.NewRequest("GET", "http//localhost", nil)
	if err != nil {
		t.Errorf("failed to create http request")
		return
	}
	r2, err := http.NewRequest("GET", "http//localhost", nil)
	if err != nil {
		t.Errorf("failed to create http request")
		return
	}

	r2.AddCookie(&http.Cookie{
		Name:  CookieName,
		Value: "12345",
	})
	r3, err := http.NewRequest("GET", "http//localhost", nil)
	if err != nil {
		t.Errorf("failed to create http request")
		return
	}

	r3.AddCookie(&http.Cookie{
		Name:  CookieName,
		Value: "",
	})
	tests := []struct {
		name    string
		args    args
		want    context.Context
		wantErr bool
	}{
		{
			name:    "with no context and no cookie",
			args:    args{ctx: nil, req: r1},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "with no context but with cookie",
			args:    args{ctx: nil, req: r2},
			want:    context.WithValue(context.Background(), CookieKey, "12345"),
			wantErr: false,
		},
		{
			name:    "with no context but with empty cookie",
			args:    args{ctx: nil, req: r3},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetCookieFromHttpToCtx(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetCookieFromHttpToCtx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetCookieFromHttpToCtx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSession(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "not authorized",
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
		{
			name: "authorized",
			args: args{ctx: SetSessionInCtx(context.Background(), &client.Session{
				Active: pointer.Bool(true),
			})},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetSession(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("GetSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
