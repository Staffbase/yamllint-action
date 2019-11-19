# yamllint-action

This github action is linting your yaml files and then annotates every finding in the changed files view.

# usage

Create a new workflow with the following settings
```
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
          TARGETPATH: [RELATIVE_FOLDER_PATH]

```