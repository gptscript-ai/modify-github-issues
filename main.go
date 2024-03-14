package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/google/go-github/v60/github"
	"github.com/sirupsen/logrus"
)

type args struct {
	Command string `json:"command"`
	Number  string `json:"number"`
	Repo    string `json:"repo"`
	Comment string `json:"comment"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if len(os.Args) != 2 {
		logrus.Errorf("Usage: %s <JSON parameters>", os.Args[0])
		os.Exit(1)
	}

	var a args
	if err := json.Unmarshal([]byte(os.Args[1]), &a); err != nil {
		logrus.Errorf("Error parsing JSON: %s", err)
		os.Exit(1)
	}

	gh := github.NewClient(nil)
	if os.Getenv("GPTSCRIPT_GITHUB_TOKEN") != "" {
		gh = gh.WithAuthToken(os.Getenv("GPTSCRIPT_GITHUB_TOKEN"))
	}

	switch a.Command {
	case "comment":
		if err := comment(ctx, gh, a); err != nil {
			logrus.Errorf("Error commenting: %s", err)
			os.Exit(1)
		}
		fmt.Println("Commented successfully")
	case "close":
		if err := closeIssue(ctx, gh, a); err != nil {
			logrus.Errorf("Error closing issue: %s", err)
			os.Exit(1)
		}
		fmt.Println("Issue closed successfully")
	default:
		logrus.Errorf("Unknown command: %s", a.Command)
		os.Exit(1)
	}
}

func comment(ctx context.Context, gh *github.Client, a args) error {
	if len(strings.Split(a.Repo, "/")) != 2 {
		return fmt.Errorf("invalid repo format (should be 'owner/repo'): %s", a.Repo)
	}

	owner, repo := strings.Split(a.Repo, "/")[0], strings.Split(a.Repo, "/")[1]
	num, err := strconv.Atoi(a.Number)
	if err != nil {
		return fmt.Errorf("invalid issue number: %s", a.Number)
	}

	_, _, err = gh.Issues.CreateComment(ctx, owner, repo, num, &github.IssueComment{
		Body: &a.Comment,
	})
	return err
}

func closeIssue(ctx context.Context, gh *github.Client, a args) error {
	if len(strings.Split(a.Repo, "/")) != 2 {
		return fmt.Errorf("invalid repo format (should be 'owner/repo'): %s", a.Repo)
	}

	owner, repo := strings.Split(a.Repo, "/")[0], strings.Split(a.Repo, "/")[1]
	num, err := strconv.Atoi(a.Number)
	if err != nil {
		return fmt.Errorf("invalid issue number: %s", a.Number)
	}

	_, _, err = gh.Issues.Edit(ctx, owner, repo, num, &github.IssueRequest{
		State: github.String("closed"),
	})
	return err
}
