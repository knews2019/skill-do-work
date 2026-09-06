// Package dependencygraph derives deterministic dependency evidence from one
// repository snapshot without rescanning the filesystem.
package dependencygraph

import (
	"fmt"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/schemanormalization"
)

// DependencyEdge is one de-duplicated dependency relation.
type DependencyEdge struct {
	RequestID    string
	DependencyID string
}

// DependencyNode is the complete dependency evidence for one request id.
type DependencyNode struct {
	RequestID             string
	RequestStatus         string
	ImplementationCommit  string
	DependencyIDs         []string
	DependentIDs          []string
	UnmetDependencies     []string
	MissingTargets        []string
	AmbiguousTargets      []string
	DependenciesSatisfied bool
	IsReady               bool
	DependencyDepth       int
	IsCyclic              bool
	IsAmbiguous           bool
}

// DependencyCycle is one strongly connected request-id component in stable order.
type DependencyCycle struct {
	RequestIDs []string
}

// DependencyGraph is the stable dependency projection of a repository snapshot.
type DependencyGraph struct {
	NodesByID        map[string]*DependencyNode
	OrderedIDs       []string
	DependencyEdges  []DependencyEdge
	DependencyCycles []DependencyCycle
	WarningMessages  []string
}

// BuildGraph constructs dependency evidence exclusively from snapshot.
func BuildGraph(snapshot *repositorymodel.RepositorySnapshot) *DependencyGraph {
	graph := &DependencyGraph{NodesByID: map[string]*DependencyNode{}}
	if snapshot == nil {
		return graph
	}
	ambiguousIDs := map[string]bool{}
	for _, collision := range snapshot.CollisionEntries {
		ambiguousIDs[collision.RequestID] = true
	}
	for requestID, requestFiles := range snapshot.RequestsByID {
		if len(requestFiles) > 1 {
			ambiguousIDs[requestID] = true
		}
	}
	for _, requestFile := range snapshot.RequestFiles {
		requestID := requestFile.TypedRecord.RequestID
		if requestID == "" {
			requestID = requestFile.FilenameID
		}
		if requestID == "" {
			continue
		}
		if _, exists := graph.NodesByID[requestID]; exists {
			continue
		}
		isAmbiguous := ambiguousIDs[requestID]
		requestStatus := requestFile.TypedRecord.RequestStatus
		if isAmbiguous {
			requestStatus = ""
			graph.WarningMessages = append(graph.WarningMessages, fmt.Sprintf("%s has multiple repository records; dependency state is unresolved", requestID))
		}
		graph.NodesByID[requestID] = &DependencyNode{
			RequestID: requestID, RequestStatus: requestStatus,
			ImplementationCommit: requestFile.TypedRecord.ImplementationCommit,
			DependencyDepth:      -1, IsAmbiguous: isAmbiguous,
		}
		graph.OrderedIDs = append(graph.OrderedIDs, requestID)
	}
	sortRequestIDs(graph.OrderedIDs)

	for _, requestID := range graph.OrderedIDs {
		requestFile := firstRequestFile(snapshot, requestID)
		if requestFile == nil {
			continue
		}
		node := graph.NodesByID[requestID]
		if node.IsAmbiguous {
			continue
		}
		node.DependencyIDs = uniqueSortedIDs(requestFile.TypedRecord.DependsOn)
		for _, dependencyID := range node.DependencyIDs {
			graph.DependencyEdges = append(graph.DependencyEdges, DependencyEdge{RequestID: requestID, DependencyID: dependencyID})
			if dependencyNode := graph.NodesByID[dependencyID]; dependencyNode != nil {
				dependencyNode.DependentIDs = append(dependencyNode.DependentIDs, requestID)
			}
		}
	}
	for _, node := range graph.NodesByID {
		node.DependentIDs = uniqueSortedIDs(node.DependentIDs)
	}
	sort.Slice(graph.DependencyEdges, func(leftIndex, rightIndex int) bool {
		if graph.DependencyEdges[leftIndex].RequestID == graph.DependencyEdges[rightIndex].RequestID {
			return repositorymodel.RequestIDLess(graph.DependencyEdges[leftIndex].DependencyID, graph.DependencyEdges[rightIndex].DependencyID)
		}
		return repositorymodel.RequestIDLess(graph.DependencyEdges[leftIndex].RequestID, graph.DependencyEdges[rightIndex].RequestID)
	})

	detectCycles(graph)
	for _, requestID := range graph.OrderedIDs {
		node := graph.NodesByID[requestID]
		for _, dependencyID := range node.DependencyIDs {
			dependencyNode := graph.NodesByID[dependencyID]
			duplicateSatisfied := duplicateStatusesSatisfied(snapshot, dependencyID)
			switch {
			case ambiguousIDs[dependencyID] && !duplicateSatisfied:
				node.AmbiguousTargets = append(node.AmbiguousTargets, dependencyID)
				node.UnmetDependencies = append(node.UnmetDependencies, dependencyID)
				graph.WarningMessages = append(graph.WarningMessages, fmt.Sprintf("%s depends on ambiguous target %s", requestID, dependencyID))
			case dependencyNode == nil:
				node.MissingTargets = append(node.MissingTargets, dependencyID)
				node.UnmetDependencies = append(node.UnmetDependencies, dependencyID)
				graph.WarningMessages = append(graph.WarningMessages, fmt.Sprintf("%s depends on missing target %s", requestID, dependencyID))
			case !duplicateSatisfied && !schemanormalization.DependencySourceReady(dependencyNode.RequestStatus, dependencyNode.ImplementationCommit):
				node.UnmetDependencies = append(node.UnmetDependencies, dependencyID)
			}
		}
		node.DependenciesSatisfied = len(node.UnmetDependencies) == 0 && !node.IsCyclic && !node.IsAmbiguous
		node.IsReady = node.RequestStatus == "pending" && node.DependenciesSatisfied
	}

	depthMemo := map[string]int{}
	for _, requestID := range graph.OrderedIDs {
		graph.NodesByID[requestID].DependencyDepth = dependencyDepth(graph, ambiguousIDs, requestID, depthMemo, map[string]bool{})
	}
	return graph
}

