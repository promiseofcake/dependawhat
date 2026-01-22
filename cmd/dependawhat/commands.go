package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/promiseofcake/dependawhat/internal/scm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	checkCmd = &cobra.Command{
		Use:   "check [owner/repo...]",
		Short: "Check for open Dependabot PRs across repositories",
		Long: `Check for open Dependabot pull requests across multiple repositories.

If no repositories are specified as arguments, checks all repositories
configured in the 'repositories' section of your config file.

You can specify multiple repositories: check owner1/repo1 owner2/repo2

This is a read-only operation - it only displays PR information and does
not perform any actions on the PRs.`,
		RunE: runCheck,
	}
)

func runCheck(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get GitHub token
	token := viper.GetString("github-token")
	if token == "" {
		return fmt.Errorf("GitHub token not provided. Use --github-token flag or set USER_GITHUB_TOKEN environment variable")
	}

	// Get list of repositories to check
	var repos []string

	if len(args) > 0 {
		// Use command-line arguments
		repos = args
	} else {
		// Get all configured repositories from the main config
		repoMap := viper.GetStringMap("repositories")
		for repo := range repoMap {
			repos = append(repos, repo)
		}

		// If no repositories in main config, check for legacy check.repositories
		if len(repos) == 0 {
			repos = viper.GetStringSlice("check.repositories")
		}
	}

	if len(repos) == 0 {
		return fmt.Errorf("no repositories specified. Use command-line arguments or configure repositories in config file")
	}

	// Create GitHub client
	c := scm.NewGithubClient(http.DefaultClient, token)

	fmt.Println("Open Dependabot PRs:")
	fmt.Println("-------------------------")

	for _, repoPath := range repos {
		parts := strings.Split(repoPath, "/")
		if len(parts) != 2 {
			fmt.Printf("  Invalid repository format: %s (expected owner/repo)\n\n", repoPath)
			continue
		}

		owner, repo := parts[0], parts[1]
		fmt.Printf("%s/%s\n", owner, repo)

		// Build query with deny lists
		repoKey := fmt.Sprintf("%s/%s", owner, repo)

		// Get deny lists - merge global and repo-specific
		deniedPackages := getStringSlice("global.denied_packages")
		deniedOrgs := getStringSlice("global.denied_orgs")

		// Add repo-specific denies
		deniedPackages = append(deniedPackages, getStringSlice("repositories."+repoKey+".denied_packages")...)
		deniedOrgs = append(deniedOrgs, getStringSlice("repositories."+repoKey+".denied_orgs")...)

		// Remove duplicates
		deniedPackages = removeDuplicates(deniedPackages)
		deniedOrgs = removeDuplicates(deniedOrgs)

		q := scm.DependencyUpdateQuery{
			Owner:          owner,
			Repo:           repo,
			DeniedPackages: deniedPackages,
			DeniedOrgs:     deniedOrgs,
		}

		// Get open Dependabot PRs with deny list info
		prs, err := c.GetDependabotPRsWithDenyList(ctx, q)
		if err != nil {
			fmt.Printf("   Error: %v\n\n", err)
			continue
		}

		if len(prs) == 0 {
			fmt.Println("   (no open Dependabot PRs)")
		} else {
			for _, pr := range prs {
				if pr.Skipped {
					fmt.Printf("   #%d: %s\n", pr.Number, pr.Title)
					fmt.Printf("   %s\n", pr.URL)
					fmt.Printf("   Status: SKIPPED (%s)\n", pr.SkipReason)
				} else {
					fmt.Printf("   #%d: %s\n", pr.Number, pr.Title)
					fmt.Printf("   %s\n", pr.URL)
					if pr.Status != "" {
						statusIcon := "[pending]"
						if pr.Status == "success" {
							statusIcon = "[success]"
						} else if pr.Status == "failure" {
							statusIcon = "[failure]"
						}
						fmt.Printf("   Status: %s %s\n", statusIcon, pr.Status)
					}
				}
				fmt.Println()
			}
		}
		fmt.Println()
	}

	return nil
}

// Helper functions

func getStringSlice(key string) []string {
	if viper.IsSet(key) {
		return viper.GetStringSlice(key)
	}
	return []string{}
}

func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range slice {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if normalized != "" && !seen[normalized] {
			seen[normalized] = true
			result = append(result, item)
		}
	}

	return result
}
