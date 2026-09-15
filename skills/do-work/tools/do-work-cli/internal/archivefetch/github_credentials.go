package archivefetch

import "net/url"

// TrustedGitHubURL limits ambient GitHub credentials to the HTTPS endpoints used
// for repository downloads. Other hosts, including upstream mirrors, stay anonymous.
func TrustedGitHubURL(destination *url.URL) bool {
	if destination == nil || destination.Scheme != "https" || destination.User != nil {
		return false
	}
	if port := destination.Port(); port != "" && port != "443" {
		return false
	}
	switch destination.Hostname() {
	case "github.com", "api.github.com", "codeload.github.com", "raw.githubusercontent.com":
		return true
	default:
		return false
	}
}
