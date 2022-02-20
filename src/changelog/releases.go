package changelog

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/russross/blackfriday/v2"
)

const (
	githubReleasesEndpoint string = "https://api.github.com/repos/xvello/lbi-release-tests/releases?per_page=100"
)

type githubRelease struct {
	HtmlUrl     string    `json:"html_url"`
	Id          int       `json:"id"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Body        string    `json:"body"`
}

type Release struct {
	Id          int
	Link        string
	Description string
	CreatedAt   time.Time
	PublishedAt time.Time
}

func (r Release) Date() string {
	return r.CreatedAt.Format("02 Jan. 2006")
}

func RetrieveReleases() ([]Release, error) {
	resp, err := retryablehttp.Get(githubReleasesEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var githubReleases []githubRelease
	if err = json.NewDecoder(resp.Body).Decode(&githubReleases); err != nil {
		return nil, err
	}

	renderer := initRenderer()
	var releases []Release
	for _, r := range githubReleases {
		if r.Prerelease || r.Draft {
			continue
		}
		// Cleanup \r that mess up with blackfriday parsing
		body := strings.ReplaceAll(r.Body, "\r\n", "\n")
		desc := blackfriday.Run([]byte(body), renderer)
		releases = append(releases, Release{
			Id:          r.Id,
			Link:        r.HtmlUrl,
			Description: string(desc),
			CreatedAt:   r.CreatedAt,
			PublishedAt: r.PublishedAt,
		})
	}
	return releases, nil
}
