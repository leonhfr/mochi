package mochi

import (
	"encoding/json"
	"maps"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_errorResponse_error(t *testing.T) {
	tests := []struct {
		name string
		er   *errorResponse
		want string
	}{
		{
			name: "should return nil (nil)",
			er:   nil,
			want: "",
		},
		{
			name: "should return nil (empty struct)",
			er:   &errorResponse{},
			want: "",
		},
		{
			name: "should return a server error",
			er:   &errorResponse{errors: []string{"ERROR_MESSAGE_1", "ERROR_MESSAGE_2"}},
			want: "mochi: ERROR_MESSAGE_1 ERROR_MESSAGE_2",
		},
		{
			name: "should return a validation error",
			er:   &errorResponse{validation: map[string]string{"FIELD_1": "ERROR_MESSAGE_1", "FIELD_2": "ERROR_MESSAGE_2"}},
			want: "mochi(validation): FIELD_1: ERROR_MESSAGE_1 FIELD_2: ERROR_MESSAGE_2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.er.error()
			if tt.want != "" {
				assert.EqualError(t, got, tt.want)
			} else {
				assert.NoError(t, got)
			}
		})
	}
}

func Test_errorResponse_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  errorResponse
		err   string
	}{
		{
			name:  "should parse server errors",
			input: `{"errors":["ERROR_MESSAGE"]}`,
			want:  errorResponse{errors: []string{"ERROR_MESSAGE"}},
			err:   "",
		},
		{
			name:  "should parse validation errors",
			input: `{"errors":{"FIELD":"ERROR_MESSAGE"}}`,
			want:  errorResponse{validation: map[string]string{"FIELD": "ERROR_MESSAGE"}},
			err:   "",
		},
		{
			name:  "should return an error",
			input: "",
			want:  errorResponse{},
			err:   "unexpected end of JSON input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var er errorResponse
			err := json.Unmarshal([]byte(tt.input), &er)

			assert.Equal(t, tt.want, er)
			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type createItemTestCase[CreateItem any] struct {
	status int
	req    CreateItem
	res    any
	want   any
	err    string
}

func testCreateItem[CreateItem any](path string, test createItemTestCase[CreateItem], method func(*Client, CreateItem) (any, error)) func(t *testing.T) {
	return func(t *testing.T) {
		defer gock.Off()

		token := "TOKEN"
		gock.New(baseURL).
			BasicAuth(token, "").
			Post(path).
			MatchType("json").
			JSON(test.req).
			Reply(test.status).
			JSON(test.res)

		got, err := method(New(token), test.req)

		assert.Equal(t, test.want, got)
		if test.err != "" {
			assert.EqualError(t, err, test.err)
		} else {
			assert.NoError(t, err)
		}
		require.True(t, gock.IsDone())
	}
}

type getItemTestCase struct {
	status int
	id     string
	res    any
	want   any
	err    string
}

func testGetItem(path string, test getItemTestCase, method func(*Client, string) (any, error)) func(t *testing.T) {
	return func(t *testing.T) {
		defer gock.Off()

		token := "TOKEN"
		gock.New(baseURL).
			BasicAuth(token, "").
			Get(path).
			Path(test.id).
			Reply(test.status).
			JSON(test.res)

		got, err := method(New(token), test.id)

		assert.Equal(t, test.want, got)
		if test.err != "" {
			assert.EqualError(t, err, test.err)
		} else {
			assert.NoError(t, err)
		}
		require.True(t, gock.IsDone())
	}
}

type listItemTestCaseResponse struct {
	status int
	params map[string]string
	res    any
	want   any
	err    error
}

type listItemTestCase struct {
	id        string
	params    map[string]string
	responses []listItemTestCaseResponse
	total     int
	err       string
}

func testListItem[Item any](path string, test listItemTestCase, method func(*Client, string, func([]Item) error) error) func(t *testing.T) {
	return func(t *testing.T) {
		defer gock.Off()

		token := "TOKEN"
		for _, res := range test.responses {
			params := make(map[string]string)
			maps.Copy(params, res.params)
			maps.Copy(params, test.params)

			gock.New(baseURL).
				BasicAuth(token, "").
				Get(path).
				MatchParams(params).
				Reply(res.status).
				JSON(res.res)
		}

		var index int
		err := method(New(token), test.id, func(items []Item) error {
			want := test.responses[index].want
			err := test.responses[index].err
			assert.Equal(t, want, items)
			index++
			return err
		})

		assert.Equal(t, test.total, index)
		if test.err != "" {
			assert.EqualError(t, err, test.err)
		} else {
			assert.NoError(t, err)
		}
		require.True(t, gock.IsDone())
	}
}

type updateItemTestCase[UpdateItem any] struct {
	status int
	id     string
	req    UpdateItem
	res    any
	want   any
	err    string
}

func testUpdateItem[UpdateItem any](path string, test updateItemTestCase[UpdateItem], method func(*Client, UpdateItem) (any, error)) func(t *testing.T) {
	return func(t *testing.T) {
		defer gock.Off()

		token := "TOKEN"
		gock.New(baseURL).
			BasicAuth(token, "").
			Path(test.id).
			Post(path).
			MatchType("json").
			JSON(test.req).
			Reply(test.status).
			JSON(test.res)

		got, err := method(New(token), test.req)

		assert.Equal(t, test.want, got)
		if test.err != "" {
			assert.EqualError(t, err, test.err)
		} else {
			assert.NoError(t, err)
		}
		require.True(t, gock.IsDone())
	}
}

type deleteItemTestCase struct {
	status int
	id     string
	res    any
	err    string
}

func testDeleteItem(path string, test deleteItemTestCase, method func(*Client, string) error) func(t *testing.T) {
	return func(t *testing.T) {
		defer gock.Off()

		token := "TOKEN"
		r := gock.New(baseURL).
			BasicAuth(token, "").
			Delete(path).
			Path(test.id).
			Reply(test.status)

		if test.res != nil {
			r.JSON(test.res)
		}

		err := method(New(token), test.id)

		if test.err != "" {
			assert.EqualError(t, err, test.err)
		} else {
			assert.NoError(t, err)
		}
		require.True(t, gock.IsDone())
	}
}
