package gcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type SeriesInstance struct {
	APIURL           string   `json:"api_url"`
	Name             string   `json:"name"`
	Country          string   `json:"country"`
	Language         string   `json:"language"`
	ActiveIssues     []string `json:"active_issues"`
	IssueDescriptors []string `json:"issue_descriptors"`
	Color            string   `json:"color"`
	Dimensions       string   `json:"dimensions"`
	PaperStock       string   `json:"paper_stock"`
	Binding          string   `json:"binding"`
	PublishingFormat string   `json:"publishing_format"`
	Notes            string   `json:"notes"`
	YearBegan        int      `json:"year_began"`
	YearEnded        int      `json:"year_ended"`
	Publisher        string   `json:"publisher"`
}

type SeriesReq struct {
	ID      int // series ID should not be provided together with the series Name
	IssueNo int // this is *not* the issue ID
	Name    string
	Year    int

	Format string // optional: "api" or "json"
	Page   int
}

func (r SeriesReq) URL(prefix string) (string, error) {
	if r.ID < 0 {
		return "", errors.New("if ID is provided, it needs to be greater than zero")
	}

	if r.ID > 0 && r.Name != "" {
		return "", errors.New("cannot specify both ID and Name")
	}

	url := prefix
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}

	url += "/series"

	if r.ID > 0 {
		url += "/" + strconv.Itoa(r.ID)
	}

	if r.Name != "" {
		url += "/name/" + r.Name
	}

	if r.IssueNo > 0 {
		url += "/issue/" + strconv.Itoa(r.IssueNo)
	}

	if r.Year > 0 {
		url += "/year/" + strconv.Itoa(r.Year)
	}

	url += "/"

	var params []string

	if r.Format != "" {
		params = append(params, "format="+r.Format)
	}

	if r.Page > 0 {
		params = append(params, "page="+strconv.Itoa(r.Page))
	}

	if len(params) > 0 {
		url += "?" + strings.Join(params, "&")
	}

	return url, nil
}

type SeriesResp struct {
	Count    int              `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous,omitempty"`
	Results  []SeriesInstance `json:"results"`
}

func (a API) SeriesFromURL(ctx context.Context, url string) (SeriesResp, error) {
	var seriesResp SeriesResp

	resp, err := a.req(ctx, url)
	if err != nil {
		return seriesResp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return seriesResp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return seriesResp, fmt.Errorf("io.ReadAll: %w", err)
	}

	if err := json.Unmarshal(data, &seriesResp); err != nil {
		return seriesResp, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return seriesResp, nil
}

func (a API) Series(ctx context.Context, req SeriesReq) (SeriesResp, error) {
	prefix := a.Prefix
	if prefix == "" {
		prefix = DefaultPrefix
	}

	uu, err := req.URL(prefix)
	if err != nil {
		return SeriesResp{}, fmt.Errorf("failed to construct URL: %w", err)
	}

	return a.SeriesFromURL(ctx, uu)
}

func (a API) SeriesInstanceFromURL(ctx context.Context, url string) (SeriesInstance, error) {
	var seriesInstance SeriesInstance

	resp, err := a.req(ctx, url)
	if err != nil {
		return seriesInstance, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return seriesInstance, fmt.Errorf("io.ReadAll: %w", err)
	}

	if err := json.Unmarshal(data, &seriesInstance); err != nil {
		return seriesInstance, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return seriesInstance, nil
}

func (a API) SeriesInstance(ctx context.Context, id int) (SeriesInstance, error) {
	uu := a.Prefix
	if uu[len(uu)-1] == '/' {
		uu = uu[:len(uu)-1]
	}

	uu += "/series/" + strconv.Itoa(id)
	uu += "/"

	return a.SeriesInstanceFromURL(ctx, uu)
}
