package actions

import (
	"encoding/json"
	"fmt"
	"path/filepath"
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

func (ic *indexConfig) format() (string, error) {
	settings, err := ic.formatSettings()

	if err != nil {
		return "", err
	}

	mappings, err := ic.formatMappings()

	if err != nil {
		return "", err
	}

	aliases, err := ic.formatAliases()

	if err != nil {
		return "", err
	}

	return mapToString(map[string]interface{}{
		"settings": settings,
		"mappings": mappings,
		"aliases":  aliases,
	})
}

func (ic *indexConfig) formatSettings() (map[string]interface{}, error) {
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

func (ic *indexConfig) formatMappings() (map[string]interface{}, error) {
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

func (ic *indexConfig) formatAliases() (map[string]interface{}, error) {
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

	pool := grpool.NewPool(id.options.Concurrency, len(configs))

	defer pool.Release()

	pool.WaitCount(len(configs))

	for _, config := range configs {
		pool.JobQueue <- func() {
			defer pool.JobDone()

			if err := id.importDump(config); err != nil {
				id.errorContainer.push(id.name, id.indexName(config.name), err)
			}
		}
	}

	pool.WaitAll()

	return id
}

// ApplyFilters implementation of the Actionable interface
func (id *importDump) ApplyFilters() error {
	return nil
}

func (id *importDump) importDump(config *indexConfig) error {
	schema, err := config.format()

	if err != nil {
		return err
	}

	if err := id.indexer.CreateIndex(id.indexName(config.name), schema); err != nil {
		return err
	}

	time.Sleep(1)

	// Implement importing the data

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
