// SPDX-FileCopyrightText: 2020 Alvar Penning
//
// SPDX-License-Identifier: GPL-3.0-or-later

package git

import (
	"context"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// FetchGitHubRepos information from a GitHub organization.
//
// The ref specifies the reference to be fetched, e.g., refs/heads/master.
func FetchGitHubRepos(orga, ref, token string) (repos []Repo, err error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	type repo struct {
		Name   string
		SshUrl string
		Ref    struct {
			Name   string
			Target struct {
				Oid string
			}
		} `graphql:"ref(qualifiedName: $ref)"`
	}

	var query struct {
		Viewer struct {
			Organization struct {
				Repositories struct {
					Nodes    []repo
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
				} `graphql:"repositories(first: 100, after: $repoCursor)"`
			} `graphql:"organization(login: $orgaName)"`
		}
	}

	var opts = map[string]interface{}{
		"orgaName":   githubv4.String(orga),
		"repoCursor": (*githubv4.String)(nil),
		"ref":        githubv4.String(ref),
	}

	for {
		if err = client.Query(context.Background(), &query, opts); err != nil {
			repos = nil
			return
		}

		for _, r := range query.Viewer.Organization.Repositories.Nodes {
			repo := Repo{
				Name:   r.Name,
				SshUrl: r.SshUrl,
			}
			if r.Ref.Name != "" {
				repo.Rev = map[string]string{
					r.Ref.Name: r.Ref.Target.Oid,
				}
			}

			repos = append(repos, repo)
		}

		if !query.Viewer.Organization.Repositories.PageInfo.HasNextPage {
			return
		}
		opts["repoCursor"] = query.Viewer.Organization.Repositories.PageInfo.EndCursor
	}
}
