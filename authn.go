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

package kratox

import (
	"context"
	"github.com/w6d-io/x/logx"
	"net/http"
)

func AuthRequestFunc(ctx context.Context, r *http.Request) context.Context {
	ctx2, err := SetCookieFromHttpToCtx(ctx, r)
	if err != nil {
		logx.WithName(ctx, "OptionAuthn").Info("get kratos cookie from http request failed")
		return ctx
	}
	ctx = ctx2
	session, err := Kratox.GetSessionFromHTTP(ctx, r)
	if err != nil {
		logx.WithName(ctx, "OptionAuthn").Info("get session from kratos failed")
		return ctx
	}
	if session == nil {
		logx.WithName(ctx, "OptionAuthn").Info("get session from kratos failed")
		return ctx
	}
	return SetSessionInCtx(ctx, session)
}
