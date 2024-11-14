package gcd

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestPrefix = "https://example.org/api"

func TestSeriesReq_URL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		req       SeriesReq
		expected  string
		shouldErr bool
	}{
		{
			"name-only",
			SeriesReq{
				Name: "Batman",
			},
			"https://example.org/api/series/name/Batman/",
			false,
		},
		{
			"name-year",
			SeriesReq{
				Name: "Batman",
				Year: 2000,
			},
			"https://example.org/api/series/name/Batman/year/2000/",
			false,
		},
		{
			"name-year-issue",
			SeriesReq{
				Name:    "Batman",
				IssueNo: 12,
				Year:    2000,
			},
			"https://example.org/api/series/name/Batman/issue/12/year/2000/",
			false,
		},
		{
			"id",
			SeriesReq{
				ID: 7096,
			},
			"https://example.org/api/series/7096/",
			false,
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp, err := tt.req.URL(TestPrefix)
			if tt.shouldErr {
				require.Error(t, err, "expected error")

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, resp)
		})
	}
}

func TestAPI_Series_SessionID(t *testing.T) {
	t.Parallel()

	var requests []*http.Request

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{}`)
	}))
	t.Cleanup(server.Close)

	api := API{
		Prefix:    "http://" + server.Listener.Addr().String() + "/api/",
		SessionID: "foobar123",
	}

	_, err := api.Series(context.Background(), SeriesReq{Name: "Batman"})
	require.NoError(t, err)
	require.Len(t, requests, 1)
	req0 := requests[0]
	assert.Equal(t, `/api/series/name/Batman/`, req0.URL.Path)
	assert.Len(t, req0.Cookies(), 1, "cookies length")
	assert.Equal(t, http.MethodGet, req0.Method, "method")

	if cookies := req0.Cookies(); len(cookies) > 0 {
		cookie := cookies[0]
		assert.Equal(t, "gcdsessionid", cookie.Name)
		assert.Equal(t, "foobar123", cookie.Value)
	}
}

const supermanSeriesList = `{
    "count": 2,
    "next": null,
    "previous": null,
	"results": [
		{
            "api_url": "https://www.comics.org/api/series/197230/",
            "name": "Adventures of Superman: Jon Kent",
            "country": "us",
            "language": "en",
            "active_issues": [
                "https://www.comics.org/api/issue/2500292/",
                "https://www.comics.org/api/issue/2513856/"
            ],
            "issue_descriptors": [
                "1 [Clayton Henry Cover]",
                "2 [Clayton Henry Cover]"
			],
            "color": "color",
            "dimensions": "standard Modern Age US",
            "paper_stock": "glossy cover; matte paper interiors",
            "binding": "saddle-stitched",
            "publishing_format": "limited series",
            "notes": "",
            "year_began": 2023,
            "year_ended": 2023,
            "publisher": "https://www.comics.org/api/publisher/54/"
        },
		{
            "api_url": "https://www.comics.org/api/series/196803/",
            "name": "Superman",
            "country": "us",
            "language": "en",
            "active_issues": [
                "https://www.comics.org/api/issue/2495111/",
                "https://www.comics.org/api/issue/2507431/"
            ],
            "issue_descriptors": [
                "1 [Jamal Campbell Cover]",
                "2 [Jamal Campbell Cover]"
            ],
            "color": "color",
            "dimensions": "standard Modern Age US",
            "paper_stock": "glossy cover; matte paper interiors (#1-#9), glossy interiors (starting with #10)",
            "binding": "saddle-stitched",
            "publishing_format": "ongoing series",
            "notes": "",
            "year_began": 2023,
            "year_ended": null,
            "publisher": "https://www.comics.org/api/publisher/54/"
		}
	]
}`

const supermanSeriesInstance = `{
	"api_url": "https://www.comics.org/api/series/196803/",
	"name": "Superman",
	"country": "us",
	"language": "en",
	"active_issues": [
		"https://www.comics.org/api/issue/2495111/",
		"https://www.comics.org/api/issue/2507431/"
	],
	"issue_descriptors": [
		"1 [Jamal Campbell Cover]",
		"2 [Jamal Campbell Cover]"
	],
	"color": "color",
	"dimensions": "standard Modern Age US",
	"paper_stock": "glossy cover; matte paper interiors (#1-#9), glossy interiors (starting with #10)",
	"binding": "saddle-stitched",
	"publishing_format": "ongoing series",
	"notes": "",
	"year_began": 2023,
	"year_ended": null,
	"publisher": "https://www.comics.org/api/publisher/54/"
}`

func TestAPI_Series(t *testing.T) {
	t.Parallel()

	var requests []*http.Request

	responses := map[string]string{
		"/api/series/name/Superman/year/2023/": supermanSeriesList,
		"/api/series/196803/":                  supermanSeriesInstance,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r)

		respData, ok := responses[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, respData)
	}))
	t.Cleanup(server.Close)

	api := API{
		Prefix: "http://" + server.Listener.Addr().String() + "/api/",
	}

	t.Run("series", func(t *testing.T) {
		t.Parallel()

		resp, err := api.Series(context.Background(), SeriesReq{
			Name: "Superman",
			Year: 2023,
		})
		require.NoError(t, err, "api.Series")
		assert.Equal(t, 2, resp.Count, "resp.Count")
		assert.Len(t, resp.Results, resp.Count, "resp.Results")
	})

	t.Run("instance", func(t *testing.T) {
		t.Parallel()

		resp, err := api.SeriesInstance(context.Background(), 196803)
		require.NoError(t, err, "api.SeriesInstance")

		assert.Len(t, resp.ActiveIssues, 2, "resp.ActiveIssues")
		assert.Equal(t, "Superman", resp.Name, "resp.Name")
	})
}
