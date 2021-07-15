# Lint All Your YAML Files

Using this GitHub Action in your workflow to lint all yaml files and then annotates every finding in the changed files view.

![annotation](images/annotation.png)

## Usage

Create a new workflow with the following content:

```yaml
name: YAMLlint

on:
  push:
    branches:
      - '**'
    tags-ignore:
      - '**'

jobs:
  yamllint:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v2

      - name: Lint and Annotate
        uses: staffbase/yamllint-action@v1
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          ## The target path is processed recursively
          target-path: <relative-folder-path>
```

## Credits

This action is using

- [adrienverge/yamllint](https://github.com/adrienverge/yamllint)
- [sdesbure/docker_yamllint](https://github.com/sdesbure/docker_yamllint)
