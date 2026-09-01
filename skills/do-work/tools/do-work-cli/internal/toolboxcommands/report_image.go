package toolboxcommands

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/atomicfile"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type imageRequest struct {
	Name        string
	Description string
}

type imageGenerationResult struct {
	Path        string
	Interrupted bool
	Err         error
}

var reportImageLookPath = exec.LookPath

func handleReportImage(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	rest, dryRun, commit, err := parseMutationFlags(arguments)
	if err != nil || len(rest) != 3 {
		if err == nil {
			err = fmt.Errorf("Usage: generate-report-image <absolute-output-png> <style-brief> <sanitized-description> [--dry-run|--commit]")
		}
		return usageResult(CommandReportImage, err.Error())
	}
	if !filepath.IsAbs(rest[0]) {
		return usageResult(CommandReportImage, "Output path must be absolute.")
	}
	if info, statErr := os.Stat(filepath.Dir(rest[0])); statErr != nil || !info.IsDir() {
		return usageResult(CommandReportImage, "Output directory must already exist.")
	}
	relative, absolute, pathErr := repositoryPath(executionContext.RepositoryRoot, rest[0])
	if pathErr != nil {
		return usageResult(CommandReportImage, pathErr.Error())
	}
	if info, statErr := os.Lstat(absolute); statErr == nil && !info.Mode().IsRegular() {
		return usageResult(CommandReportImage, "output target must be a regular file or absent")
	}
	_, preexistingError := os.Lstat(absolute)
	preexisting := preexistingError == nil

	invocationContext, stopSignals, interrupted := imageSignalContext()
	defer stopSignals()
	result := runTransaction(CommandReportImage, executionContext.RepositoryRoot, []string{relative}, nil, dryRun, commit, "[do-work] Generate report image", func(recorder *gittransaction.MutationRecorder) error {
		generated := generateImage(invocationContext, absolute, rest[1], rest[2])
		if generated.Interrupted {
			return context.Canceled
		}
		if generated.Err != nil {
			return generated.Err
		}
		if preexisting {
			return recorder.RecordTouched(relative)
		}
		return recorder.RecordCreated(relative)
	})
	if status := interrupted(); status != 0 {
		result.ExitCodeOverride = status
	}
	if result.Outcome == resultmodel.OutcomeSuccess {
		output := absolute + "\n"
		if dryRun {
			output = "would generate " + absolute + "\n"
		}
		result.ExactTextOutput = &output
	}
	return result
}

