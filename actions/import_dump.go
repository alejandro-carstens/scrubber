package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/options"
	"github.com/alejandro-carstens/scrubber/filesystem"
	"github.com/ivpusic/grpool"
)

type indexConfig struct {
	name     string
	settings *gabs.Container
	aliases  *gabs.Container
	mappings *gabs.Container
}

func (ic *indexConfig) transform() (string, error) {
	settings, err := ic.transformSettings()

	if err != nil {
		return "", err
	}

	mappings, err := ic.transformMappings()

	if err != nil {
		return "", err
	}

	aliases, err := ic.transformAliases()

	if err != nil {
		return "", err
	}

	return mapToString(map[string]interface{}{
		"settings": settings,
		"mappings": mappings,
		"aliases":  aliases,
	})
}

func (ic *indexConfig) transformSettings() (map[string]interface{}, error) {
	settings := map[string]map[string]interface{}{}

	if err := json.Unmarshal(ic.settings.Bytes(), &settings); err != nil {
		return nil, err
	}

	indexSettings, valid := settings["index"]

	if !valid {
		return map[string]interface{}{}, nil
	}

	delete(indexSettings, "creation_date")
	delete(indexSettings, "provided_name")
	delete(indexSettings, "uuid")
	delete(indexSettings, "version")

	return indexSettings, nil
}

func (ic *indexConfig) transformMappings() (map[string]interface{}, error) {
	mappings := map[string]map[string]interface{}{}

	if err := json.Unmarshal(ic.mappings.Bytes(), &mappings); err != nil {
		return nil, err
	}

	indexMappings, valid := mappings[ic.name]["mappings"]

	if !valid {
		return map[string]interface{}{}, nil
	}

	return indexMappings.(map[string]interface{}), nil
}

func (ic *indexConfig) transformAliases() (map[string]interface{}, error) {
	aliases, err := ic.aliases.S("Indices", ic.name, "Aliases").Children()

	if err != nil {
		return nil, err
	}

	indexAliases := map[string]interface{}{}

	for _, alias := range aliases {
		indexAlias := map[string]interface{}{}

		if err := json.Unmarshal(alias.Bytes(), &indexAlias); err != nil {
			return nil, err
		}

		delete(indexAlias, "IsWriteIndex")

		for _, name := range indexAlias {
			indexAliases[name.(string)] = map[string]interface{}{}
		}
	}

	return indexAliases, nil
}

type importDump struct {
	action
	options *options.ImportDumpOptions
}

// ApplyOptions implementation of the Actionable interface
func (id *importDump) ApplyOptions() Actionable {
	id.options = id.context.Options().(*options.ImportDumpOptions)

	id.indexer.SetOptions(&golastic.IndexOptions{Timeout: id.options.TimeoutInSeconds()})

	return id
}

// Perform implementation of the Actionable interface
func (id *importDump) Perform() Actionable {
	configs, err := id.getIndexConfigs()

	if err != nil {
		id.errorContainer.push(id.name, "_all", err)

		return id
	}

	for _, config := range configs {
		if err := id.importDump(config); err != nil {
			id.errorContainer.push(id.name, id.indexName(config.name), err)

			break
		}
	}

	return id
}

// ApplyFilters implementation of the Actionable interface
func (id *importDump) ApplyFilters() error {
	return nil
}

func (id *importDump) importDump(config *indexConfig) error {
	schema, err := config.transform()

	if err != nil {
		return err
	}

	if err := id.indexer.CreateIndex(id.indexName(config.name), schema); err != nil {
		return err
	}

	time.Sleep(1)

	fs, err := filesystem.Build(id.filesystemConfig())

	if err != nil {
		return err
	}

	dataFiles, err := fs.List(config.name)

	if err != nil {
		return err
	}

	for i, dataFile := range dataFiles {
		if strings.Contains(dataFile, "data") {
			dataFiles[i] = dataFiles[len(dataFiles)-1]

			dataFiles[len(dataFiles)-1] = ""

			dataFiles = dataFiles[:len(dataFiles)-1]
		}
	}

	pool := grpool.NewPool(id.options.Concurrency, len(dataFiles))

	defer pool.Release()

	pool.WaitCount(len(dataFiles))

	errs := []error{}

	for _, dataFile := range dataFiles {
		pool.JobQueue <- func() {
			defer pool.JobDone()

			if err := id.importData(config.name, dataFile); err != nil {
				id.errorContainer.push(id.name, id.indexName(config.name), err)

				errs = append(errs, err)
			}
		}
	}

	pool.WaitAll()

	if len(errs) > 0 {
		return errors.New("an error occurred while importing data")
	}

	return nil
}

func (id *importDump) importData(name string, dataFile string) error {
	// Implement data import
	return nil
}

func (id *importDump) getIndexConfigs() ([]*indexConfig, error) {
	fs, err := filesystem.Build(id.filesystemConfig())

	if err != nil {
		return nil, err
	}

	list, err := fs.List("")

	if err != nil {
		return nil, err
	}

	errs := []error{}
	properties := []*indexConfig{}

	var wg sync.WaitGroup
	var mux sync.Mutex

	wg.Add(len(list))

	for _, index := range list {
		go func() {
			defer wg.Done()

			contents, err := id.extractIndexConfig(index)

			if err != nil {
				errs = append(errs, err)

				return
			}

			mux.Lock()

			properties = append(properties, contents)

			mux.Unlock()
		}()
	}

	wg.Wait()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	return properties, nil
}

func (id *importDump) extractIndexConfig(index string) (*indexConfig, error) {
	fs, err := filesystem.Build(id.filesystemConfig())

	if err != nil {
		return nil, err
	}

	settings, err := fileToJSON(fs, filepath.Join(index, filepath.FromSlash("settings.json")))

	if err != nil {
		return nil, err
	}

	aliases, err := fileToJSON(fs, filepath.Join(index, filepath.FromSlash("aliases.json")))

	if err != nil {
		return nil, err
	}

	mappings, err := fileToJSON(fs, filepath.Join(index, filepath.FromSlash("mappings.json")))

	if err != nil {
		return nil, err
	}

	return &indexConfig{
		name:     index,
		settings: settings,
		aliases:  aliases,
		mappings: mappings,
	}, nil
}

func (id *importDump) filesystemConfig() filesystem.Configurable {
	if id.options.Repository == "gcs" {
		return &filesystem.GCS{
			Context:             id.ctx,
			Bucket:              id.options.Bucket,
			CredentialsFilePath: id.options.CredentialsFilePath,
			Directory:           id.options.Name,
		}
	}

	return &filesystem.Local{
		Path: filepath.Join(id.options.Path, filepath.FromSlash(id.options.Name)),
	}
}

func (id *importDump) indexName(index string) string {
	return fmt.Sprintf("import-%v-%v", id.options.Name, index)
}
