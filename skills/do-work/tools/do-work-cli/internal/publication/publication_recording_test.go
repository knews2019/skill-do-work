package publication

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestPublicationCreationIntentPreservesForeignDestinationAfterPreflight(t *testing.T) {
	for _, mutationKind := range []MutationKind{MutationCreate, MutationMove} {
		t.Run(string(mutationKind), func(t *testing.T) {
			repositoryRoot := initializedGitRepository(t)
			const destinationPath = "destination.txt"
			const sourcePath = "source.txt"
			publishedBytes := []byte("our publication\n")
			mutation := PlannedMutation{Kind: mutationKind, Path: destinationPath, Contents: publishedBytes, Mode: 0o644}
			if mutationKind == MutationMove {
				writeFixture(t, repositoryRoot, sourcePath, publishedBytes, 0o644)
				runGitFixture(t, repositoryRoot, "add", sourcePath)
				runGitFixture(t, repositoryRoot, "commit", "-qm", "move source")
				mutation = PlannedMutation{Kind: MutationMove, Path: sourcePath, DestinationPath: destinationPath, ExpectedBytes: publishedBytes}
			}
			plan := finalizePlan(PublicationPlan{Operation: OperationRelease, RepositoryRoot: repositoryRoot, Mutations: []PlannedMutation{mutation}})
			previous := beforePublicationMutation
			beforePublicationMutation = func(_ int, _ PlannedMutation) error {
				writeFixture(t, repositoryRoot, destinationPath, []byte("foreign destination\n"), 0o644)
				return nil
			}
			t.Cleanup(func() { beforePublicationMutation = previous })
			result := ApplyPlan(t.Context(), plan, false, false)
			if contents, err := os.ReadFile(filepath.Join(repositoryRoot, destinationPath)); err != nil || string(contents) != "foreign destination\n" {
				t.Errorf("creation intent claimed and removed a foreign destination: contents=%q, err=%v", contents, err)
			}
			if result.Outcome != resultmodel.OutcomeRisk || result.Rollback.Status != resultmodel.RollbackIncomplete {
				t.Errorf("foreign destination did not report incomplete rollback: %#v", result)
			}
			if rollbackErrors := strings.Join(result.Rollback.Errors, "\n"); !strings.Contains(rollbackErrors, "created target was not identity-recorded; preserved object: "+destinationPath) {
				t.Errorf("rollback did not report unowned destination: %s", rollbackErrors)
			}
			if mutationKind == MutationMove {
				if contents, err := os.ReadFile(filepath.Join(repositoryRoot, sourcePath)); err != nil || string(contents) != string(publishedBytes) {
					t.Errorf("failed move did not preserve source: contents=%q, err=%v", contents, err)
				}
			}
		})
	}
}

func TestPublicationRecordingFailureKeepsCreatedTargetsVisibleToRollback(t *testing.T) {
	for _, mutationKind := range []MutationKind{MutationCreate, MutationMove} {
		for _, failureStage := range []string{"earlier object replaced", "published identity unavailable"} {
			t.Run(string(mutationKind)+"/"+failureStage, func(t *testing.T) {
				repositoryRoot := initializedGitRepository(t)
				const firstPath = "a-first.txt"
				const destinationPath = "b-published.txt"
				const sourcePath = "source.txt"
				publishedBytes := []byte("published bytes\n")
				mutation := PlannedMutation{Kind: mutationKind, Path: destinationPath, Contents: publishedBytes, Mode: 0o644}
				if mutationKind == MutationMove {
					writeFixture(t, repositoryRoot, sourcePath, publishedBytes, 0o644)
					runGitFixture(t, repositoryRoot, "add", sourcePath)
					runGitFixture(t, repositoryRoot, "commit", "-qm", "move source")
					mutation = PlannedMutation{Kind: MutationMove, Path: sourcePath, DestinationPath: destinationPath, ExpectedBytes: publishedBytes}
				}
				plan := finalizePlan(PublicationPlan{Operation: OperationRelease, RepositoryRoot: repositoryRoot, Mutations: []PlannedMutation{
					{Kind: MutationCreate, Path: firstPath, Contents: []byte("first publication\n"), Mode: 0o644},
					mutation,
				}})
				publicationObserved := false
				previous := beforePublicationRecording
				beforePublicationRecording = func(index int, _ PlannedMutation) {
					if index != 1 {
						return
					}
					publicationObserved = true
					publishedPath := filepath.Join(repositoryRoot, destinationPath)
					if contents, err := os.ReadFile(publishedPath); err != nil || string(contents) != string(publishedBytes) {
						t.Fatalf("recording seam ran before publication: contents=%q, err=%v", contents, err)
					}
					if mutationKind == MutationMove {
						if _, err := os.Stat(filepath.Join(repositoryRoot, sourcePath)); !os.IsNotExist(err) {
							t.Fatalf("recording seam ran before source removal: %v", err)
						}
					}
					if failureStage == "earlier object replaced" {
						// Revalidation must fail after capturing the new destination's identity.
						foreignPath := filepath.Join(repositoryRoot, firstPath)
						if err := os.Remove(foreignPath); err != nil {
							t.Fatal(err)
						}
						writeFixture(t, repositoryRoot, firstPath, []byte("foreign replacement\n"), 0o644)
						return
					}
					// A directory cannot be identity-recorded as the file just published.
					if err := os.Remove(publishedPath); err != nil {
						t.Fatal(err)
					}
					writeFixture(t, repositoryRoot, destinationPath+"/foreign.txt", []byte("foreign directory\n"), 0o644)
				}
				t.Cleanup(func() { beforePublicationRecording = previous })
				result := ApplyPlan(t.Context(), plan, false, false)
				if !publicationObserved {
					t.Fatal("post-publication recording failure was not exercised")
				}
				if result.Outcome != resultmodel.OutcomeRisk || result.Rollback.Status != resultmodel.RollbackIncomplete {
					t.Fatalf("recording failure did not report incomplete rollback: %#v", result)
				}
				rollbackErrors := strings.Join(result.Rollback.Errors, "\n")
				if failureStage == "earlier object replaced" {
					if contents, err := os.ReadFile(filepath.Join(repositoryRoot, firstPath)); err != nil || string(contents) != "foreign replacement\n" {
						t.Fatalf("foreign replacement was not preserved: contents=%q, err=%v", contents, err)
					}
					if !strings.Contains(rollbackErrors, "created target changed after created-object capture; preserved replacement: "+firstPath) {
						t.Errorf("rollback did not identify the foreign replacement: %s", rollbackErrors)
					}
					if _, err := os.Stat(filepath.Join(repositoryRoot, destinationPath)); !os.IsNotExist(err) {
						t.Errorf("recording failure left the published destination behind: %v", err)
					}
				} else {
					if contents, err := os.ReadFile(filepath.Join(repositoryRoot, destinationPath, "foreign.txt")); err != nil || string(contents) != "foreign directory\n" {
						t.Fatalf("unknown-ownership replacement was not preserved: contents=%q, err=%v", contents, err)
					}
					if !strings.Contains(rollbackErrors, "created target was not identity-recorded; preserved object: "+destinationPath) {
						t.Errorf("rollback did not identify unknown destination ownership: %s", rollbackErrors)
					}
				}
				if mutationKind == MutationMove {
					if contents, err := os.ReadFile(filepath.Join(repositoryRoot, sourcePath)); err != nil || string(contents) != string(publishedBytes) {
						t.Errorf("move source was not restored: contents=%q, err=%v", contents, err)
					}
				}
			})
		}
	}
}
