# Lint All Your YAML Files

Using this GitHub Action in your workflow to lint all yaml files and then annotates every finding in the changed files view.

![annotation](images/annotation.png)

## Usage

Create a new workflow with the following content:

```yaml
name: YAMLlint

on: push

jobs:
  yamllint:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v1
      - name: Lint and annotate
        uses: "docker://registry.staffbase.com/public/yamllint-action:latest"
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          ## The target path is processed recursively
          TARGETPATH: <relative-folder-path>

```

## Credits

This action is using [adrienverge/yamllint](https://github.com/adrienverge/yamllint).
