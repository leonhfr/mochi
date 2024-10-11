package mochi

import (
	"context"
	"net/http"
	"testing"
)

func Test_GetTemplate(t *testing.T) {
	tests := []struct {
		name string
		test getItemTestCase
	}{
		{
			name: "should get a template",
			test: getItemTestCase{
				status: http.StatusOK,
				id:     "TEMPLATE_ID",
				res:    Template{ID: "TEMPLATE_ID", Name: "TemplateName", Content: "Template content"},
				want:   Template{ID: "TEMPLATE_ID", Name: "TemplateName", Content: "Template content"},
				err:    "",
			},
		},
		{
			name: "should return an error",
			test: getItemTestCase{
				status: http.StatusBadRequest,
				id:     "TEMPLATE_ID",
				res:    `{"errors":["ERROR_MESSAGE"]}`,
				want:   Template{},
				err:    "mochi: ERROR_MESSAGE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testGetItem("/api/templates", tt.test, func(client *Client, id string) (any, error) {
			return client.GetTemplate(context.Background(), id)
		}))
	}
}

func Test_ListTemplates(t *testing.T) {
	tests := []struct {
		name string
		test listItemTestCase[Template]
	}{
		{
			name: "should call the endpoint once",
			test: listItemTestCase[Template]{
				responses: []listItemTestCaseResponse{
					{
						status: http.StatusOK,
						params: map[string]string{"limit": "100"},
						res: listResponse[Template]{
							Docs: []Template{{ID: "TEMPLATE_ID", Name: "TemplateName", Content: "Template content"}},
						},
					},
				},
				want: []Template{
					{ID: "TEMPLATE_ID", Name: "TemplateName", Content: "Template content"},
				},
			},
		},
		{
			name: "should call the endpoint several times",
			test: listItemTestCase[Template]{
				responses: []listItemTestCaseResponse{
					{
						status: http.StatusOK,
						params: map[string]string{"limit": "100"},
						res: listResponse[Template]{
							Docs:     []Template{{ID: "TEMPLATE_ID_1", Name: "TemplateName1", Content: "Template content"}},
							Bookmark: "BOOKMARK_1",
						},
					},
					{
						status: http.StatusOK,
						params: map[string]string{"limit": "100", "bookmark": "BOOKMARK_1"},
						res: listResponse[Template]{
							Docs: []Template{{ID: "TEMPLATE_ID_2", Name: "TemplateName2", Content: "Template content"}},
						},
					},
				},
				want: []Template{
					{ID: "TEMPLATE_ID_1", Name: "TemplateName1", Content: "Template content"},
					{ID: "TEMPLATE_ID_2", Name: "TemplateName2", Content: "Template content"},
				},
			},
		},
		{
			name: "should return an error",
			test: listItemTestCase[Template]{
				responses: []listItemTestCaseResponse{
					{
						status: http.StatusBadRequest,
						params: map[string]string{"limit": "100"},
						res:    `{"errors":["ERROR_MESSAGE"]}`,
					},
				},
				err: "mochi: ERROR_MESSAGE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testListItem("/api/templates", tt.test, func(client *Client, _ string) ([]Template, error) {
			return client.ListTemplates(context.Background())
		}))
	}
}
