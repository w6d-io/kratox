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
	"net/http"

	client "github.com/ory/kratos-client-go"
	"github.com/pkg/errors"

	"github.com/w6d-io/kratox"
)

type kratosMock struct {
	kratox.Helper
	behaviour string
	token     bool
	subject   string
	code      int
	provider  string
}

func (k kratosMock) GetSessionFromHTTP(_ context.Context, _ *http.Request) (*client.Session, error) {
	switch k.behaviour {
	case "ko":
		return nil, errors.New("failed to connect")
	case "sessionNil":
		return nil, nil
	default:
		return session, nil
	}
}
