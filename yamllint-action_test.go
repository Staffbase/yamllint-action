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
    
    if report.LinterResults[7].AssertionResults[0].Severity != "failure" {
      t.Errorf("unexpected severity found: got %s, wanted %s", report.LinterResults[4].AssertionResults[0].Severity, "failure")
    }
  })
}
