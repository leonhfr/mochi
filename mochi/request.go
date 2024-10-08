package mochi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/carlmjohnson/requests"
)

func buildRequest(c *Client) *requests.Builder {
	return requests.
		URL(c.baseURL).
		Client(c.client).
		Transport(c.transport).
		Accept("application/json").
		BasicAuth(c.token, "")
}

func createItem[Item any](ctx context.Context, c *Client, path string, req any) (Item, error) {
	var item Item
	rb := buildRequest(c).
		Path(path).
		Method(http.MethodPost).
		BodyJSON(req).
		ToJSON(&item)
	err := executeRequest(ctx, rb)
	return item, err
}

func getItem[Item any](ctx context.Context, c *Client, path, id string) (Item, error) {
	var item Item
	rb := buildRequest(c).
		Pathf("%s/%s", path, id).
		Method(http.MethodGet).
		ToJSON(&item)
	err := executeRequest(ctx, rb)
	return item, err
}

type listResponse[Item any] struct {
	Bookmark string `json:"bookmark"`
	Docs     []Item `json:"docs"`
}

func listItems[Item any](ctx context.Context, c *Client, path string, params url.Values, cb func([]Item) error) error {
	var bookmark string
	for {
		var res listResponse[Item]
		rb := buildRequest(c).
			Path(path).
			Method(http.MethodGet).
			Params(params).
			ParamOptional("limit", "100").
			ParamOptional("bookmark", bookmark).
			ToJSON(&res)
		if err := executeRequest(ctx, rb); err != nil {
			return err
		}
		if err := cb(res.Docs); err != nil {
			return err
		}
		bookmark = res.Bookmark
		if bookmark == "" || bookmark == "nil" || len(res.Docs) == 0 {
			break
		}
	}
	return nil
}

func updateItem[Item any](ctx context.Context, c *Client, path, id string, req any) (Item, error) {
	var item Item
	rb := buildRequest(c).
		Pathf("%s/%s", path, id).
		Method(http.MethodPost).
		BodyJSON(req).
		ToJSON(&item)
	err := executeRequest(ctx, rb)
	return item, err
}

func deleteItem(ctx context.Context, c *Client, path, id string) error {
	rb := buildRequest(c).
		Pathf("%s/%s", path, id).
		Method(http.MethodDelete)
	err := executeRequest(ctx, rb)
	return err
}

func executeRequest(ctx context.Context, rb *requests.Builder) error {
	var errRes errorResponse
	err := rb.
		AddValidator(requests.ErrorJSON(&errRes)).
		Fetch(ctx)
	switch {
	case errors.Is(err, requests.ErrInvalidHandled):
		return errRes.error()
	case err != nil:
		return err
	default:
		return nil
	}
}

type errorResponse struct {
	errors     []string
	validation map[string]string
}

func (er *errorResponse) error() error {
	if er == nil {
		return nil
	}
	if er.errors != nil {
		return fmt.Errorf("mochi: %s", strings.Join(er.errors, " "))
	}
	if er.validation != nil {
		var errors []string
		for field, error := range er.validation {
			errors = append(errors, fmt.Sprintf("%s: %s", field, error))
		}
		slices.Sort(errors)
		return fmt.Errorf("mochi(validation): %s", strings.Join(errors, " "))
	}
	return nil
}

func (er *errorResponse) UnmarshalJSON(input []byte) error {
	var parsedErrors struct {
		Errors []string `json:"errors"`
	}
	err := json.Unmarshal(input, &parsedErrors)
	if err == nil {
		er.errors = parsedErrors.Errors
		er.validation = nil
		return nil
	}

	type validationErrors struct {
		Errors map[string]string `json:"errors"`
	}
	parsedValidationErrors := validationErrors{Errors: map[string]string{}}
	err = json.Unmarshal(input, &parsedValidationErrors)
	if err == nil {
		er.errors = nil
		er.validation = parsedValidationErrors.Errors
		return nil
	}

	return err
}