// duplicateStatusesSatisfied resolves only exact duplicate records. Filename/frontmatter
// collisions remain ambiguous, while duplicate copies satisfy dependents only when every
// discovered record exposes source-ready dependency evidence.
func duplicateStatusesSatisfied(snapshot *repositorymodel.RepositorySnapshot, requestID string) bool {
	requestFiles := snapshot.RequestsByID[requestID]
	if len(requestFiles) < 2 {
		return false
	}
	for _, collision := range snapshot.CollisionEntries {
		if collision.RequestID == requestID && len(collision.ClaimPaths) != len(requestFiles) {
			return false
		}
	}
	for _, requestFile := range requestFiles {
		if requestFile.FilenameID != requestID || requestFile.TypedRecord.RequestID != requestID {
			return false
		}
		if !schemanormalization.DependencySourceReady(requestFile.TypedRecord.RequestStatus, requestFile.TypedRecord.ImplementationCommit) {
			return false
		}
	}
	return true
}

func firstRequestFile(snapshot *repositorymodel.RepositorySnapshot, requestID string) *repositorymodel.RequestFile {
	requestFiles := snapshot.RequestsByID[requestID]
	if len(requestFiles) > 0 {
		return requestFiles[0]
	}
	for _, requestFile := range snapshot.RequestFiles {
		candidateID := requestFile.TypedRecord.RequestID
		if candidateID == "" {
			candidateID = requestFile.FilenameID
		}
		if candidateID == requestID {
			return requestFile
		}
	}
	return nil
}

