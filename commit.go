package main

type GitCommit struct {
	Sha    string `json:"sha"`
	Commit struct {
		Message   string `json:"message"`
		Committer struct {
			Name string `json:"name"`
		} `json:"committer"`
	} `json:"commit"`
}

type multiCommits []GitCommit

type NewCommit struct {
	Sha       string
	Message   string
	Committer string
}

type newCommits []NewCommit
