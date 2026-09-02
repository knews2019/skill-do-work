package nextselection

import (
	"fmt"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func resolveTargets(snapshot *repositorymodel.RepositorySnapshot, graph *dependencygraph.DependencyGraph, options SelectionOptions) ([]selectionCandidate, []resultmodel.SelectionExclusion) {
	if len(options.TargetTokens) == 0 {
		candidates := queueCandidates(snapshot, ProvenanceDefault)
		if options.SimpleOnly {
			for index := range candidates {
				candidates[index].Provenance = ProvenanceSimple
			}
		}
		return candidates, nil
	}

	requestByNumber := map[int]*repositorymodel.RequestFile{}
	for _, requestFile := range snapshot.RequestFiles {
		if requestFile.TreeSection != "queue" {
			continue
		}
		if number, valid := numericID(requestID(requestFile), "REQ-"); valid {
			requestByNumber[number] = requestFile
		}
	}
	ordered := []selectionCandidate{}
	indexByID := map[string]int{}
	exclusions := []resultmodel.SelectionExclusion{}
	explicitIDs := map[string]bool{}
	for _, token := range options.TargetTokens {
		if !hasIDPrefix(token, "REQ-") {
			continue
		}
		number, _ := numericID(token, "REQ-")
		if requestFile := requestByNumber[number]; requestFile != nil {
			explicitIDs[requestID(requestFile)] = true
		}
	}
	for _, token := range options.TargetTokens {
		switch {
		case hasIDPrefix(token, "REQ-"):
			number, _ := numericID(token, "REQ-")
			requestFile := requestByNumber[number]
			canonical, _ := canonicalToken(token, "REQ-")
			if requestFile == nil {
				exclusions = append(exclusions, targetNotFound(canonical, ProvenanceExplicit))
				continue
			}
			identifier := requestID(requestFile)
			if _, found := indexByID[identifier]; found {
				continue
			}
			indexByID[identifier] = len(ordered)
			ordered = append(ordered, selectionCandidate{RequestFile: requestFile, RequestID: identifier, Provenance: ProvenanceExplicit, SourceToken: canonical, Priority: selectionPriority(requestFile)})
		case hasIDPrefix(token, "UR-"):
			userRequestNumber, _ := numericID(token, "UR-")
			canonical, _ := canonicalToken(token, "UR-")
			members := []*repositorymodel.RequestFile{}
			for _, requestFile := range snapshot.RequestFiles {
				if requestFile.TreeSection != "queue" {
					continue
				}
				memberNumber, valid := numericID(requestFile.TypedRecord.UserRequestID, "UR-")
				if valid && memberNumber == userRequestNumber {
					members = append(members, requestFile)
				}
			}
			sort.SliceStable(members, func(leftIndex, rightIndex int) bool {
				leftPriority := priorityRank(selectionPriority(members[leftIndex]))
				rightPriority := priorityRank(selectionPriority(members[rightIndex]))
				if leftPriority != rightPriority {
					return leftPriority < rightPriority
				}
				leftID := requestID(members[leftIndex])
				rightID := requestID(members[rightIndex])
				leftDepth := graphDepth(graph, leftID)
				rightDepth := graphDepth(graph, rightID)
				if leftDepth != rightDepth {
					return leftDepth < rightDepth
				}
				return requestIDLess(leftID, rightID)
			})
			if len(members) == 0 {
				exclusions = append(exclusions, targetNotFound(canonical, ProvenanceUserRequest))
			}
			for _, requestFile := range members {
				identifier := requestID(requestFile)
				if explicitIDs[identifier] {
					continue
				}
				if _, found := indexByID[identifier]; found {
					continue
				}
				indexByID[identifier] = len(ordered)
				ordered = append(ordered, selectionCandidate{RequestFile: requestFile, RequestID: identifier, Provenance: ProvenanceUserRequest, SourceToken: canonical, Priority: selectionPriority(requestFile)})
			}
		}
	}
	return ordered, exclusions
}

func queueCandidates(snapshot *repositorymodel.RepositorySnapshot, provenance string) []selectionCandidate {
	candidates := []selectionCandidate{}
	for _, requestFile := range snapshot.RequestFiles {
		if requestFile.TreeSection == "queue" {
			candidates = append(candidates, selectionCandidate{RequestFile: requestFile, RequestID: requestID(requestFile), Provenance: provenance, Priority: selectionPriority(requestFile)})
		}
	}
	sort.SliceStable(candidates, func(leftIndex, rightIndex int) bool {
		return requestIDLess(candidates[leftIndex].RequestID, candidates[rightIndex].RequestID)
	})
	return candidates
}

func selectionPriority(requestFile *repositorymodel.RequestFile) string {
	if requestFile != nil && requestFile.TypedRecord.RepositoryGateRepairValue == "true" {
		return PriorityRepositoryGateRepair
	}
	if requestFile != nil && requestFile.TypedRecord.GateDeferredValue == "true" {
		return PriorityDeferredParent
	}
	return PriorityOrdinary
}

func hasIDPrefix(token, prefix string) bool {
	_, valid := numericID(token, prefix)
	return valid
}

func graphDepth(graph *dependencygraph.DependencyGraph, requestID string) int {
	if graph == nil || graph.NodesByID[requestID] == nil {
		return -1
	}
	return graph.NodesByID[requestID].DependencyDepth
}

func targetNotFound(identifier, provenance string) resultmodel.SelectionExclusion {
	nextArgv := []string{"do-work", "roadmap"}
	return resultmodel.SelectionExclusion{
		RequestID: identifier, Provenance: provenance, Code: "TARGET-NOT-FOUND",
		Reason:   fmt.Sprintf("%s did not resolve to a queued request", identifier),
		NextArgv: nextArgv, NextJustRecipe: justRecipeFor(nextArgv),
		VerificationArgv: []string{"do-work-cli", "--format", "json", "next", identifier},
	}
}

func validateTargetTokens(tokens []string) error {
	for _, token := range tokens {
		if hasIDPrefix(token, "REQ-") || hasIDPrefix(token, "UR-") {
			continue
		}
		return fmt.Errorf("unrecognized argument %q; expected REQ-NNN or UR-NNN", strings.TrimSpace(token))
	}
	return nil
}
