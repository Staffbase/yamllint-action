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
	"strings"
	"testing"
)

func TestParseInput(t *testing.T) {
	t.Run("standard case", func(t *testing.T) {
		testData := `sources/prod/testapp-deploy.yaml:17:11: [warning] wrong indentation: expected 12 but found 10 (indentation)
sources/prod/testapp2-deploy.yaml:16:11: [warning] wrong indentation: expected 12 but found 10 (indentation)
sources/prod/testapp3-deploy.yaml:16:11: [warning] wrong indentation: expected 12 but found 10 (indentation)
sources/prod/testapp4-deploy.yaml:17:11: [warning] wrong indentation: expected 12 but found 10 (indentation)
sources/stage/testapp-deploy.yaml:17:11: [error] wrong indentation: expected 12 but found 10 (indentation)
sources/stage/testapp2-deploy.yaml:16:11: [warning] wrong indentation: expected 12 but found 10 (indentation)
sources/stage/testapp3-deploy.yaml:16:11: [warning] wrong indentation: expected 12 but found 10 (indentation)
sources/stage/testapp4-deploy.yaml:17:11: [warning] wrong indentation: expected 12 but found 10 (indentation)
sources/dev/testapp-deploy.yaml:17:11: [warning] wrong indentation: expected 12 but found 10 (indentation)
sources/dev/testapp2-deploy.yaml:16:11: [warning] wrong indentation: expected 12 but found 10 (indentation)`

		report := parseInput(strings.NewReader(testData))
		if report.NumFailedLines != 10 {
			t.Errorf("unexpected amount of failed lines: got %d, wanted %d", report.NumFailedLines, 10)
		}

		if report.LinterResults[0].AssertionResults[0].Severity != "warning" {
			t.Errorf("unexpected severity found: got %s, wanted %s", report.LinterResults[0].AssertionResults[0].Severity, "warning")
		}

		if report.LinterResults[6].AssertionResults[0].Severity != "failure" {
			t.Errorf("unexpected severity found: got %s, wanted %s", report.LinterResults[4].AssertionResults[0].Severity, "failure")
		}

		if report.LinterResults[6].AssertionResults[0].Message != "wrong indentation: expected 12 but found 10 (indentation)" {
			t.Errorf("unexpected severity found: got %s, wanted %s", report.LinterResults[4].AssertionResults[0].Message, "wrong indentation: expected 12 but found 10 (indentation)")
		}
	})
}
