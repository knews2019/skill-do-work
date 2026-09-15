package publication

import (
	"fmt"
	"os"
	"os/exec"
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
				runGitFixture(t, repositoryRoot, "add", "--", destinationPath)
				return nil
			}
			t.Cleanup(func() { beforePublicationMutation = previous })
			result := ApplyPlan(t.Context(), plan, false, false)
			staged, err := exec.Command("git", "-C", repositoryRoot, "diff", "--cached", "--name-only").Output()
			if err != nil || string(staged) != destinationPath+"\n" {
				t.Errorf("foreign staged entry changed: staged=%q err=%v", staged, err)
			}
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
	for _, sourceState := range []string{"create", "clean move", "dirty move", "dirty rewritten move"} {
		mutationKind := MutationMove
		if sourceState == "create" {
			mutationKind = MutationCreate
		}
		for _, failureStage := range []string{"earlier object replaced", "published destination replaced"} {
			t.Run(sourceState+"/"+failureStage, func(t *testing.T) {
				repositoryRoot := initializedGitRepository(t)
				const firstPath = "a-first.txt"
				const destinationPath = "b-published.txt"
				const sourcePath = "source.txt"
				publishedBytes := []byte("published bytes\n")
				sourceBytes := publishedBytes
				var dirtyPaths []string
				mutation := PlannedMutation{Kind: mutationKind, Path: destinationPath, Contents: publishedBytes, Mode: 0o644}
				if mutationKind == MutationMove {
					writeFixture(t, repositoryRoot, sourcePath, publishedBytes, 0o644)
					runGitFixture(t, repositoryRoot, "add", sourcePath)
					runGitFixture(t, repositoryRoot, "commit", "-qm", "move source")
					if strings.HasPrefix(sourceState, "dirty") {
						sourceBytes = []byte("dirty preimage\n")
						dirtyPaths = []string{sourcePath}
						writeFixture(t, repositoryRoot, sourcePath, sourceBytes, 0o644)
						if sourceState == "dirty move" {
							publishedBytes = sourceBytes
						}
					}
					mutation = PlannedMutation{Kind: MutationMove, Path: sourcePath, DestinationPath: destinationPath, ExpectedBytes: sourceBytes, Contents: publishedBytes}
				}
				plan := finalizePlan(PublicationPlan{Operation: OperationRelease, RepositoryRoot: repositoryRoot, ExistingDirtyTargetPaths: dirtyPaths, Mutations: []PlannedMutation{
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
					// The recorded file identity must not adopt a replacement directory.
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
					if !strings.Contains(rollbackErrors, "created target changed after created-object capture; preserved replacement: "+destinationPath) {
						t.Errorf("rollback did not identify the destination replacement: %s", rollbackErrors)
					}
				}
				if mutationKind == MutationMove {
					if contents, err := os.ReadFile(filepath.Join(repositoryRoot, sourcePath)); err != nil || string(contents) != string(sourceBytes) {
						t.Errorf("move source was not restored: contents=%q, err=%v", contents, err)
					}
				}
			})
		}
	}
}

func TestPublicationRecordingPreservesForeignRegularReplacements(t *testing.T) {
	for _, kind := range []MutationKind{MutationCreate, MutationMove} {
		for _, foreignContents := range []string{"our publication\n", "foreign replacement\n"} {
			t.Run(string(kind)+"/"+strings.TrimSpace(foreignContents), func(t *testing.T) {
				root := initializedGitRepository(t)
				published := []byte("our publication\n")
				mutation := PlannedMutation{Kind: kind, Path: "destination.txt", Contents: published, Mode: 0o644}
				if kind == MutationMove {
					writeFixture(t, root, "source.txt", published, 0o644)
					runGitFixture(t, root, "add", "source.txt")
					runGitFixture(t, root, "commit", "-qm", "source")
					mutation = PlannedMutation{Kind: kind, Path: "source.txt", DestinationPath: "destination.txt", ExpectedBytes: published}
				}
				plan := finalizePlan(PublicationPlan{Operation: OperationRelease, RepositoryRoot: root, Mutations: []PlannedMutation{mutation}})
				before, after := beforePublicationRecording, afterPublicationMutation
				t.Cleanup(func() { beforePublicationRecording, afterPublicationMutation = before, after })
				beforePublicationRecording = func(_ int, _ PlannedMutation) {
					// Keep the original inode alive so even equal bytes are a different object.
					if err := os.Rename(filepath.Join(root, "destination.txt"), filepath.Join(root, "held-original.txt")); err != nil {
						t.Fatal(err)
					}
					writeFixture(t, root, "destination.txt", []byte(foreignContents), 0o644)
				}
				if foreignContents == string(published) {
					afterPublicationMutation = func(_ int, _ PlannedMutation) error { return fmt.Errorf("later mutation failed") }
				}
				result := ApplyPlan(t.Context(), plan, false, false)
				contents, err := os.ReadFile(filepath.Join(root, "destination.txt"))
				if err != nil || string(contents) != foreignContents {
					t.Errorf("foreign replacement lost: contents=%q err=%v", contents, err)
				}
				if result.Outcome != resultmodel.OutcomeRisk || result.Rollback.Status != resultmodel.RollbackIncomplete {
					t.Errorf("foreign replacement reported restored: outcome=%s rollback=%+v", result.Outcome, result.Rollback)
				}
				if kind == MutationMove {
					contents, err := os.ReadFile(filepath.Join(root, "source.txt"))
					if err != nil || string(contents) != string(published) {
						t.Errorf("move source not restored: contents=%q err=%v", contents, err)
					}
				}
			})
		}
	}
}
