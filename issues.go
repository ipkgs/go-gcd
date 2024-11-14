package gcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type StorySet struct {
	Type           string `json:"type"`
	Title          string `json:"title"`
	Feature        string `json:"feature"`
	SequenceNumber int    `json:"sequence_number"`
	PageCount      string `json:"page_count"`
	Script         string `json:"script"`
	Pencils        string `json:"pencils"`
	Inks           string `json:"inks"`
	Colors         string `json:"colors"`
	Letters        string `json:"letters"`
	Editing        string `json:"editing"`
	JobNumber      string `json:"job_number"`
	Genre          string `json:"genre"`
	Characters     string `json:"characters"`
	Synopsis       string `json:"synopsis"`
	Notes          string `json:"notes"`
}

type IssueReq struct {
	ID int

	Format string // optional: "api" or "json"
}

func (r IssueReq) URL(prefix string) (string, error) {
	if r.ID <= 0 {
		return "", errors.New("invalid ID")
	}

	url := prefix
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}

	url += "/issue/" + strconv.Itoa(r.ID)
	url += "/"

	var params []string

	if r.Format != "" {
		params = append(params, "format="+r.Format)
	}

	if len(params) > 0 {
		url += "?" + strings.Join(params, "&")
	}

	return url, nil
}

type IssueResp struct {
	APIURL           string      `json:"api_url"`
	SeriesName       string      `json:"series_name"`
	Descriptor       string      `json:"descriptor"`
	PublicationDate  string      `json:"publication_date"`
	Price            string      `json:"price"`
	PageCount        string      `json:"page_count"`
	Editing          string      `json:"editing"`
	Brand            string      `json:"brand"`
	ISBN             string      `json:"isbn"`
	Barcode          string      `json:"barcode"`
	Rating           string      `json:"rating"`
	OnSaleDate       string      `json:"on_sale_date"`
	Notes            string      `json:"notes"`
	VariantOf        interface{} `json:"variant_of"`
	Series           string      `json:"series"`
	Cover            string      `json:"cover"`
	StorySet         []StorySet  `json:"story_set"`
	IndiciaPublisher string      `json:"indicia_publisher"`
	IndiciaFrequency string      `json:"indicia_frequency"`
}

func (a API) IssueFromURL(ctx context.Context, url string) (IssueResp, error) {
	var issueResp IssueResp

	resp, err := a.req(ctx, url)
	if err != nil {
		return issueResp, fmt.Errorf("client.Do: %w", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return issueResp, fmt.Errorf("io.ReadAll: %w", err)
	}

	if err := json.Unmarshal(data, &issueResp); err != nil {
		return issueResp, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return issueResp, nil
}

func (a API) Issue(ctx context.Context, req IssueReq) (IssueResp, error) {
	prefix := a.Prefix
	if prefix == "" {
		prefix = DefaultPrefix
	}

	uu, err := req.URL(prefix)
	if err != nil {
		return IssueResp{}, fmt.Errorf("failed to construct URL: %w", err)
	}

	return a.IssueFromURL(ctx, uu)
}
