package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/v27/github"
	"github.com/ldez/ghactions"
)

type Report struct {
	NumFailedLines int
	Success        bool
	LinterResults  []*LinterResult
}

type LinterResult struct {
	AssertionResults []*AssertionResult
	FilePath         string
}

type AssertionResult struct {
	Message  string
	Status   string
	Column   int
	Line     int
	Severity string
}

func mapSeverity(severity string) string {
	if severity == "warning" {
		return "warning"
	}

	return "failure"
}

func getSortedKeySlice(items map[string]*LinterResult) []string {
	keys := make([]string, len(items))

	i := 0
	for k := range items {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	return keys
}

func parseInput(r io.Reader) Report {
	scanner := bufio.NewScanner(r)
	files := make(map[string]*LinterResult)
	failedLines := 0
	re := regexp.MustCompile(` \[(.*)\]`)

	for scanner.Scan() {
		cols := strings.Split(scanner.Text(), ":")

		if len(cols) < 4 {
			break
		}

		codeLine, _ := strconv.Atoi(cols[1])
		codeCol, _ := strconv.Atoi(cols[2])

		assertionResult := AssertionResult{
			Message:  strings.Split(cols[3], "] ")[1] + ":" + cols[4],
			Line:     codeLine,
			Column:   codeCol,
			Severity: mapSeverity(re.FindStringSubmatch(cols[3])[1]),
		}

		if _, exist := files[cols[0]]; exist == false {
			files[cols[0]] = &LinterResult{FilePath: cols[0]}
		}

		files[cols[0]].AssertionResults = append(files[cols[0]].AssertionResults, &assertionResult)
		failedLines++
	}

	report := Report{
		NumFailedLines: failedLines,
		Success:        failedLines == 0,
	}

	keys := getSortedKeySlice(files)

	for _, key := range keys {
		report.LinterResults = append(report.LinterResults, files[key])
	}

	return report
}

func main() {
	report := parseInput(os.Stdin)

	ctx := context.Background()
	action := ghactions.NewAction(ctx)
	action.OnPush(func(client *github.Client, event *github.PushEvent) error {
		return handlePush(ctx, client, event, report)
	})

	if err := action.Run(); err != nil {
		log.Fatal(err)
	}
}

func handlePush(ctx context.Context, client *github.Client, event *github.PushEvent, report Report) error {
	if report.Success {
		return nil
	}

	head := os.Getenv(ghactions.GithubSha)
	owner, repoName := ghactions.GetRepoInfo()

	// find the action's checkrun
	checkName := "yamllint"
	result, _, err := client.Checks.ListCheckRunsForRef(ctx, owner, repoName, head, &github.ListCheckRunsOptions{
		CheckName: github.String(checkName),
		Status:    github.String("in_progress"),
	})
	if err != nil {
		return err
	}

	if len(result.CheckRuns) == 0 {
		return fmt.Errorf("Unable to find check run for action: %s", checkName)
	}
	checkRun := result.CheckRuns[0]

	// add annotations for test failures
	workspacePath := os.Getenv(ghactions.GithubWorkspace) + "/"
	var annotations []*github.CheckRunAnnotation
	for _, t := range report.LinterResults {
		path := strings.TrimPrefix(t.FilePath, workspacePath)

		if len(t.AssertionResults) > 0 {
			for _, a := range t.AssertionResults {
				annotations = append(annotations, &github.CheckRunAnnotation{
					Path:            github.String(path),
					StartLine:       github.Int(a.Line),
					EndLine:         github.Int(a.Line),
					AnnotationLevel: github.String(a.Severity),
					Title:           github.String(""),
					Message:         github.String(a.Message),
				})
			}
		}
	}

	summary := fmt.Sprintf(
		"Tested lines: %d failed\n",
		report.NumFailedLines,
	)

	// add annotations in #50 chunks
	for i := 0; i < len(annotations); i += 50 {
		end := i + 50

		if end > len(annotations) {
			end = len(annotations)
		}

		output := &github.CheckRunOutput{
			Title:       github.String("Result"),
			Summary:     github.String(summary),
			Annotations: annotations[i:end],
		}

		_, _, err = client.Checks.UpdateCheckRun(ctx, owner, repoName, checkRun.GetID(), github.UpdateCheckRunOptions{
			Name:    checkName,
			HeadSHA: github.String(head),
			Output:  output,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
