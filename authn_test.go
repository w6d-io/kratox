/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 14/12/2023
*/

package kratox_test

import (
	"context"
	"github.com/google/uuid"
	client "github.com/ory/kratos-client-go"
	"github.com/w6d-io/kratox"
	"k8s.io/utils/pointer"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var now = time.Now()
var expire = now.Add(time.Hour)
var id, _ = uuid.NewUUID()
var session = &client.Session{
	Active:          pointer.Bool(true),
	AuthenticatedAt: &now,
	ExpiresAt:       &expire,
	Id:              id.String(),
	Identity: client.Identity{
		Id: id.String(),
	},
	IssuedAt: &now,
}

func TestAuthRequestFunc(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *http.Request
	}

	req, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)
	reqWithCookie, _ := http.NewRequest(http.MethodGet, "http://localhost", nil)
	reqWithCookie.AddCookie(&http.Cookie{
		Name:  kratox.CookieName,
		Value: "test",
	})
	tests := []struct {
		name string
		args args
		mock kratosMock
		want context.Context
	}{
		{
			name: "no kratos cookie",
			args: args{r: req},
			mock: kratosMock{},
			want: nil,
		},
		{
			name: "fail to get session from http",
			args: args{r: reqWithCookie},
			mock: kratosMock{
				behaviour: "ko",
			},
			want: context.WithValue(context.Background(), kratox.CookieKey, "test"),
		},
		{
			name: "get nil session",
			args: args{r: reqWithCookie},
			mock: kratosMock{
				behaviour: "sessionNil",
			},
			want: context.WithValue(context.Background(), kratox.CookieKey, "test"),
		},
		{
			name: "with kratos cookie",
			args: args{r: reqWithCookie},
			mock: kratosMock{
				behaviour: "ok",
			},
			want: context.WithValue(context.WithValue(context.Background(), kratox.CookieKey, "test"), kratox.SessionKey, session),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kratox.Kratox = &tt.mock
			if got := kratox.AuthRequestFunc(tt.args.ctx, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthRequestFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
