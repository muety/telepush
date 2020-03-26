package bitbucket_webhook

import "time"

type Link struct {
	Href string `json:"href"`
}

type Links struct {
	Self    Link `json:"self"`
	HTML    Link `json:"html"`
	Avatar  Link `json:"avatar"`
	Diff    Link `json:"diff"`
	Commits Link `json:"commits"`
}

type User struct {
	Type        string `json:"type"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	DisplayName string `json:"display_name"`
	UUID        string `json:"uuid"`
	Links       Links  `json:"links"`
}

type Project struct {
	Type    string `json:"type"`
	Project string `json:"project"`
	UUID    string `json:"uuid"`
	Links   Links  `json:"links"`
	Key     string `json:"key"`
}

type Repository struct {
	Type      string  `json:"type"`
	Links     Links   `json:"links"`
	UUID      string  `json:"uuid"`
	Project   Project `json:"project"`
	FullName  string  `json:"full_name"`
	Name      string  `json:"name"`
	Website   string  `json:"website"`
	Owner     User    `json:"owner"`
	Scm       string  `json:"scm"`
	IsPrivate bool    `json:"is_private"`
}

type Payload struct {
	Actor        User          `json:"actor"`
	Repository   Repository    `json:"repository"`
	CommitStatus *CommitStatus `json:"commit_status"`
	Push         *Push         `json:"push"`
	Fork         *Fork         `json:"fork"`
}

type CommitStatus struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	State       string    `json:"state"`
	Key         string    `json:"key"`
	URL         string    `json:"url"`
	Type        string    `json:"type"`
	CreatedOn   time.Time `json:"created_on"`
	UpdatedOn   time.Time `json:"updated_on"`
	Links       Links     `json:"links"`
}

type Push struct {
	Changes []Change `json:"changes"`
}

type Change struct {
	New       State    `json:"new"`
	Old       State    `json:"old"`
	Links     Links    `json:"links"`
	Created   bool     `json:"created"`
	Forced    bool     `json:"forced"`
	Closed    bool     `json:"closed"`
	Commits   []Commit `json:"commits"`
	Truncated bool     `json:"truncated"`
}

type State struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Target Commit `json:"target"`
	Links  Links  `json:"links"`
}

type Commit struct {
	Type    string    `json:"type"`
	Hash    string    `json:"hash"`
	Author  User      `json:"author"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
	Parents []Commit  `json:"parents"`
	Links   Links     `json:"links"`
}

type Fork Repository
