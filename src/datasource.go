package beholder

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

var compendiums = map[string]string{
	"data.xml":        "https://raw.githubusercontent.com/storskegg/DnDAppFiles/master/Compendiums/Full%20Compendium.xml",
	"xgte-spells.xml": "https://raw.githubusercontent.com/storskegg/DnDAppFiles/master/Spells/XGtE%20Spells.xml",
}

// DataSource abstracts fetching and loading Entities
type DataSource interface {
	GetEntities() ([]Entity, error)
}

// NewDataSource creates a new default DataSource
func NewDataSource() (DataSource, error) {
	sources := make([]DataSource, 0, 3+len(compendiums))

	for localName, url := range compendiums {
		sources = append(sources, newNetworkDataSource(url, localName))
	}

	sources = append(sources,
		ActionsDataSource,
		ConditionsDataSource,
		RuleDataSource,
	)

	return MergeDataSources(sources...), nil
}

func newNetworkDataSource(url, localName string) *networkDataSource {
	relativePath := fmt.Sprintf("~/.config/beholder/%s", localName)
	localPath, err := homedir.Expand(relativePath)
	if err != nil {
		return nil
	}

	return &networkDataSource{
		compendiumURL: url,
		localPath:     localPath,
		delegate: &diskDataSource{
			localPath: localPath,
		},
	}
}

type staticDataSource struct {
	entities []Entity
}

func (d *staticDataSource) GetEntities() ([]Entity, error) {
	return d.entities, nil
}

// NewStaticDataSource .
func NewStaticDataSource(entities []Entity) DataSource {
	return &staticDataSource{entities}
}

type networkDataSource struct {
	compendiumURL string
	localPath     string

	delegate DataSource
}

func (d *networkDataSource) GetEntities() ([]Entity, error) {

	// do we need to fetch?
	if _, err := os.Stat(d.localPath); err != nil {
		if !os.IsNotExist(err) {
			// some unexpected issue; report it:
			return nil, err
		}

		// yes, let's fetch
		dir := filepath.Dir(d.localPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}

		// create the local file
		out, err := os.Create(d.localPath)
		if err != nil {
			return nil, err
		}
		defer out.Close()

		// request the remote file
		resp, err := http.Get(d.compendiumURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// write the response body
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return nil, err
		}
	}

	return d.delegate.GetEntities()
}

type diskDataSource struct {
	localPath string
}

func (d *diskDataSource) GetEntities() ([]Entity, error) {
	f, err := os.Open(d.localPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ParseXML(bufio.NewReader(f))
}

type mergeDataSource struct {
	sources []DataSource
}

func (d *mergeDataSource) GetEntities() ([]Entity, error) {
	all := make([]Entity, 0, 0)

	for _, s := range d.sources {
		result, err := s.GetEntities()
		if err != nil {
			return nil, err
		}

		all = append(all, result...)
	}

	return all, nil
}

// MergeDataSources creates a DataSource that combines the
// results of all the provided DataSource instances
func MergeDataSources(sources ...DataSource) DataSource {
	return &mergeDataSource{
		sources: sources,
	}
}
