package toolboxcommands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const auditThresholdUnset = -1
const auditDefaultTop = 10
const auditDefaultWindow = "12 months"

type auditOptions struct {
	kind                                                                string
	repoRoot                                                            string
	excludes                                                            []string
	top                                                                 int
	window                                                              string
	watchLines, flagLines, watchWords, flagWords, watchFiles, flagFiles int
	present                                                             map[string]bool
}

func handleAuditMetrics(ctx commandruntime.ExecutionContext, args []string) resultmodel.CommandResult {
	options, err := parseAuditOptions(ctx.RepositoryRoot, args)
	if err != nil {
		return usageResult(CommandAuditMetrics, err.Error())
	}
	root, err := auditRepositoryRoot(options.repoRoot)
	if err != nil {
		return auditFailure(options.kind, err)
	}
	options.repoRoot = root
	var payload *resultmodel.AuditMetricsResult
	var markdown string
	switch options.kind {
	case "inventory", "folders":
		report, computeErr := computeAuditInventory(root, options.excludes)
		if computeErr != nil {
			return auditFailure(options.kind, computeErr)
		}
		payload, markdown = renderAuditInventory(options, report)
	case "churn", "hotspots":
		report, computeErr := computeAuditChurn(root, options.window, options.excludes)
		if computeErr != nil {
			return auditFailure(options.kind, computeErr)
		}
		payload, markdown = renderAuditChurn(options, report)
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, AuditMetrics: payload, ExactTextOutput: &markdown}
}

func parseAuditOptions(defaultRoot string, args []string) (auditOptions, error) {
	o := auditOptions{repoRoot: defaultRoot, top: auditDefaultTop, window: auditDefaultWindow,
		watchLines: -1, flagLines: -1, watchWords: -1, flagWords: -1, watchFiles: -1, flagFiles: -1, present: map[string]bool{}}
	if len(args) == 0 {
		return o, fmt.Errorf("audit-metrics requires inventory | folders | churn | hotspots")
	}
	o.kind, args = args[0], args[1:]
	if o.kind != "inventory" && o.kind != "folders" && o.kind != "churn" && o.kind != "hotspots" {
		return o, fmt.Errorf("unknown audit-metrics subcommand %q", o.kind)
	}
	for i := 0; i < len(args); i++ {
		name := args[i]
		if !strings.HasPrefix(name, "--") {
			return o, fmt.Errorf("audit-metrics %s: unexpected argument(s): %s", o.kind, strings.Join(args[i:], " "))
		}
		value, valueErr := optionValue(args, &i, strings.SplitN(name, "=", 2)[0])
		if valueErr != nil {
			return o, valueErr
		}
		optionName := strings.SplitN(name, "=", 2)[0]
		o.present[optionName] = true
		switch optionName {
		case "--repo-root":
			o.repoRoot = value
		case "--exclude-path":
			o.excludes = append(o.excludes, value)
		case "--top-count":
			o.top, valueErr = parseAuditCount("--top-count", value)
		case "--since-window":
			o.window = value
		case "--watch-lines":
			o.watchLines, valueErr = parseAuditCount("--watch-lines", value)
		case "--flag-lines":
			o.flagLines, valueErr = parseAuditCount("--flag-lines", value)
		case "--watch-words":
			o.watchWords, valueErr = parseAuditCount("--watch-words", value)
		case "--flag-words":
			o.flagWords, valueErr = parseAuditCount("--flag-words", value)
		case "--watch-files":
			o.watchFiles, valueErr = parseAuditCount("--watch-files", value)
		case "--flag-files":
			o.flagFiles, valueErr = parseAuditCount("--flag-files", value)
		default:
			return o, fmt.Errorf("unknown option %q", name)
		}
		if valueErr != nil {
			return o, valueErr
		}
	}
	if o.kind != "inventory" && (o.present["--watch-lines"] || o.present["--flag-lines"] || o.present["--watch-words"] || o.present["--flag-words"]) {
		return o, fmt.Errorf("file band flags are valid only for inventory")
	}
	if o.kind != "folders" && (o.present["--watch-files"] || o.present["--flag-files"]) {
		return o, fmt.Errorf("folder band flags are valid only for folders")
	}
	if o.kind != "churn" && o.kind != "hotspots" && o.present["--since-window"] {
		return o, fmt.Errorf("--since-window is valid only for churn or hotspots")
	}
	return o, nil
}

func parseAuditCount(name, value string) (int, error) {
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s requires an integer", name)
	}
	return n, nil
}

func auditFailure(kind string, err error) resultmodel.CommandResult {
	f := toolboxFinding(CommandAuditMetrics, "AUDIT-METRICS-FAILED", resultmodel.SeverityError, nil,
		fmt.Sprintf("audit-metrics %s: %v", kind, err), resultmodel.FixabilityManual, "Git or the filesystem declined to provide complete measurements")
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: []resultmodel.CommandFinding{f}}
}

func auditGit(root string, args ...string) (string, error) {
	command := exec.Command("git", append([]string{"-C", root}, args...)...)
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, stderr.String())
	}
	return stdout.String(), nil
}

func auditRepositoryRoot(path string) (string, error) {
	output, err := auditGit(path, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func auditTrackedFiles(root string) ([]string, error) {
	output, err := auditGit(root, "ls-files", "-z")
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, path := range strings.Split(output, "\x00") {
		if path != "" {
			paths = append(paths, path)
		}
	}
	return paths, nil
}

func auditExcluded(path string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

type auditFile struct {
	path         string
	lines, words int
	binary       bool
}

func measureAuditFile(root, relative string) (auditFile, error) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
	if err != nil {
		return auditFile{}, err
	}
	measurement := auditFile{path: relative}
	sniff := data
	if len(sniff) > 8192 {
		sniff = sniff[:8192]
	}
	if bytes.IndexByte(sniff, 0) >= 0 {
		measurement.binary = true
		return measurement, nil
	}
	measurement.lines = bytes.Count(data, []byte{'\n'})
	if len(data) > 0 && data[len(data)-1] != '\n' {
		measurement.lines++
	}
	measurement.words = len(bytes.Fields(data))
	return measurement, nil
}

type auditDistribution struct{ median, p90, p95, max int }

func auditPercentile(values []int, percent int) int {
	if len(values) == 0 {
		return 0
	}
	sorted := append([]int(nil), values...)
	sort.Ints(sorted)
	rank := (percent*len(sorted) + 99) / 100
	if rank < 1 {
		rank = 1
	}
	if rank > len(sorted) {
		rank = len(sorted)
	}
	return sorted[rank-1]
}
func auditSummary(values []int) auditDistribution {
	return auditDistribution{auditPercentile(values, 50), auditPercentile(values, 90), auditPercentile(values, 95), auditPercentile(values, 100)}
}
func auditBand(value, watch, flag int) string {
	if flag != auditThresholdUnset && value > flag {
		return "FLAG"
	}
	if watch != auditThresholdUnset && value > watch {
		return "WATCH"
	}
	return ""
}
