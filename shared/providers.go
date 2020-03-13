package shared

type Provider int

const (
	GitHubProvider Provider = iota
	GitLabProvider
	BitBucketProvider
)
