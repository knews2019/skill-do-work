#!/bin/sh
set -eu
cd "$(git rev-parse --show-toplevel)/skills/do-work/tools/do-work-cli"
run_named_tests() {
  package="$1"
  shift
  names=$(IFS='|'; echo "$*")
  output=$(go test -count=1 -v -run "^($names)$" "$package" 2>&1) || { echo "$output" >&2; return 1; }
  for test_name in "$@"; do
    echo "$output" | grep -F -- "--- PASS: $test_name ("
  done
}
run_named_tests ./internal/toolboxcommands TestHandlersRegisterCanonicalToolboxCommands TestAIReportGalleryListUsesProjectOverrideAndSkipsBrokenRecords TestAIReportGalleryPageLinksLatestRealReportsAndAssets TestAIReportGalleryEmptyListWritesNothingAndSourceCanCreatePage TestAIReportGalleryRefusesHandMadeOrLinkedOutput TestAIReportGalleryPublicationPreservesAnEditAfterOwnershipCheck TestAIReportIndexFindListsMatchesNewestFirstIncludingSuperseded TestAIReportIndexRefusesHandMadeCatalog
run_named_tests ./internal/suiteinstall TestReportGalleryRebuildsFromTheWrittenSuite TestReportGalleryWriteProblemsWarnOnceWithoutRecoveringTheSuite TestAnAvailableUpdateInstallsAndVerifiesTheNewVersion TestInstallAndUpdateRecommendOrganizeWhenReportsAreLoose
