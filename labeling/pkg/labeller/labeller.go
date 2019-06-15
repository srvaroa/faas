package labeller

import (
	"context"
	"regexp"

	gh "github.com/google/go-github/v26/github"
)

type LabellerConfig map[string]LabelMatcher
type LabelMatcher struct {
	Title string "json:title"
}

// LabelUpdates Represents a request to update the set of labels
type LabelUpdates struct {
	set map[string]bool
}

type Labeller struct {
	// TODO: use the upstream config
	fetchRepoConfig    func(owner string, repoName string) (LabellerConfig, error)
	replaceLabelsForPr func(owner string, repoName string, prNumber int, labels []string) error
	getCurrentLabels   func(owner string, repoName string, prNumber int) ([]string, error)
}

func NewLabeller(github *gh.Client) *Labeller {
	l := Labeller{

		fetchRepoConfig: func(owner string, repoName string) (LabellerConfig, error) {
			return LabellerConfig{}, nil
		},

		replaceLabelsForPr: func(owner string, repoName string, prNumber int, labels []string) error {
			_, _, err := github.Issues.ReplaceLabelsForIssue(
				context.Background(), owner, repoName, prNumber, labels)
			return err
		},

		getCurrentLabels: func(owner string, repoName string, prNumber int) ([]string, error) {
			opts := gh.ListOptions{} // WARN: ignoring pagination here
			currLabels, _, err := github.Issues.ListLabelsByIssue(
				context.Background(), owner, repoName, prNumber, &opts)

			labels := []string{}
			for _, label := range currLabels {
				labels = append(labels, *label.Name)
			}
			return labels, err
		},
	}
	return &l
}

func (l *Labeller) HandleEvent(
	eventName string,
	payload *[]byte) error {

	event, err := gh.ParseWebHook(eventName, *payload)
	if err != nil {
		panic("ARGH")
	}
	switch event := event.(type) {
	case *gh.PullRequestEvent:

		prRepo := event.PullRequest.Base.Repo
		owner := prRepo.GetOwner().GetLogin()
		repoName := *prRepo.Name
		prNumber := *event.PullRequest.Number

		config, err := l.fetchRepoConfig(owner, repoName)

		labelUpdates, err := l.findMatches(event.PullRequest, &config)
		if err != nil {
			panic("ARGH")
		}

		currLabels, err := l.getCurrentLabels(owner, repoName, prNumber)

		newLabels := map[string]bool{}
		for _, label := range currLabels {
			newLabels[label] = true
		}
		for label, b := range labelUpdates.set {
			// Add new ones, or override current ones with false (removal)
			newLabels[label] = b
		}

		err = l.replaceLabelsForPr(owner, repoName, prNumber, collect(&newLabels))

		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Labeller) findMatches(pr *gh.PullRequest, config *LabellerConfig) (LabelUpdates, error) {
	labelUpdates := LabelUpdates{
		set: map[string]bool{},
	}
	for label, matcher := range *config {
		matched, _ := regexp.Match(matcher.Title, []byte(pr.GetTitle()))
		labelUpdates.set[label] = matched
	}
	return labelUpdates, nil
}

// collect takes the values from the set that are set to true
func collect(s *map[string]bool) []string {
	res := []string{}
	for k, v := range *s {
		if v {
			res = append(res, k)
		}
	}
	return res
}
