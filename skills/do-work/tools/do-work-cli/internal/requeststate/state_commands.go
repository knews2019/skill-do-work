package requeststate

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

var discoverRepository = repositorymodel.DiscoverRepository

func Handlers() map[string]commandruntime.CommandHandler {
	handlers := map[string]commandruntime.CommandHandler{}
	for _, transition := range []Transition{TransitionClaim, TransitionRecover, TransitionUnblock, TransitionComplete, TransitionFail, TransitionCancel} {
		selectedTransition := transition
		handlers[string(transition)] = func(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
			return handleStateCommand(executionContext, selectedTransition, arguments)
		}
	}
	return handlers
}

func handleStateCommand(executionContext commandruntime.ExecutionContext, transition Transition, arguments []string) resultmodel.CommandResult {
	options, parseError := parseStateOptions(transition, arguments)
	if parseError != nil {
		return commandFailure(executionContext.RepositoryRoot, transition, "STATE-USAGE", parseError.Error())
	}
	snapshot, discoveryError := discoverRepository(executionContext.RepositoryRoot)
	if discoveryError != nil {
		return commandFailure(executionContext.RepositoryRoot, transition, "STATE-DISCOVERY-FAILED", discoveryError.Error())
	}
	return ApplyPlan(context.Background(), BuildPlan(snapshot, dependencygraph.BuildGraph(snapshot), options))
}

func parseStateOptions(transition Transition, arguments []string) (StateOptions, error) {
	options := StateOptions{Transition: transition, Provenance: ProvenanceDefault, TerminalStatus: "completed"}
	if len(arguments) == 0 || strings.HasPrefix(arguments[0], "-") {
		return options, fmt.Errorf("%s requires one REQ-NNN id", transition)
	}
	options.RequestID = arguments[0]
	for index := 1; index < len(arguments); index++ {
		argument := arguments[index]
		switch argument {
		case "--dry-run":
			options.DryRun = true
		case "--commit":
			options.Commit = true
		case "--checkpoint-unlabeled":
			options.CheckpointUnlabeled = true
		case "--checkpoint-absent":
			options.CheckpointAbsent = true
		case "--assume-sole-writer":
			options.AssumeSoleWriter = true
		case "--record-commit-hash":
			options.RecordCommitHashOnly = true
		case "--unblock-required":
			options.UnblockRequired = true
		case "--confirmed":
			options.CancellationConfirmed = true
		case "--request-path", "--provenance", "--original-status", "--probe-status", "--source", "--terminal-status", "--implementation-hash", "--error", "--error-type", "--reason", "--reason-summary", "--dependent-disposition", "--writer", "--checkpoint-writer", "--at":
			index++
			if index >= len(arguments) {
				return options, fmt.Errorf("%s requires a value", argument)
			}
			value := arguments[index]
			switch argument {
			case "--request-path":
				options.RequestPath = value
			case "--provenance":
				options.Provenance = SelectionProvenance(value)
			case "--original-status":
				options.OriginalStatus = value
			case "--probe-status":
				options.ProbeStatus = resultmodel.SelectionProbeStatus(value)
			case "--source":
				options.UnblockSource = UnblockSource(value)
			case "--terminal-status":
				options.TerminalStatus = value
			case "--implementation-hash":
				options.ImplementationHash = value
			case "--error":
				options.FailureError = value
			case "--error-type":
				options.FailureType = value
			case "--reason":
				options.CancellationReason = value
			case "--reason-summary":
				options.CancellationSummary = value
			case "--dependent-disposition":
				options.DependentDisposition = value
			case "--writer":
				options.WriterLabel = value
			case "--checkpoint-writer":
				options.CheckpointWriter = value
			case "--at":
				parsed, err := time.Parse(time.RFC3339, value)
				if err != nil {
					return options, fmt.Errorf("--at requires RFC3339: %w", err)
				}
				options.Now = parsed
			}
		default:
			return options, fmt.Errorf("unknown %s option %q", transition, argument)
		}
	}
	if options.DryRun && options.Commit {
		return options, fmt.Errorf("--dry-run and --commit cannot be combined")
	}
	if options.Provenance != ProvenanceDefault && options.Provenance != ProvenanceExplicit && options.Provenance != ProvenanceURExpanded {
		return options, fmt.Errorf("unknown selection provenance %q", options.Provenance)
	}
	return options, nil
}
