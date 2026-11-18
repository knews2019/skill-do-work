package toolboxcommands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"
)

type workedReportCard struct {
	Name, Kind, Subject, Report, Image string
	Width, Height                      int
}

var workedReportName = regexp.MustCompile(`^[a-z]+(-[a-z]+){1,2}$`)

func readWorkedReports(repository *os.Root) ([]workedReportCard, error) {
	record, err := fs.ReadFile(repository.FS(), reportGalleryFolder + "/preview-manifest.json")
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var manifest struct {
		Width, Height int
		Cards         []workedReportCard
	}
	if err := json.Unmarshal(record, &manifest); err != nil {
		return nil, err
	}
	if manifest.Width <= 0 || manifest.Height <= 0 || len(manifest.Cards) != 10 {
		return nil, errors.New("preview manifest requires positive dimensions and ten worked reports")
	}
	seenNames, seenReports := map[string]bool{}, map[string]bool{}
	kinds := map[string]int{}
	for index := range manifest.Cards {
		card := &manifest.Cards[index]
		if len(card.Name) > 24 || !workedReportName.MatchString(card.Name) || seenNames[card.Name] {
			return nil, fmt.Errorf("invalid or duplicate worked report name %q", card.Name)
		}
		if !strings.HasPrefix(card.Report, "ai-reports/") || path.Clean(card.Report) != card.Report ||
			path.Base(card.Report) != "index.html" || path.Dir(card.Report) == "ai-reports" || seenReports[card.Report] {
			return nil, fmt.Errorf("invalid or duplicate worked report path %q", card.Report)
		}
		if card.Image != reportGalleryFolder+"/previews/"+card.Name+".png" || strings.TrimSpace(card.Subject) == "" {
			return nil, fmt.Errorf("invalid preview or subject for %q", card.Name)
		}
		seenNames[card.Name], seenReports[card.Report] = true, true
		kinds[card.Kind]++
		card.Report = reportGalleryFileLink(repository, card.Report)
		card.Image = reportGalleryFileLink(repository, card.Image)
		if card.Report == "" || card.Image == "" {
			return nil, fmt.Errorf("missing source or screenshot for %q", card.Name)
		}
		card.Width, card.Height = manifest.Width, manifest.Height
	}
	for _, kind := range []string{"Completed work", "Proposal", "Root cause", "Architecture"} {
		if kinds[kind] < 2 {
			return nil, fmt.Errorf("preview manifest needs two %s reports", kind)
		}
	}
	return manifest.Cards, nil
}
