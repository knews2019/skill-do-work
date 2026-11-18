package toolboxcommands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const (
	reportGalleryGeneratedPrefix = "report:"
	reportGalleryDesignRecord    = "gallery-design.md"
)

type reportGalleryGeneratedSource struct {
	bundle   reportBundle
	identity string
}

type reportGalleryGeneratedDesign struct {
	recordPath string
	title      string
	useWhen    string
	borrow     []string
}

func reportGalleryGeneratedSources(repository *os.Root, bundles []reportBundle) []reportGalleryGeneratedSource {
	if info, err := repository.Lstat(reportsFolder); err != nil || !info.IsDir() {
		return nil
	}
	var sources []reportGalleryGeneratedSource
	basenameCount := map[string]int{}
	for _, bundle := range bundles {
		if bundle.Entry == "" || path.Base(bundle.Entry) == reportGalleryDesignRecord ||
			strings.IndexFunc(bundle.Entry, unicode.IsControl) >= 0 {
			continue
		}
		basenameCount[path.Base(bundle.Path)]++
		sources = append(sources, reportGalleryGeneratedSource{bundle: bundle})
	}
	for index := range sources {
		source := &sources[index]
		basename := path.Base(source.bundle.Path)
		source.identity = reportGalleryGeneratedPrefix + basename
		// A unique basename survives the organizer moving a bundle into a month.
		// Qualify collisions and names ending in whitespace so every starter round-trips.
		if basenameCount[basename] > 1 || strings.TrimSpace(basename) != basename {
			source.identity = source.qualifiedIdentity()
		}
	}
	return sources
}

func (source reportGalleryGeneratedSource) qualifiedIdentity() string {
	return reportGalleryGeneratedPrefix + source.bundle.Path + "/"
}

func (source reportGalleryGeneratedSource) readDesign(repository *os.Root) (reportGalleryGeneratedDesign, error) {
	record := path.Join(source.bundle.Path, reportGalleryDesignRecord)
	document, err := readReportGalleryRecord(repository, filepath.FromSlash(record))
	if errors.Is(err, fs.ErrNotExist) {
		// A legacy bundle's own primary page is the structural reference.
		design := reportGalleryGeneratedDesign{
			recordPath: source.bundle.Entry,
			title:      source.bundle.Title,
			useWhen:    "Adapt the layout and explanatory moves of this report.",
		}
		info, err := repository.Lstat(filepath.FromSlash(design.recordPath))
		if err == nil && !info.Mode().IsRegular() {
			err = errors.New("the report design source must be a regular file")
		}
		return design, err
	}
	design := reportGalleryGeneratedDesign{recordPath: record}
	if err != nil {
		return design, err
	}
	design.title, err = reportGalleryScalar(document, "title")
	if err != nil {
		return design, err
	}
	design.useWhen, err = reportGalleryScalar(document, "use_when")
	if err != nil {
		return design, err
	}
	design.borrow, err = readReportGalleryBorrow(document.BodyBytes())
	return design, err
}

func resolveReportGalleryGenerated(repositoryRoot string, repository *os.Root, identity string) resultmodel.CommandResult {
	info, err := repository.Lstat(reportsFolder)
	if errors.Is(err, fs.ErrNotExist) {
		return reportGalleryFailure(reportsFolder, fmt.Errorf("no generated report design matches %q", identity))
	}
	if err != nil {
		return reportGalleryFailure(reportsFolder, err)
	}
	if !info.IsDir() {
		return reportGalleryFailure(reportsFolder, errors.New("ai-reports must be a directory"))
	}
	bundles, _, err := reportIndexWalk(repositoryRoot)
	if err != nil {
		return reportGalleryFailure(reportsFolder, err)
	}
	var matches []reportGalleryGeneratedSource
	for _, source := range reportGalleryGeneratedSources(repository, bundles) {
		if identity == source.qualifiedIdentity() || identity == reportGalleryGeneratedPrefix+path.Base(source.bundle.Path) {
			matches = append(matches, source)
		}
	}
	switch len(matches) {
	case 0:
		return reportGalleryFailure(reportsFolder, fmt.Errorf("no generated report design matches %q", identity))
	case 1:
		source := matches[0]
		design, err := source.readDesign(repository)
		if err != nil {
			return reportGalleryFailure(design.recordPath, err)
		}
		return exactOutputResult(fmt.Sprintf("%s\t%s\t%s\n", source.identity, design.useWhen, design.recordPath), nil)
	default:
		options := make([]string, 0, len(matches))
		for _, source := range matches {
			options = append(options, source.qualifiedIdentity())
		}
		return reportGalleryFailure(reportsFolder, fmt.Errorf("generated report design %q is ambiguous. Use one of these identities: %s", identity, strings.Join(options, ", ")))
	}
}
