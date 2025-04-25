package main_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPractice(t *testing.T) {
	testCases := []struct {
		name          string
		requestPath   string
		requestBody   string
		requestMethod string
		wantCode      int
		wantBody      string
	}{
		{
			name:          "Create item #1",
			requestPath:   "/items",
			requestMethod: http.MethodPost,
			requestBody:   `{"name": "item1", "price": 100}`,
			wantCode:      http.StatusOK,
			wantBody:      "OK",
		},
		{
			name:          "Create item #2",
			requestPath:   "/items",
			requestMethod: http.MethodPost,
			requestBody:   `{"name": "item2", "price": 200}`,
			wantCode:      http.StatusOK,
			wantBody:      "OK",
		},
		{
			name:          "View items",
			requestPath:   "/items/view",
			requestMethod: http.MethodGet,
			requestBody:   `{"name": "item2", "price": 200}`,
			wantCode:      http.StatusOK,
			wantBody: `
<!DOCTYPE html>
<html>
<body>

<h1>Список Товаров</h1>

<div>
    <h2>item1</h2>
    <p>Цена $100</p>
</div>
<div>
    <h2>item2</h2>
    <p>Цена $200</p>
</div>

</body>
</html>
`,
		},
		{
			name:          "Create item #3",
			requestPath:   "/items",
			requestMethod: http.MethodPost,
			requestBody:   `{"name": "item3", "price": 300}`,
			wantCode:      http.StatusOK,
			wantBody:      "OK",
		},
		{
			name:          "View items #2",
			requestPath:   "/items/view",
			requestMethod: http.MethodGet,
			requestBody:   `{"name": "item3", "price": 300}`,
			wantCode:      http.StatusOK,
			wantBody: `
<!DOCTYPE html>
<html>
<body>

<h1>Список Товаров</h1>

<div>
    <h2>item1</h2>
    <p>Цена $100</p>
</div>
<div>
    <h2>item2</h2>
    <p>Цена $200</p>
</div>
<div>
    <h2>item3</h2>
    <p>Цена $300</p>
</div>

</body>
</html>
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, tErr := http.NewRequest(
				tc.requestMethod,
				"http://localhost:8080"+tc.requestPath,
				strings.NewReader(tc.requestBody),
			)
			tr := require.New(t)
			tr.NoError(tErr)

			req.Header.Set("Content-Type", "application/json")

			httpClient := http.Client{}
			resp, tErr := httpClient.Do(req)
			tr.NoError(tErr)
			defer resp.Body.Close()

			bodyBytes, tErr := io.ReadAll(resp.Body)
			tr.NoError(tErr)

			tr.Equal(tc.wantCode, resp.StatusCode)
			tr.Equal(strings.ReplaceAll(tc.wantBody, "\n", ""), strings.ReplaceAll(string(bodyBytes), "\n", ""))
		})
	}
}
