package model

// RepoType the type of the repository
type RepoType string

type repoTypeValuesType struct {
	GitHub           RepoType
	TfsGit           RepoType
	Bitbucket        RepoType
	GitHubEnterprise RepoType
}

// RepoTypeValues enum of the type of the repository
var RepoTypeValues = repoTypeValuesType{
	GitHub:           "GitHub",
	TfsGit:           "TfsGit",
	Bitbucket:        "Bitbucket",
	GitHubEnterprise: "GitHubEnterprise",
}