func detectCycles(graph *DependencyGraph) {
	nextIndex := 0
	indexes := map[string]int{}
	lowLinks := map[string]int{}
	onStack := map[string]bool{}
	var stack []string
	var visit func(string)
	visit = func(requestID string) {
		indexes[requestID] = nextIndex
		lowLinks[requestID] = nextIndex
		nextIndex++
		stack = append(stack, requestID)
		onStack[requestID] = true
		for _, dependencyID := range graph.NodesByID[requestID].DependencyIDs {
			if graph.NodesByID[dependencyID] == nil {
				continue
			}
			dependencyIndex, visited := indexes[dependencyID]
			if !visited {
				visit(dependencyID)
				if lowLinks[dependencyID] < lowLinks[requestID] {
					lowLinks[requestID] = lowLinks[dependencyID]
				}
			} else if onStack[dependencyID] && dependencyIndex < lowLinks[requestID] {
				lowLinks[requestID] = dependencyIndex
			}
		}
		if lowLinks[requestID] != indexes[requestID] {
			return
		}
		var component []string
		for len(stack) > 0 {
			lastIndex := len(stack) - 1
			componentID := stack[lastIndex]
			stack = stack[:lastIndex]
			onStack[componentID] = false
			component = append(component, componentID)
			if componentID == requestID {
				break
			}
		}
		selfCycle := len(component) == 1 && containsID(graph.NodesByID[component[0]].DependencyIDs, component[0])
		if len(component) <= 1 && !selfCycle {
			return
		}
		sortRequestIDs(component)
		for _, componentID := range component {
			graph.NodesByID[componentID].IsCyclic = true
		}
		graph.DependencyCycles = append(graph.DependencyCycles, DependencyCycle{RequestIDs: component})
		graph.WarningMessages = append(graph.WarningMessages, fmt.Sprintf("dependency cycle: %s", strings.Join(component, " -> ")))
	}
	for _, requestID := range graph.OrderedIDs {
		if _, visited := indexes[requestID]; !visited {
			visit(requestID)
		}
	}
	sort.Slice(graph.DependencyCycles, func(leftIndex, rightIndex int) bool {
		return repositorymodel.RequestIDLess(graph.DependencyCycles[leftIndex].RequestIDs[0], graph.DependencyCycles[rightIndex].RequestIDs[0])
	})
}

func dependencyDepth(graph *DependencyGraph, ambiguousIDs map[string]bool, requestID string, memo map[string]int, visiting map[string]bool) int {
	if depth, found := memo[requestID]; found {
		return depth
	}
	node := graph.NodesByID[requestID]
	if node == nil || node.IsCyclic || ambiguousIDs[requestID] || visiting[requestID] {
		memo[requestID] = -1
		return -1
	}
	visiting[requestID] = true
	maximumDependencyDepth := -1
	for _, dependencyID := range node.DependencyIDs {
		if graph.NodesByID[dependencyID] == nil || ambiguousIDs[dependencyID] {
			delete(visiting, requestID)
			memo[requestID] = -1
			return -1
		}
		dependencyNodeDepth := dependencyDepth(graph, ambiguousIDs, dependencyID, memo, visiting)
		if dependencyNodeDepth < 0 {
			delete(visiting, requestID)
			memo[requestID] = -1
			return -1
		}
		if dependencyNodeDepth > maximumDependencyDepth {
			maximumDependencyDepth = dependencyNodeDepth
		}
	}
	delete(visiting, requestID)
	depth := maximumDependencyDepth + 1
	memo[requestID] = depth
	return depth
}

func uniqueSortedIDs(requestIDs []string) []string {
	seen := map[string]bool{}
	uniqueIDs := make([]string, 0, len(requestIDs))
	for _, requestID := range requestIDs {
		trimmedID := strings.TrimSpace(requestID)
		if trimmedID == "" || seen[trimmedID] {
			continue
		}
		seen[trimmedID] = true
		uniqueIDs = append(uniqueIDs, trimmedID)
	}
	sortRequestIDs(uniqueIDs)
	return uniqueIDs
}

func sortRequestIDs(requestIDs []string) {
	sort.Slice(requestIDs, func(leftIndex, rightIndex int) bool {
		return repositorymodel.RequestIDLess(requestIDs[leftIndex], requestIDs[rightIndex])
	})
}

func containsID(requestIDs []string, targetID string) bool {
	for _, requestID := range requestIDs {
		if requestID == targetID {
			return true
		}
	}
	return false
}