func handleReportImageBatch(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	rest, dryRun, commit, err := parseMutationFlags(arguments)
	if err != nil || len(rest) < 3 {
		if err == nil {
			err = fmt.Errorf("Usage: generate-report-image-batch <report-directory> <style-brief> <target-name>:<prompt>... [--dry-run|--commit]")
		}
		return usageResult(CommandReportImageBatch, err.Error())
	}
	reportDirectory := absoluteFromRoot(executionContext.RepositoryRoot, rest[0])
	if info, statErr := os.Stat(reportDirectory); statErr != nil || !info.IsDir() {
		return usageResult(CommandReportImageBatch, "report directory must exist")
	}
	requests := make([]imageRequest, 0, len(rest)-2)
	seen := map[string]struct{}{}
	for _, specification := range rest[2:] {
		name, description, ok := strings.Cut(specification, ":")
		if !ok || name == "" || name == "." || filepath.Base(name) != name || strings.ContainsAny(name, `/\\`) {
			return usageResult(CommandReportImageBatch, "Image target must be a bare filename: "+name)
		}
		if _, duplicate := seen[name]; duplicate {
			return usageResult(CommandReportImageBatch, "duplicate image target: "+name)
		}
		seen[name] = struct{}{}
		requests = append(requests, imageRequest{Name: name, Description: description})
	}
	generatedDirectory := filepath.Join(reportDirectory, "generated")
	if _, statErr := os.Lstat(generatedDirectory); !os.IsNotExist(statErr) {
		return usageResult(CommandReportImageBatch, "REFUSING: generated/ already exists")
	}
	generatedRel, _, pathErr := repositoryPath(executionContext.RepositoryRoot, generatedDirectory)
	if pathErr != nil {
		return usageResult(CommandReportImageBatch, pathErr.Error())
	}
	targets := make([]string, 0, len(requests))
	for _, request := range requests {
		targets = append(targets, filepath.ToSlash(filepath.Join(generatedRel, request.Name)))
	}

	invocationContext, stopSignals, interrupted := imageSignalContext()
	defer stopSignals()
	succeeded := 0
	result := runTransaction(CommandReportImageBatch, executionContext.RepositoryRoot, targets, []string{generatedRel}, dryRun, commit, "[do-work] Generate report image batch", func(recorder *gittransaction.MutationRecorder) error {
		stage, createErr := os.MkdirTemp(reportDirectory, ".generated.staging.*")
		if createErr != nil {
			return createErr
		}
		defer os.RemoveAll(stage)
		outcomes := make([]imageGenerationResult, len(requests))
		var wait sync.WaitGroup
		for index, request := range requests {
			wait.Add(1)
			go func(index int, request imageRequest) {
				defer wait.Done()
				outcomes[index] = generateImage(invocationContext, filepath.Join(stage, request.Name), rest[1], request.Description)
			}(index, request)
		}
		wait.Wait()
		for index, outcome := range outcomes {
			if outcome.Interrupted {
				return context.Canceled
			}
			if outcome.Err != nil {
				_ = os.Remove(filepath.Join(stage, requests[index].Name))
				continue
			}
			succeeded++
		}
		if succeeded == 0 {
			return nil
		}
		if _, statErr := os.Lstat(generatedDirectory); !os.IsNotExist(statErr) {
			return errorsNewGeneratedAppeared()
		}
		if renameErr := os.Rename(stage, generatedDirectory); renameErr != nil {
			return renameErr
		}
		if err := recorder.RecordCreatedDirectory(generatedRel); err != nil {
			return err
		}
		for index, outcome := range outcomes {
			if outcome.Err == nil {
				if err := recorder.RecordCreated(targets[index]); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if status := interrupted(); status != 0 {
		result.ExitCodeOverride = status
	}
	if result.Outcome == resultmodel.OutcomeSuccess {
		output := ""
		if dryRun {
			output = "would generate " + generatedDirectory + "\n"
		} else if succeeded > 0 {
			output = generatedDirectory + "\n"
		}
		result.ExactTextOutput = &output
	}
	return result
}

func errorsNewGeneratedAppeared() error {
	return fmt.Errorf("REFUSING: generated/ appeared before publication")
}

func generateImage(ctx context.Context, outputPath, style, description string) imageGenerationResult {
	parent := filepath.Dir(outputPath)
	stage, err := os.CreateTemp(parent, "."+filepath.Base(outputPath)+".generating.*")
	if err != nil {
		return imageGenerationResult{Path: outputPath, Err: err}
	}
	stagePath := stage.Name()
	if closeErr := stage.Close(); closeErr != nil {
		os.Remove(stagePath)
		return imageGenerationResult{Path: outputPath, Err: closeErr}
	}
	defer os.Remove(stagePath)
	prompt := style + " Content: " + description
	if imagegen, lookupErr := reportImageLookPath("imagegen"); lookupErr == nil {
		outcome := runOwnedProcess(ctx, "", imagegen, "--output", stagePath, "--prompt", prompt)
		if outcome.Interrupted {
			return imageGenerationResult{Path: outputPath, Interrupted: true, Err: outcome.Err}
		}
		if outcome.Status == 0 && nonemptyRegular(stagePath) {
			return publishGeneratedImage(stagePath, outputPath)
		}
	}
	if os.Getenv("DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND") != "1" {
		return imageGenerationResult{Path: outputPath, Err: fmt.Errorf("no report image backend succeeded")}
	}
	codex, lookupErr := reportImageLookPath("codex")
	if lookupErr != nil {
		return imageGenerationResult{Path: outputPath, Err: lookupErr}
	}
	agentDirectory, createErr := os.MkdirTemp("", "do-work-ai-report-image.*")
	if createErr != nil {
		return imageGenerationResult{Path: outputPath, Err: createErr}
	}
	defer os.RemoveAll(agentDirectory)
	outcome := runOwnedProcess(ctx, agentDirectory, codex, "exec", "--dangerously-bypass-approvals-and-sandbox", "Generate a 16:9 image and save the PNG EXACTLY to ./generated.png. "+prompt)
	if outcome.Interrupted {
		return imageGenerationResult{Path: outputPath, Interrupted: true, Err: outcome.Err}
	}
	agentOutput := filepath.Join(agentDirectory, "generated.png")
	if outcome.Status != 0 || !nonemptyRegular(agentOutput) {
		return imageGenerationResult{Path: outputPath, Err: fmt.Errorf("agentic report image backend failed")}
	}
	input, openErr := os.Open(agentOutput)
	if openErr != nil {
		return imageGenerationResult{Path: outputPath, Err: openErr}
	}
	destination, openErr := os.OpenFile(stagePath, os.O_WRONLY|os.O_TRUNC, 0)
	if openErr != nil {
		input.Close()
		return imageGenerationResult{Path: outputPath, Err: openErr}
	}
	_, copyErr := io.Copy(destination, input)
	closeOutputErr := destination.Close()
	closeInputErr := input.Close()
	if copyErr != nil || closeOutputErr != nil || closeInputErr != nil {
		return imageGenerationResult{Path: outputPath, Err: fmt.Errorf("copy generated image: %v %v %v", copyErr, closeOutputErr, closeInputErr)}
	}
	return publishGeneratedImage(stagePath, outputPath)
}

func publishGeneratedImage(stagePath, outputPath string) imageGenerationResult {
	contents, err := os.ReadFile(stagePath)
	if err != nil || len(contents) == 0 {
		return imageGenerationResult{Path: outputPath, Err: fmt.Errorf("generated image is empty: %v", err)}
	}
	if info, statErr := os.Lstat(outputPath); statErr == nil {
		if !info.Mode().IsRegular() {
			return imageGenerationResult{Path: outputPath, Err: fmt.Errorf("output is not a regular file")}
		}
		err = atomicfile.ReplaceExisting(outputPath, contents)
	} else if os.IsNotExist(statErr) {
		err = atomicfile.CreateExclusive(outputPath, contents, 0o600)
	} else {
		err = statErr
	}
	return imageGenerationResult{Path: outputPath, Err: err}
}

func nonemptyRegular(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular() && info.Size() > 0
}

func imageSignalContext() (context.Context, func(), func() int) {
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	observed := make(chan int, 1)
	signal.Notify(signals, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case signalValue := <-signals:
			status := 130
			switch signalValue {
			case syscall.SIGHUP:
				status = 129
			case syscall.SIGTERM:
				status = 143
			}
			observed <- status
			cancel()
		case <-ctx.Done():
		}
	}()
	status := 0
	return ctx, func() { signal.Stop(signals); cancel() }, func() int {
		select {
		case status = <-observed:
		default:
		}
		return status
	}
}
