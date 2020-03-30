package models

type Provider int

const (
	GitHubProvider Provider = iota
	GitLabProvider
	BitBucketProvider
)
