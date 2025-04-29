/*
Copyright 2021, Staffbase GmbH and contributors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

	"github.com/google/go-github/v71/github"
	"github.com/ldez/ghactions"
)

type Report struct {
	NumFailedLines  int
	Success         bool
	ErrorHasOccured bool
	LinterResults   []*LinterResult
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
	ErrorHasOccured := false
	re := regexp.MustCompile(` \[(.*)]`)

	for scanner.Scan() {
		cols := strings.Split(scanner.Text(), ":")

		if len(cols) < 4 {
			log.Println(scanner.Text())
			break
		}

		codeLine, _ := strconv.Atoi(cols[1])
		codeCol, _ := strconv.Atoi(cols[2])
		fileName := cols[0]
		message := strings.Split(cols[3], "] ")[1]

		if len(cols) == 5 {
			message += ":" + cols[4]
		}

		severity := mapSeverity(re.FindStringSubmatch(cols[3])[1])

		if severity == "failure" {
			ErrorHasOccured = true
		}

		assertionResult := AssertionResult{
			Message:  message,
			Line:     codeLine,
			Column:   codeCol,
			Severity: severity,
		}

		if _, exist := files[fileName]; !exist {
			files[fileName] = &LinterResult{FilePath: fileName}
		}

		files[fileName].AssertionResults = append(files[fileName].AssertionResults, &assertionResult)
		failedLines++
	}

	report := Report{
		NumFailedLines:  failedLines,
		Success:         failedLines == 0,
		ErrorHasOccured: ErrorHasOccured,
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
		return handlePush(ctx, client, report)
	})

	if err := action.Run(); err != nil {
		log.Fatal(err)
	}
}

func handlePush(ctx context.Context, client *github.Client, report Report) error {
	if report.Success {
		return nil
	}

	head := os.Getenv(ghactions.GithubSha)
	owner, repoName := ghactions.GetRepoInfo()

	// find the action's checkrun
	checkName := os.Getenv("ACTION_NAME")
	result, _, err := client.Checks.ListCheckRunsForRef(ctx, owner, repoName, head, &github.ListCheckRunsOptions{
		CheckName: github.Ptr(checkName),
		Status:    github.Ptr("in_progress"),
	})
	if err != nil {
		return err
	}

	if len(result.CheckRuns) == 0 {
		return fmt.Errorf("unable to find check run for action: %s", checkName)
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
					Path:            github.Ptr(path),
					StartLine:       github.Ptr(a.Line),
					EndLine:         github.Ptr(a.Line),
					AnnotationLevel: github.Ptr(a.Severity),
					Title:           github.Ptr(""),
					Message:         github.Ptr(a.Message),
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
			Title:       github.Ptr("Result"),
			Summary:     github.Ptr(summary),
			Annotations: annotations[i:end],
		}

		_, _, err = client.Checks.UpdateCheckRun(ctx, owner, repoName, checkRun.GetID(), github.UpdateCheckRunOptions{
			Name:   checkName,
			Output: output,
		})
		if err != nil {
			return err
		}
	}

	if report.ErrorHasOccured {
		return fmt.Errorf("report contains error: %s", summary)
	} else {
		return nil
	}
}
