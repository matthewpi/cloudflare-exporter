//
// Copyright (c) 2021 Matthew Penner
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package cloudflare

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

// Auth .
type Auth interface {
	// Authorize .
	Authorize(context.Context, http.Header) error
}

// KeyAuthorization .
type KeyAuthorization struct {
	Email string
	Key   string
}

var _ Auth = (*KeyAuthorization)(nil)

// NewKeyAuthorization .
func NewKeyAuthorization(email, key string) (*KeyAuthorization, error) {
	if email == "" {
		return nil, errors.New("cloudflare: missing email")
	}
	if key == "" {
		return nil, errors.New("cloudflare: missing key")
	}
	return &KeyAuthorization{
		Email: email,
		Key:   key,
	}, nil
}

// Authorize .
func (a *KeyAuthorization) Authorize(_ context.Context, h http.Header) error {
	h.Set("X-Auth-Email", a.Email)
	h.Set("X-Auth-Key", a.Key)
	return nil
}

// TokenAuthorization .
type TokenAuthorization struct {
	Token string
}

var _ Auth = (*TokenAuthorization)(nil)

// NewTokenAuthorization .
func NewTokenAuthorization(token string) (*TokenAuthorization, error) {
	if token == "" {
		return nil, errors.New("cloudflare: missing token")
	}
	return &TokenAuthorization{
		Token: token,
	}, nil
}

// Authorize .
func (a *TokenAuthorization) Authorize(_ context.Context, h http.Header) error {
	h.Set("Authorization", "Bearer "+a.Token)
	return nil
}
