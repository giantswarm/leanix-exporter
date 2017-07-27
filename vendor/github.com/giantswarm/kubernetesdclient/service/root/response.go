package root

type VersionResponse struct {
	GitCommit string `json:"git_commit"`
	GoVersion string `json:"go_version"`
	OsArch    string `json:"os_arch"`
	Version   string `json:"version"`
}

type Response struct {
	Description string          `json:"description"`
	Name        string          `json:"name"`
	Source      string          `json:"source"`
	Version     VersionResponse `json:"version"`
}
