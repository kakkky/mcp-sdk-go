package utilities

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUriTemplate_Match(t *testing.T) {
	tests := []struct {
		name          string
		template      string
		uri           string
		expected      map[string]any
		expectedError bool
	}{
		{
			name:     "normal : basic path parameter",
			template: "/users/{userId}",
			uri:      "/users/123",
			expected: map[string]any{"userId": "123"},
		},
		{
			name:     "normal : multiple path parameters",
			template: "/users/{userId}/posts/{postId}",
			uri:      "/users/123/posts/456",
			expected: map[string]any{"userId": "123", "postId": "456"},
		},
		{
			name:     "semi normal : non-matching URI",
			template: "/users/{userId}",
			uri:      "/products/123",
		},
		{
			name:     "semi normal : mismatched parameter count",
			template: "/users/{userId}/posts/{postId}",
			uri:      "/users/123/posts",
		},
		{
			name:     "normal : exploded variable",
			template: "/tags/{tags*}",
			uri:      "/tags/golang,testing,uri",
			expected: map[string]any{"tags": []string{"golang", "testing", "uri"}},
		},
		{
			name:     "normal : variable in middle of path",
			template: "/users/{userId}/profile",
			uri:      "/users/123/profile",
			expected: map[string]any{"userId": "123"},
		},
		{
			name:     "normal : match without variables",
			template: "/about",
			uri:      "/about",
			expected: map[string]any{},
		},
		{
			name:     "normal : variable with special characters",
			template: "/search/{query}",
			uri:      "/search/testing+uri_templates",
			expected: map[string]any{"query": "testing+uri_templates"},
		},
		{
			name:     "normal : query parameter style variables",
			template: "/api{?page,limit}",
			uri:      "/api?page=2&limit=10",
			expected: map[string]any{"page": "2", "limit": "10"},
		},
		{
			name:          "semi normal : invalid template",
			template:      "/users/{userId/posts", // missing closing brace
			uri:           "/users/123/posts",
			expectedError: true,
		},
		{
			name:     "normal : fragment",
			template: "/page{#fragment}",
			uri:      "/page#section1",
			expected: map[string]any{"fragment": "section1"},
		},
		{
			name:     "normal : query parameter",
			template: "/search{?q}",
			uri:      "/search?q=test",
			expected: map[string]any{"q": "test"},
		},
		{
			name:     "normal : path parameter with query combination",
			template: "/users/{userId}{?fields}",
			uri:      "/users/123?fields=name,email",
			expected: map[string]any{
				"userId": "123",
				"fields": "name,email",
			},
		},
		{
			name:     "normal : multiple query parameters",
			template: "/api{?page,limit,sort}",
			uri:      "/api?page=2&limit=10&sort=name",
			expected: map[string]any{
				"page":  "2",
				"limit": "10",
				"sort":  "name",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut, err := NewUriTemplate(tt.template)
			if err != nil && !tt.expectedError {
				t.Fatalf("Failed to parse template: %v", err)
			}
			if err != nil && tt.expectedError {
				return
			}

			got, err := sut.Match(tt.uri)

			if (err != nil) != tt.expectedError {
				t.Errorf("Error expected: %v, got: %v, error: %v", tt.expectedError, err != nil, err)
				return
			}

			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("Mismatch in matched variables (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestUriTemplate_Expand(t *testing.T) {
	tests := []struct {
		name          string
		template      string
		values        map[string]any
		expected      string
		expectedError bool
	}{
		{
			name:     "normal : basic path parameter",
			template: "/users/{userId}",
			values:   map[string]any{"userId": "123"},
			expected: "/users/123",
		},
		{
			name:     "normal : multiple path parameters",
			template: "/users/{userId}/posts/{postId}",
			values:   map[string]any{"userId": "123", "postId": "456"},
			expected: "/users/123/posts/456",
		},
		{
			name:     "normal : parameter with special characters",
			template: "/search/{query}",
			values:   map[string]any{"query": "test+space"},
			expected: "/search/test%2Bspace",
		},
		{
			name:     "normal : exploded variable with array",
			template: "/tags/{tags*}",
			values:   map[string]any{"tags": []any{"golang", "testing", "uri"}},
			expected: "/tags/golang,testing,uri",
		},
		{
			name:     "normal : query parameter",
			template: "/search{?q}",
			values:   map[string]any{"q": "test"},
			expected: "/search?q=test",
		},
		{
			name:     "normal : multiple query parameters",
			template: "/api{?page,limit,sort}",
			values:   map[string]any{"page": "2", "limit": "10", "sort": "name"},
			expected: "/api?page=2&limit=10&sort=name",
		},
		{
			name:     "normal : partial query parameters",
			template: "/api{?page,limit,sort}",
			values:   map[string]any{"page": "2", "sort": "name"},
			expected: "/api?page=2&sort=name",
		},
		{
			name:     "normal : fragment",
			template: "/page{#fragment}",
			values:   map[string]any{"fragment": "section1"},
			expected: "/page#section1",
		},
		{
			name:     "normal : path parameter with query",
			template: "/users/{userId}{?fields}",
			values:   map[string]any{"userId": "123", "fields": "name,email"},
			expected: "/users/123?fields=name%2Cemail",
		},
		{
			name:     "normal : optional parameter not provided",
			template: "/users{/userId}",
			values:   map[string]any{},
			expected: "/users",
		},
		{
			name:     "normal : map expansion",
			template: "/search{?params*}",
			values:   map[string]any{"params": map[string]any{"q": "test", "lang": "go"}},
			expected: "/search?q=test&lang=go",
		},
		{
			name:     "normal : numeric values",
			template: "/items/{itemId}",
			values:   map[string]any{"itemId": 42},
			expected: "/items/42",
		},
		{
			name:          "semi normal : invalid template",
			template:      "/users/{userId/posts", // missing closing brace
			values:        map[string]any{"userId": "123"},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut, err := NewUriTemplate(tt.template)
			if err != nil && !tt.expectedError {
				t.Fatalf("Failed to parse template: %v", err)
			}
			if err != nil && tt.expectedError {
				return
			}

			got, err := sut.Expand(tt.values)

			if (err != nil) != tt.expectedError {
				t.Errorf("Error expected: %v, got: %v, error: %v", tt.expectedError, err != nil, err)
				return
			}

			if err == nil && got != tt.expected {
				t.Errorf("Expected: %q, got: %q", tt.expected, got)
			}
		})
	}
}
