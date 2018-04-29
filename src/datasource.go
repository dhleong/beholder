package beholder

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

const compendiumURL = "https://raw.githubusercontent.com/storskegg/DnDAppFiles/master/Compendiums/Full%20Compendium.xml"

// DataSource abstracts fetching and loading Entities
type DataSource interface {
	GetEntities() ([]Entity, error)
}

// NewDataSource creates a new default DataSource
func NewDataSource() (DataSource, error) {
	localPath, err := homedir.Expand("~/.config/beholder/data.xml")
	if err != nil {
		return nil, err
	}

	ds := &networkDataSource{
		compendiumURL: compendiumURL,
		localPath:     localPath,
		delegate: &diskDataSource{
			localPath: localPath,
		},
	}

	return ds, nil
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
		dir := path.Dir(d.localPath)
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
