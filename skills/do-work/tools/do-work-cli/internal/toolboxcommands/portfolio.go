package toolboxcommands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func handlePortfolio(ctx commandruntime.ExecutionContext, args []string) resultmodel.CommandResult {
	rest, dryRun, commit, err := parseMutationFlags(args)
	if err != nil {
		return usageResult(CommandPortfolio, err.Error())
	}
	if len(rest) != 3 && len(rest) != 4 {
		return usageResult(CommandPortfolio, "Usage: publish-portfolio-summary (--canonical-only <source> <canonical> | --with-snapshot <source> <canonical> <candidate>) [--dry-run|--commit]")
	}
	mode := rest[0]
	if mode != "--canonical-only" && mode != "--with-snapshot" {
		return usageResult(CommandPortfolio, "unknown portfolio mode")
	}
	if mode == "--canonical-only" && len(rest) != 3 || mode == "--with-snapshot" && len(rest) != 4 {
		return usageResult(CommandPortfolio, "portfolio mode has wrong argument count")
	}
	source := absoluteFromRoot(ctx.RepositoryRoot, rest[1])
	data, readErr := os.ReadFile(source)
	if readErr != nil {
		return usageResult(CommandPortfolio, "Portfolio source is not a regular file: "+rest[1])
	}
	if info, statErr := os.Stat(source); statErr != nil || !info.Mode().IsRegular() {
		return usageResult(CommandPortfolio, "Portfolio source is not a regular file: "+rest[1])
	}
	canonicalRel, canonicalAbs, pathErr := repositoryPath(ctx.RepositoryRoot, rest[2])
	if pathErr != nil {
		return usageResult(CommandPortfolio, pathErr.Error())
	}
	_, canonicalStatError := os.Lstat(canonicalAbs)
	canonicalExisted := canonicalStatError == nil
	if canonicalStatError != nil && !os.IsNotExist(canonicalStatError) {
		return usageResult(CommandPortfolio, canonicalStatError.Error())
	}
	targets := []string{canonicalRel}
	directories := []string{filepath.ToSlash(filepath.Dir(canonicalRel))}
	snapshotRel := ""
	if mode == "--with-snapshot" {
		snapshotRel, _, pathErr = firstFreePortfolioPath(ctx.RepositoryRoot, rest[3])
		if pathErr != nil {
			return usageResult(CommandPortfolio, pathErr.Error())
		}
		targets = append(targets, snapshotRel)
		directories = append(directories, filepath.ToSlash(filepath.Dir(snapshotRel)))
	}
	createdDirectories := absentTransactionDirectories(ctx.RepositoryRoot, directories...)
	if mode == "--with-snapshot" && !commit && !dryRun {
		snapshotDirectories := absentTransactionDirectories(ctx.RepositoryRoot, filepath.ToSlash(filepath.Dir(snapshotRel)))
		snapshotResult := runTransaction(CommandPortfolio, ctx.RepositoryRoot, []string{snapshotRel}, snapshotDirectories, false, false, "[do-work] Publish portfolio snapshot", func(recorder *gittransaction.MutationRecorder) error {
			for _, directory := range snapshotDirectories {
				if err := rootedMkdirAll(ctx.RepositoryRoot, directory, 0o755); err != nil {
					return err
				}
				if err := recorder.RecordCreatedDirectory(directory); err != nil {
					return err
				}
			}
			if err := rootedPublishFile(ctx.RepositoryRoot, snapshotRel, data, 0o644, false); err != nil {
				return err
			}
			return recorder.RecordCreated(snapshotRel)
		})
		if snapshotResult.Outcome != resultmodel.OutcomeSuccess {
			return snapshotResult
		}
		canonicalDirectories := absentTransactionDirectories(ctx.RepositoryRoot, filepath.ToSlash(filepath.Dir(canonicalRel)))
		canonicalResult := runTransaction(CommandPortfolio, ctx.RepositoryRoot, []string{canonicalRel}, canonicalDirectories, false, false, "[do-work] Publish portfolio canonical", func(recorder *gittransaction.MutationRecorder) error {
			for _, directory := range canonicalDirectories {
				if err := rootedMkdirAll(ctx.RepositoryRoot, directory, 0o755); err != nil {
					return err
				}
				if err := recorder.RecordCreatedDirectory(directory); err != nil {
					return err
				}
			}
			if canonicalExisted {
				if err := rootedPublishFile(ctx.RepositoryRoot, canonicalRel, data, 0o644, true); err != nil {
					return err
				}
				return recorder.RecordTouched(canonicalRel)
			}
			if err := rootedPublishFile(ctx.RepositoryRoot, canonicalRel, data, 0o644, false); err != nil {
				return err
			}
			return recorder.RecordCreated(canonicalRel)
		})
		canonicalResult.Changes = append(snapshotResult.Changes, canonicalResult.Changes...)
		output := snapshotRel + "\n"
		if canonicalResult.Outcome == resultmodel.OutcomeSuccess {
			output = canonicalRel + "\n" + snapshotRel + "\n"
		} else {
			canonicalResult.Findings = append(canonicalResult.Findings, toolboxFinding(CommandPortfolio, "PORTFOLIO-SNAPSHOT-RETAINED", resultmodel.SeverityWarning, []string{snapshotRel}, "immutable snapshot published before canonical failure and intentionally retained", resultmodel.FixabilityManual, "repair the canonical destination without deleting the snapshot"))
		}
		canonicalResult.ExactTextOutput = &output
		return canonicalResult
	}
	result := runTransaction(CommandPortfolio, ctx.RepositoryRoot, targets, createdDirectories, dryRun, commit, "[do-work] Publish portfolio summary", func(recorder *gittransaction.MutationRecorder) error {
		for _, dir := range createdDirectories {
			if err := rootedMkdirAll(ctx.RepositoryRoot, dir, 0o755); err != nil {
				return err
			}
			if err := recorder.RecordCreatedDirectory(dir); err != nil {
				return err
			}
		}
		if snapshotRel != "" {
			if err := rootedPublishFile(ctx.RepositoryRoot, snapshotRel, data, 0o644, false); err != nil {
				return err
			}
			if err := recorder.RecordCreated(snapshotRel); err != nil {
				return err
			}
		}
		if canonicalExisted {
			info, err := os.Lstat(canonicalAbs)
			if err != nil {
				return err
			}
			if info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
				return fmt.Errorf("Portfolio canonical path is unsafe: %s", canonicalRel)
			}
			if err := rootedPublishFile(ctx.RepositoryRoot, canonicalRel, data, 0o644, true); err != nil {
				return err
			}
			return recorder.RecordTouched(canonicalRel)
		}
		if err := rootedPublishFile(ctx.RepositoryRoot, canonicalRel, data, 0o644, false); err != nil {
			return err
		}
		return recorder.RecordCreated(canonicalRel)
	})
	if result.Outcome == resultmodel.OutcomeSuccess {
		paths := []string{canonicalRel}
		if snapshotRel != "" {
			paths = append(paths, snapshotRel)
		}
		output := strings.Join(paths, "\n") + "\n"
		if dryRun {
			output = "would publish\n" + output
		}
		result.ExactTextOutput = &output
	}
	return result
}

func absoluteFromRoot(root, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(root, path)
}
func repositoryPath(root, path string) (string, string, error) {
	absolute := absoluteFromRoot(root, path)
	relative, err := filepath.Rel(root, absolute)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", "", fmt.Errorf("path must stay inside repository root: %s", path)
	}
	relative = filepath.ToSlash(relative)
	if err := validateNoLinkedAncestors(root, relative, false); err != nil {
		return "", "", err
	}
	return relative, absolute, nil
}
func firstFreePortfolioPath(root, candidate string) (string, string, error) {
	_, absolute, err := repositoryPath(root, candidate)
	if err != nil {
		return "", "", err
	}
	extension := filepath.Ext(absolute)
	stem := strings.TrimSuffix(absolute, extension)
	for sequence := 1; ; sequence++ {
		chosen := absolute
		if sequence > 1 {
			chosen = fmt.Sprintf("%s-%d%s", stem, sequence, extension)
		}
		if _, err := os.Lstat(chosen); os.IsNotExist(err) {
			chosenRel, _ := filepath.Rel(root, chosen)
			return filepath.ToSlash(chosenRel), chosen, nil
		}
	}
}
