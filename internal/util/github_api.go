package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Callback func(resp *http.Response) error

type Repo struct {
	Id   int64  `json:"id,string"`
	Name string `json:"name"`
}

type GitHubApi struct {
	Token string
	Org   string
}

func (h *GitHubApi) ListRepos() ([]Repo, error) {
	repos := make([]Repo, 0)
	err := h.getJson(fmt.Sprintf("orgs/%s/repos", h.Org), func(resp *http.Response) error {
		if resp.StatusCode != 200 {
			return fmt.Errorf("wrong response code %d", resp.StatusCode)
		}
		return json.NewDecoder(resp.Body).Decode(&repos)
	})
	return repos, err
}

func (h *GitHubApi) getJson(path string, cb Callback) error {

	tr := &http.Transport{
		DisableKeepAlives:  true,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", "https://api.github.com", path), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2026-03-10")
	if h.Token != "" {
		req.Header.Set("Authorization", "Bearer "+h.Token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return cb(resp)
}
