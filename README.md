# dependawhat

A read-only CLI tool to check for open Dependabot pull requests across your GitHub repositories.

## Installation

```bash
go install github.com/promiseofcake/dependawhat/cmd/dependawhat@latest
```

## Usage

Set your GitHub token:

```bash
export USER_GITHUB_TOKEN=your_token_here
```

Check specific repositories:

```bash
dependawhat owner/repo1 owner/repo2
```

Or configure repositories in `~/.dependawhat/config.yaml` and run:

```bash
dependawhat
```

## Configuration

Create `~/.dependawhat/config.yaml`:

```yaml
# Global deny lists apply to all repositories
global:
  denied_packages:
    - github.com/pkg/errors
    - github.com/dgrijalva/jwt-go
  denied_orgs:
    - datadog
    - elastic

# Repositories to check (when running `dependawhat check` with no args)
repositories:
  myorg/repo1: {}
  myorg/repo2:
    denied_packages:
      - github.com/aws/aws-sdk-go
    denied_orgs:
      - hashicorp
```

## Output

```
Open Dependabot PRs:
-------------------------
myorg/repo1
   #123: Bump golang.org/x/net from 0.17.0 to 0.23.0
   https://github.com/myorg/repo1/pull/123
   Status: [success] success

   #124: Bump github.com/datadog/datadog-go from 5.0.0 to 5.1.0
   https://github.com/myorg/repo1/pull/124
   Status: SKIPPED (org 'datadog' is denied)
```
