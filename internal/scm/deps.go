package scm

// DependencyUpdateQuery contains parameters for querying dependency PRs
type DependencyUpdateQuery struct {
	Owner          string
	Repo           string
	DeniedPackages []string // List of package names to exclude
	DeniedOrgs     []string // List of organization names to exclude (e.g., "datadog")
}

// PRInfo contains information about a pull request
type PRInfo struct {
	Number     int
	Title      string
	URL        string
	Status     string // CI status: "success", "failure", "pending", or ""
	Skipped    bool   // Whether PR would be skipped due to deny lists
	SkipReason string // Reason for skipping (denied package/org name)
}
