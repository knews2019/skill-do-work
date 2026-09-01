package doctor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type TimestampFieldChange struct {
	FieldName  string
	LineNumber int
	OldValue   string
	NewValue   string
	Source     string
}

type TimestampRepairPlan struct {
	RelativePath  string
	ExpectedBytes []byte
	UpdatedBytes  []byte
	Changes       []TimestampFieldChange
}

type TimestampScope string

const (
	TimestampScopeAll     TimestampScope = "all"
	TimestampScopeActive  TimestampScope = "active"
	TimestampScopeArchive TimestampScope = "archive"
)

var repairableTimestampPattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}(?:[ T]\d{2}:\d{2}:\d{2}Z?)?$`)

type timestampField struct {
	name       string
	evidence   requestmodel.FieldEvidence
	parsedTime time.Time
	comparable bool
}

func BuildTimestampPlan(ctx context.Context, snapshot *repositorymodel.RepositorySnapshot, now time.Time) ([]TimestampRepairPlan, []resultmodel.CommandFinding) {
	return buildTimestampPlan(ctx, snapshot, now, TimestampScopeAll)
}

// BuildTimestampPlanForScope reuses the doctor's timestamp grammar while applying the
// legacy source policy: active dirty/untracked records use mtime, archive uses Git only.
func BuildTimestampPlanForScope(ctx context.Context, snapshot *repositorymodel.RepositorySnapshot, now time.Time, scope TimestampScope) ([]TimestampRepairPlan, []resultmodel.CommandFinding) {
	if snapshot == nil {
		return buildTimestampPlan(ctx, snapshot, now, scope)
	}
	filtered := *snapshot
	filtered.RequestFiles = nil
	for _, requestFile := range snapshot.RequestFiles {
		path := requestRepositoryPath(requestFile)
		isArchive := strings.HasPrefix(path, "do-work/archive/")
		if scope == TimestampScopeArchive && isArchive || scope == TimestampScopeActive && !isArchive && (strings.HasPrefix(path, "do-work/queue/") || strings.HasPrefix(path, "do-work/working/")) {
			filtered.RequestFiles = append(filtered.RequestFiles, requestFile)
		}
	}
	return buildTimestampPlan(ctx, &filtered, now, scope)
}

func buildTimestampPlan(ctx context.Context, snapshot *repositorymodel.RepositorySnapshot, now time.Time, scope TimestampScope) ([]TimestampRepairPlan, []resultmodel.CommandFinding) {
	if snapshot == nil {
		return nil, []resultmodel.CommandFinding{doctorFinding("DOCTOR-SNAPSHOT-MISSING", resultmodel.SeverityError, nil, nil,
			"repository snapshot is required", resultmodel.FixabilityManual, "no repository evidence was available", doctorArgv(), doctorJSONArgv())}
	}
	now = now.UTC().Truncate(time.Second)
	if now.IsZero() {
		now = time.Now().UTC().Truncate(time.Second)
	}
	horizon := now.Add(2 * time.Minute)
	plans := []TimestampRepairPlan{}
	findings := []resultmodel.CommandFinding{}
	for _, requestFile := range snapshot.RequestFiles {
		if requestFile.ParsedDocument == nil {
			continue
		}
		fields := effectiveTimestampFields(requestFile)
		if len(fields) == 0 {
			continue
		}
		plannedValues := map[string]time.Time{}
		changeByName := map[string]TimestampFieldChange{}
		for _, field := range fields {
			if !field.comparable {
				findings = append(findings, timestampFinding("TIMESTAMP-UNPARSEABLE", requestFile, field,
					fmt.Sprintf("line %d %s=%q is not a supported timestamp", field.evidence.LineNumber, field.name, field.evidence.ScalarValue)))
				findings = append(findings, timestampRefusal(requestFile, field, "the value is not a calendar-valid supported timestamp"))
				continue
			}
			repairableShape := repairableTimestamp(field.evidence)
			if !repairableShape {
				findings = append(findings, timestampRefusal(requestFile, field, "numeric offsets, fractional seconds, non-ASCII padding, and unsupported shapes are diagnosis-only"))
			}
			if !field.parsedTime.After(horizon) {
				continue
			}
			findings = append(findings, timestampFinding("TIMESTAMP-FUTURE", requestFile, field,
				fmt.Sprintf("line %d %s=%q is later than the two-minute UTC skew allowance", field.evidence.LineNumber, field.name, field.evidence.ScalarValue)))
			if !repairableShape {
				continue
			}
			replacement, source, deriveError := scopedTimestampSource(ctx, snapshot.RepositoryRoot, requestRepositoryPath(requestFile), field.evidence.LineNumber, scope)
			if deriveError != nil {
				findings = append(findings, timestampRefusal(requestFile, field, deriveError.Error()))
				continue
			}
			if replacement.After(now) {
				replacement = now
			}
			plannedValues[field.name] = replacement
			changeByName[field.name] = newTimestampChange(field, replacement, source)
		}

		anchors := []string{"created_at", "claimed_at", "completed_at"}
		var predecessor *timestampField
		for _, anchor := range anchors {
			field := fieldNamed(fields, anchor)
			if field == nil || !field.comparable || !repairableTimestamp(field.evidence) {
				continue
			}
			effective := field.parsedTime
			if planned, found := plannedValues[field.name]; found {
				effective = planned
			}
			if predecessor != nil {
				predecessorTime := predecessor.parsedTime
				if planned, found := plannedValues[predecessor.name]; found {
					predecessorTime = planned
				}
				if effective.Before(predecessorTime) {
					findings = append(findings, timestampFinding("TIMESTAMP-ORDER", requestFile, *field,
						fmt.Sprintf("line %d %s=%q precedes %s=%q", field.evidence.LineNumber, field.name, field.evidence.ScalarValue, predecessor.name, predecessor.evidence.ScalarValue)))
					replacement, source, deriveError := scopedTimestampSource(ctx, snapshot.RepositoryRoot, requestRepositoryPath(requestFile), field.evidence.LineNumber, scope)
					if deriveError != nil {
						findings = append(findings, timestampRefusal(requestFile, *field, deriveError.Error()))
						predecessor = field
						continue
					}
					if replacement.After(now) {
						replacement = now
					}
					if replacement.Before(predecessorTime) {
						replacement = predecessorTime
						source += ", clamped to " + requestmodel.CanonicalTimestamp(predecessorTime)
					}
					plannedValues[field.name] = replacement
					changeByName[field.name] = newTimestampChange(*field, replacement, source)
				}
			}
			predecessor = field
		}
		if len(changeByName) == 0 {
			continue
		}
		document, parseError := requestmodel.ParseDocument(requestFile.ContentBytes)
		if parseError != nil {
			continue
		}
		changes := make([]TimestampFieldChange, 0, len(changeByName))
		for _, field := range fields {
			change, found := changeByName[field.name]
			if !found {
				continue
			}
			if setError := document.SetScalar(field.name, change.NewValue); setError != nil {
				findings = append(findings, timestampRefusal(requestFile, field, setError.Error()))
				continue
			}
			changes = append(changes, change)
		}
		if len(changes) > 0 {
			plans = append(plans, TimestampRepairPlan{RelativePath: requestRepositoryPath(requestFile), ExpectedBytes: append([]byte(nil), requestFile.ContentBytes...), UpdatedBytes: document.DocumentBytes(), Changes: changes})
		}
	}
	sort.Slice(plans, func(leftIndex, rightIndex int) bool {
		return plans[leftIndex].RelativePath < plans[rightIndex].RelativePath
	})
	sortFindings(findings)
	return plans, findings
}

func scopedTimestampSource(ctx context.Context, repositoryRoot, relativePath string, lineNumber int, scope TimestampScope) (time.Time, string, error) {
	if scope == TimestampScopeActive {
		tracked := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "ls-files", "--error-unmatch", "--", relativePath).Run() == nil
		clean := tracked && exec.CommandContext(ctx, "git", "-C", repositoryRoot, "diff", "--quiet", "HEAD", "--", relativePath).Run() == nil
		if !clean {
			info, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(relativePath)))
			if err != nil {
				return time.Time{}, "", fmt.Errorf("file mtime could not derive the field's source instant")
			}
			return info.ModTime().UTC().Truncate(time.Second), "file mtime", nil
		}
	}
	return blameTimestamp(ctx, repositoryRoot, relativePath, lineNumber)
}

func effectiveTimestampFields(requestFile *repositorymodel.RequestFile) []timestampField {
	fields := []timestampField{}
	for fieldName, evidence := range requestFile.TypedRecord.FieldEvidenceByName {
		if !strings.HasSuffix(fieldName, "_at") {
			continue
		}
		parsed, parseError := requestmodel.ParseTimestamp(evidence.ScalarValue)
		fields = append(fields, timestampField{name: fieldName, evidence: evidence, parsedTime: parsed.UTC(), comparable: parseError == nil})
	}
	sort.Slice(fields, func(leftIndex, rightIndex int) bool {
		return fields[leftIndex].evidence.LineNumber < fields[rightIndex].evidence.LineNumber
	})
	return fields
}

func repairableTimestamp(evidence requestmodel.FieldEvidence) bool {
	rawValue := strings.Trim(evidence.RawValue, " \t")
	if len(rawValue) >= 2 && ((rawValue[0] == '\'' && rawValue[len(rawValue)-1] == '\'') ||
		(rawValue[0] == '"' && rawValue[len(rawValue)-1] == '"')) {
		rawValue = strings.Trim(rawValue[1:len(rawValue)-1], " \t")
	}
	trimmed := strings.Trim(evidence.ScalarValue, " \t")
	if rawValue != trimmed {
		return false
	}
	if !repairableTimestampPattern.MatchString(trimmed) {
		return false
	}
	_, parseError := requestmodel.ParseTimestamp(trimmed)
	return parseError == nil
}

func blameTimestamp(ctx context.Context, repositoryRoot, relativePath string, lineNumber int) (time.Time, string, error) {
	command := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "--literal-pathspecs", "blame", "--line-porcelain",
		"-L", fmt.Sprintf("%d,%d", lineNumber, lineNumber), "--", relativePath)
	output, commandError := command.Output()
	if commandError != nil {
		return time.Time{}, "", fmt.Errorf("git blame could not derive the field's source instant")
	}
	lines := strings.Split(string(output), "\n")
	commit := ""
	if len(lines) > 0 && len(strings.Fields(lines[0])) > 0 {
		commit = strings.Fields(lines[0])[0]
	}
	for _, line := range lines {
		if !strings.HasPrefix(line, "author-time ") {
			continue
		}
		epoch, parseError := strconv.ParseInt(strings.TrimSpace(strings.TrimPrefix(line, "author-time ")), 10, 64)
		if parseError != nil {
			break
		}
		if len(commit) > 7 {
			commit = commit[:7]
		}
		return time.Unix(epoch, 0).UTC(), "commit " + commit + " author time", nil
	}
	return time.Time{}, "", fmt.Errorf("git blame returned no author-time evidence")
}

func fieldNamed(fields []timestampField, fieldName string) *timestampField {
	for index := range fields {
		if fields[index].name == fieldName {
			return &fields[index]
		}
	}
	return nil
}

func newTimestampChange(field timestampField, replacement time.Time, source string) TimestampFieldChange {
	return TimestampFieldChange{FieldName: field.name, LineNumber: field.evidence.LineNumber,
		OldValue: field.evidence.ScalarValue, NewValue: requestmodel.CanonicalTimestamp(replacement), Source: source}
}

func timestampFinding(code string, requestFile *repositorymodel.RequestFile, field timestampField, evidence string) resultmodel.CommandFinding {
	return doctorFinding(code, resultmodel.SeverityWarning, nonEmptyString(requestFile.TypedRecord.RequestID), []string{requestRepositoryPath(requestFile)}, evidence,
		resultmodel.FixabilityAutomatic, "diagnosis never mutates without --repair-timestamps", repairDoctorArgv(false, false), doctorJSONArgv())
}

func timestampRefusal(requestFile *repositorymodel.RequestFile, field timestampField, reason string) resultmodel.CommandFinding {
	return doctorFinding("TIMESTAMP-REPAIR-REFUSED", resultmodel.SeverityWarning, nonEmptyString(requestFile.TypedRecord.RequestID), []string{requestRepositoryPath(requestFile)},
		fmt.Sprintf("line %d %s=%q: %s", field.evidence.LineNumber, field.name, field.evidence.ScalarValue, reason),
		resultmodel.FixabilityRefused, "no provably lossless blame-derived rewrite is available", doctorArgv(), doctorJSONArgv())
}
