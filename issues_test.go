package gcd

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

const superman2023_1Issue = `{
	"api_url": "https://www.comics.org/api/issue/2495111/",
    "series_name": "Superman (2023 series)",
    "descriptor": "1 [Jamal Campbell Cover]",
    "publication_date": "April 2023",
    "price": "4.99 USD",
    "page_count": "36.000",
    "editing": "Jillian Grant (credited) (assistant editor); Paul Kaminski (credited) (editor)",
    "indicia_publisher": "DC Comics",
    "brand": "DC [circle and serifs]",
    "isbn": "",
    "barcode": "76194137950000111",
    "rating": "Ages 13+",
    "on_sale_date": "2023-02-21",
    "indicia_frequency": "monthly",
    "notes": "",
    "variant_of": null,
    "series": "https://www.comics.org/api/series/196803/",
    "story_set": [
        {
            "type": "cover",
            "title": "The Man of Steel: Back in Action!",
            "feature": "Superman",
            "sequence_number": 0,
            "page_count": "2.000",
            "script": "",
            "pencils": "Jamal Campbell (credited) (signed as JC Pryce14)",
            "inks": "Jamal Campbell (credited) (signed as JC Pryce14)",
            "colors": "Jamal Campbell (credited) (signed as JC Pryce14)",
            "letters": "?",
            "editing": "",
            "job_number": "",
            "genre": "superhero",
            "characters": "Superman [Clark Kent; Kal-El]; Jimmy Olsen; Lois Lane; Perry White; Neo Kekoa; Mercy Graves; Lex Luthor; Livewire [Leslie Willis]; Parasite; Silver Banshee; Bizarro; Dr. Pharm; Graft",
            "synopsis": "",
            "notes": "Wraparound cover."
        },
        {
            "type": "comic story",
            "title": "Chapter One: Voices in Your Head",
            "feature": "Superman",
            "sequence_number": 1,
            "page_count": "28.000",
            "script": "Joshua Williamson (credited)",
            "pencils": "Jamal Campbell (credited)",
            "inks": "Jamal Campbell (credited)",
            "colors": "Jamal Campbell (credited)",
            "letters": "Ariana Maher (credited)",
            "editing": "",
            "job_number": "",
            "genre": "superhero",
            "characters": "Superman [Clark Kent; Kal-El]; Livewire [Leslie Willis]; Lex Luthor; Duke Dixon; Neo Kekoa; Jimmy Olsen; Lois Lane; Mercy Graves; LL-01 (Lex Luthor hologram); Parasite [Rudy Jones]; Parasite children; Graft; Dr. Pharm; Bizarro; Martha Kent (flashback); Jonathan Kent (flashback); Jor-El (flashback); Lara (flashback); Superman [Jon Kent] (image); Supergirl [Kara Zor-El] (image); Super-Man of China [Kong Kenan] (image); Superboy [Conner Kent] (image); Otho-Ra (image); Osul-Ra (image); Steel [Natasha Irons] (image); Perry White (image); Silver Banshee (image)",
            "synopsis": "As Superman battles Livewire, he gets unwanted advice from Lex Luthor, who despite being in prison wants to help Superman stop threats to Metropolis.  Superman meets Neo Kekoa, the new chief of the Metropolis SCU, and then returns to the Daily Planet as Clark, where Lois Lane chafes in her new role as the new editor-in-chief.  Superman then investigates a disturbance at LexCorp, where Mercy Graves informs him that the company has been renamed SuperCorp, and Lex has dedicated its resources to serve Superman's needs, whether Superman wants the help or not.",
            "notes": ""
        },
        {
            "type": "credits, title page",
            "title": "",
            "feature": "",
            "sequence_number": 2,
            "page_count": "2.000",
            "script": "",
            "pencils": "",
            "inks": "",
            "colors": "?",
            "letters": "?; typeset",
            "editing": "",
            "job_number": "",
            "genre": "",
            "characters": "",
            "synopsis": "",
            "notes": "Title and credits, between pages 1 and 2 of the story."
        },
        {
            "type": "comic story",
            "title": "Coming to Superman",
            "feature": "Superman",
            "sequence_number": 3,
            "page_count": "2.000",
            "script": "?",
            "pencils": "?",
            "inks": "?",
            "colors": "?",
            "letters": "?",
            "editing": "",
            "job_number": "",
            "genre": "superhero",
            "characters": "Brainiac",
            "synopsis": "",
            "notes": "Two-page teaser for an upcoming storyline."
        }
    ],
    "cover": "https://files1.comics.org//img/gcd/covers_by_id/1614/w400/1614882.jpg"
}`

func TestIssueReq_URL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		req  IssueReq
		want string
	}{
		{
			name: "issue",
			req:  IssueReq{ID: 42},
			want: TestPrefix + "/issue/42/",
		},
		{
			name: "issue+format",
			req:  IssueReq{ID: 42, Format: "json"},
			want: TestPrefix + "/issue/42/?format=json",
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp, err := tt.req.URL(TestPrefix)
			require.NoError(t, err, "req.URL")

			assert.Equal(t, tt.want, resp, "req.URL response")
		})
	}
}

func TestIssues(t *testing.T) {
	t.Parallel()

	var requests []*http.Request

	responses := map[string]string{
		"/api/issue/2495111/": superman2023_1Issue,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r)

		resp, ok := responses[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintln(w, resp)
	}))
	t.Cleanup(server.Close)

	api := API{
		Prefix: "http://" + server.Listener.Addr().String() + "/api/",
	}

	resp, err := api.Issue(context.Background(), IssueReq{ID: 2495111})
	require.NoError(t, err, "api.Issue")

	assert.Len(t, resp.StorySet, 4, "resp.StorySet")
	assert.Equal(t, "April 2023", resp.PublicationDate)
}
