package main

import (
	"os"
	"path/filepath"
	"regexp"
)

// repoFileMentionPattern extracts repo-relative file-path mentions from REQ/UR
// Markdown bodies: one or more directory segments followed by a filename whose
// dot extension starts with a letter — so "and/or", bare directories, numeric
// ratios ("2.0/5.75"), and placeholder paths containing <angle-bracket> tokens
// never match. MUST stay in lock-step with the file-path alternative of
// bodyMentionPattern in web/board.js: the frontend only decorates paths this
// scanner has classified, so a drift makes mentions silently fall back to
// plain text.
var repoFileMentionPattern = regexp.MustCompile(
	`(?:[A-Za-z0-9_@-]+(?:\.[A-Za-z0-9_-]+)*/)+[A-Za-z0-9_@-][A-Za-z0-9_@.-]*\.[A-Za-z][A-Za-z0-9]{0,7}`)

// collectRepoFileMentions scans every REQ and UR body for file-path mentions
// and stats each unique one against the repo root, producing mention →
// "exists as a regular file". The frontend renders true as a clickable file
// link (serve mode) and false as a visibly-missing path; a mention absent from
// the map stays plain text. Existence is checked at board-build time, so a
// file created afterwards shows up on the next rebuild (in serve mode: the
// next do-work tree change), not instantly.
func collectRepoFileMentions(board *Board) map[string]bool {
	if board.RepoRoot == "" {
		return nil
	}
	fileMentionExists := map[string]bool{}
	scanMarkdownBody := func(markdownBody string) {
		for _, mentionPath := range repoFileMentionPattern.FindAllString(markdownBody, -1) {
			if _, alreadyChecked := fileMentionExists[mentionPath]; alreadyChecked {
				continue
			}
			fileInfo, statErr := os.Stat(filepath.Join(board.RepoRoot, mentionPath))
			fileMentionExists[mentionPath] = statErr == nil && fileInfo.Mode().IsRegular()
		}
	}
	for _, ticket := range board.AllRequests {
		scanMarkdownBody(ticket.BodyMarkdown)
	}
	for _, userRequest := range board.UserRequests {
		scanMarkdownBody(userRequest.BodyMarkdown)
	}
	return fileMentionExists
}
