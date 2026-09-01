package corehelpers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

var beforePrivateCopyPublish = func() {}

func handleCaptureScreenshot(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	sourcePath, destinationPath := "", ""
	staged, keepSource := false, false
	for index := 0; index < len(arguments); index++ {
		switch {
		case arguments[index] == "--staged":
			staged = true
		case arguments[index] == "--keep-source":
			keepSource = true
		case arguments[index] == "--source" || strings.HasPrefix(arguments[index], "--source="):
			value, err := optionValue(arguments, &index, "--source")
			if err != nil {
				return usageResult(CommandCaptureScreenshot, err.Error())
			}
			sourcePath = value
		case arguments[index] == "--destination" || strings.HasPrefix(arguments[index], "--destination="):
			value, err := optionValue(arguments, &index, "--destination")
			if err != nil {
				return usageResult(CommandCaptureScreenshot, err.Error())
			}
			destinationPath = value
		default:
			return usageResult(CommandCaptureScreenshot, "unknown option "+arguments[index])
		}
	}
	if sourcePath == "" || destinationPath == "" {
		return usageResult(CommandCaptureScreenshot, "--source and --destination are required")
	}
	if staged == keepSource {
		return usageResult(CommandCaptureScreenshot, "choose exactly one of --staged or --keep-source")
	}
	sourceAbsolute := absoluteFromRoot(executionContext.RepositoryRoot, sourcePath)
	destinationAbsolute := absoluteFromRoot(executionContext.RepositoryRoot, destinationPath)
	sourceIdentity, sourceIdentityError := os.Lstat(sourceAbsolute)
	if sourceIdentityError != nil || !sourceIdentity.Mode().IsRegular() {
		return usageResult(CommandCaptureScreenshot, "source must be a readable regular file")
	}
	if err := publishPrivateCopy(sourceAbsolute, destinationAbsolute); err != nil {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{helperFinding("SCREENSHOT-PUBLISH-FAILED", resultmodel.SeverityError, []string{sourcePath, destinationPath}, err.Error(), resultmodel.FixabilityManual, "source was preserved and no partial destination was authorized", []string{"do-work-cli", CommandCaptureScreenshot, "--keep-source", "--source", sourcePath, "--destination", destinationPath}, []string{"test", "-f", sourcePath})}}
	}
	findings := []resultmodel.CommandFinding{helperFinding("SCREENSHOT-PUBLISHED", resultmodel.SeverityInfo, []string{destinationPath}, "bytes verified and published without overwrite using private mode 0600", resultmodel.FixabilityAutomatic, "", nil, []string{"test", "-s", destinationPath})}
	changes := []resultmodel.RecordedChange{{Path: destinationPath, Kind: "created", Detail: "published private screenshot copy"}}
	if staged {
		currentIdentity, identityError := os.Lstat(sourceAbsolute)
		if identityError != nil || !os.SameFile(sourceIdentity, currentIdentity) {
			findings = append(findings, helperFinding("SCREENSHOT-SOURCE-CLEANUP-WARNING", resultmodel.SeverityWarning, []string{sourcePath}, "source identity changed after publication", resultmodel.FixabilityRefused, "replacement source is not owned by this dispatch", nil, []string{"test", "-f", destinationPath}))
		} else if err := os.Remove(sourceAbsolute); err != nil {
			findings = append(findings, helperFinding("SCREENSHOT-SOURCE-CLEANUP-WARNING", resultmodel.SeverityWarning, []string{sourcePath}, err.Error(), resultmodel.FixabilityManual, "publication succeeded; cleanup can be retried independently", []string{"rm", "--", sourcePath}, []string{"test", "-f", destinationPath}))
		} else {
			changes = append(changes, resultmodel.RecordedChange{Path: sourcePath, Kind: "deleted", Detail: "removed staged source after destination verification"})
			_ = os.Remove(filepath.Dir(sourceAbsolute))
		}
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Changes: changes, Findings: findings}
}

func publishPrivateCopy(sourcePath, destinationPath string) error {
	sourceInfo, err := os.Lstat(sourcePath)
	if err != nil {
		return fmt.Errorf("inspect source: %w", err)
	}
	if !sourceInfo.Mode().IsRegular() {
		return fmt.Errorf("source must be a regular non-symlink file")
	}
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()
	if err := os.MkdirAll(filepath.Dir(destinationPath), 0o755); err != nil {
		return err
	}
	parentRoot, err := os.OpenRoot(filepath.Dir(destinationPath))
	if err != nil {
		return err
	}
	defer parentRoot.Close()
	destinationName := filepath.Base(destinationPath)
	if _, err := parentRoot.Lstat(destinationName); err == nil {
		return fmt.Errorf("destination already exists")
	} else if !os.IsNotExist(err) {
		return err
	}
	token := make([]byte, 8)
	if _, err := rand.Read(token); err != nil {
		return err
	}
	temporaryName := "." + destinationName + ".publishing." + hex.EncodeToString(token)
	temporary, err := parentRoot.OpenFile(temporaryName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	published := false
	defer func() {
		_ = temporary.Close()
		if !published {
			_ = parentRoot.Remove(temporaryName)
		}
	}()
	copied, copyErr := io.Copy(temporary, source)
	if copyErr == nil {
		copyErr = temporary.Sync()
	}
	if closeErr := temporary.Close(); copyErr == nil {
		copyErr = closeErr
	}
	if copyErr != nil {
		return copyErr
	}
	if copied != sourceInfo.Size() {
		return fmt.Errorf("copied %d bytes, expected %d", copied, sourceInfo.Size())
	}
	sourceBytes, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	stagedBytes, err := parentRoot.ReadFile(temporaryName)
	if err != nil {
		return err
	}
	if !bytes.Equal(sourceBytes, stagedBytes) {
		return fmt.Errorf("staged bytes differ from source")
	}
	beforePrivateCopyPublish()
	if err := parentRoot.Link(temporaryName, destinationName); err != nil {
		return fmt.Errorf("publish without overwrite: %w", err)
	}
	destinationInfo, err := parentRoot.Stat(destinationName)
	temporaryInfo, temporaryStatError := parentRoot.Stat(temporaryName)
	if err != nil || temporaryStatError != nil || !os.SameFile(destinationInfo, temporaryInfo) {
		return fmt.Errorf("published destination identity could not be verified")
	}
	published = true
	if err := parentRoot.Remove(temporaryName); err != nil {
		return nil
	}
	return nil
}
