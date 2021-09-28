/*
An example of using goaction with Github APIs.

Before reading this example, please go through the (simpler example) https://github.com/posener/goaction-example

To see how this example works, feel free to create (an issue) https://github.com/posener/goaction-issues-example/issues in this repository.
*/
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v31/github"
	"github.com/posener/goaction"
	"github.com/posener/goaction/actionutil"
	"github.com/posener/goaction/log"
)

//goaction:required
//goaction:description A token for Github APIs.
var token = os.Getenv("GITHUB_TOKEN")

func main() {
	ctx := context.Background()

	// Checking if we are running under Github action mode can be done using the `goaction.CI` flag.
	// This enable us having different flows for running under Github action or in the command line
	// (using `go run` for example). Here we decide to handle only Github action flow:
	if !goaction.CI {
		log.Debugf("Not in Github action mode, quiting.")
		return
	}

	// Check which Github action flow can be done using the `goaction.Event` value. Here we decide
	// only to act in issue mode:
	if goaction.Event != goaction.EventIssues {
		log.Debugf("Not an issue action, nothing to do here.")
		return
	}

	// Since we are in issue flow, getting issue information can be done using the `GetIssue`
	// function. Each flow "Foo" has its own `Get"Foo"` function.
	issue, err := goaction.GetIssues()
	if err != nil {
		log.Fatalf("Failed getting issue information: %s", err)
	}

	// Create a Github Client using the token provided through environment.
	if token == "" {
		log.Fatalf("Token was not provided, please define the Github action 'with' 'github-token' as '${{ secrets.GITHUB_TOKEN }}'")
	}
	gh := actionutil.NewClientWithToken(ctx, token)

	// Interact with the create issue according to the triggering action:
	switch issue.GetAction() {
	case "opened":
		_, _, err := gh.IssuesCreateComment(ctx, issue.GetIssue().GetNumber(), &github.IssueComment{
			Body: github.String(fmt.Sprintf("Hey there %s! Thanks for trying %s", goaction.Actor, goaction.ActionID)),
		})
		if err != nil {
			log.Fatalf("Failed commenting: %s", err)
		}
	case "closed":
		_, _, err := gh.IssuesCreateComment(ctx, issue.GetIssue().GetNumber(), &github.IssueComment{
			Body: github.String("Thanks for cleaning up!"),
		})
		if err != nil {
			log.Fatalf("Failed commenting: %s", err)
		}
	case "reponed":
		_, _, err := gh.IssuesCreateComment(ctx, issue.GetIssue().GetNumber(), &github.IssueComment{
			Body: github.String("Welcome back!"),
		})
		if err != nil {
			log.Fatalf("Failed commenting: %s", err)
		}
	case "edited":
		_, _, err := gh.IssuesCreateComment(ctx, issue.GetIssue().GetNumber(), &github.IssueComment{
			Body: github.String("I'll always have the last word!"),
		})
		if err != nil {
			log.Fatalf("Failed commenting: %s", err)
		}
	case "labeled":
		if issue.GetLabel().GetName() == "bug" {
			_, _, err := gh.IssuesCreateComment(ctx, issue.GetIssue().GetNumber(), &github.IssueComment{
				Body: github.String("Really?? A bug? No way!"),
			})
			if err != nil {
				log.Fatalf("Failed commenting: %s", err)
			}
			_, err = gh.IssuesDeleteLabel(ctx, "bug")
			if err != nil {
				log.Fatalf("Failed commenting: %s", err)
			}
		} else {
			log.Warnf("Ignoring label %s", issue.GetLabel().GetName())
		}
	case "deleted", "transferred", "pinned", "unpinned", "assigned", "unassigned",
		"unlabeled", "locked", "unlocked", "milestoned", "demilestoned":
		log.Fatalf("Unexpected issue action %s", issue.GetAction())
	}
}
