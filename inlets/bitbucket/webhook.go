package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/muety/telepush/config"
	"github.com/muety/telepush/inlets"
	"github.com/muety/telepush/model"
	"github.com/muety/telepush/resolvers"
	"github.com/muety/telepush/util"
)

type BitbucketWebhookInlet struct{}

func New() inlets.Inlet {
	return &BitbucketWebhookInlet{}
}

func (i *BitbucketWebhookInlet) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventKey := r.Header.Get("X-Event-Key")

		var payload Payload
		j := json.NewDecoder(r.Body)
		if err := j.Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		message := &model.DefaultMessage{
			Text: buildMessage(eventKey, &payload),
			Type: resolvers.TextType,
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, config.KeyMessage, message)
		ctx = context.WithValue(ctx, config.KeyParams, &model.MessageParams{DisableLinkPreviews: true})

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func buildMessage(eventKey string, payload *Payload) string {
	switch eventKey {
	// A build system, CI tool, or another vendor recognizes that
	// a commit has a new status and updates the commit with its status
	case "repo:commit_status_created", "repo:commit_status_updated":
		if payload.CommitStatus != nil {
			var emoji string
			switch payload.CommitStatus.State {
			case "INPROGRESS":
				emoji = "⌛️"
			case "SUCCESSFUL":
				emoji = "✅"
			case "FAILED":
				emoji = "❌"
			}
			return fmt.Sprintf(
				"%s *%s*: [%s](%s)\n%s",
				emoji,
				util.EscapeMarkdown(payload.Repository.Name),
				payload.CommitStatus.State,
				payload.CommitStatus.URL,
				util.EscapeMarkdown(payload.CommitStatus.Name),
			)
		}
		fallthrough

		// A user pushes 1 or more commits to a repository
	case "repo:push":
		fallthrough

	// A user forks a repository
	case "repo:fork":
		fallthrough

	// A user updates the  Name,  Description,  Website or Language fields
	// under the Repository details page of the repository settings.
	case "repo:updated":
		fallthrough

	// A repository transfer is accepted
	case "repo:transfer":
		fallthrough

	// A user comments on a commit in a repository
	case "repo:commit_comment_created":
		fallthrough

	// A user creates an issue for a repository
	case "issue:created":
		fallthrough

	// A user updated an issue for a repository
	case "issue:updated":
		fallthrough

	// A user comments on an issue associated with a repository
	case "issue:comment_created":
		fallthrough

	// A user creates a pull request for a repository
	case "pullrequest:created":
		fallthrough

	// A user updates a pull request for a repository
	case "pullrequest:updated":
		fallthrough

	// A user approves a pull request for a repository.
	case "pullrequest:approved":
		fallthrough

	// A user removes an approval from a pull request for a repository
	case "pullrequest:unapproved":
		fallthrough

	// A user merges a pull request for a repository
	case "pullrequest:fulfilled":
		fallthrough

	// A user declines a pull request for a repository
	case "pullrequest:rejected":
		fallthrough

	// A user comments on a pull request
	case "pullrequest:comment_created":
		fallthrough

	// A user updates a comment on a pull request
	case "pullrequest:comment_updated":
		fallthrough

	// A user deletes a comment on a pull request
	case "pullrequest:comment_deleted":
		fallthrough

	default:
		return fmt.Sprintf("Event %s triggered", util.EscapeMarkdown(eventKey))
	}
}
