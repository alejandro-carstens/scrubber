package actions

import (
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

func (ic *indexConfig) format() (*gabs.Container, error) {
	return nil, nil
}

type importDump struct {
	action
	options *options.ImportOptions
}

// ApplyOptions implementation of the Actionable interface
func (id *importDump) ApplyOptions() Actionable {
	id.options = id.context.Options().(*options.ImportOptions)

	id.indexer.SetOptions(&golastic.IndexOptions{Timeout: id.options.TimeoutInSeconds()})

	return id
}

// Perform implementation of the Actionable interface
func (id *importDump) Perform() Actionable {
	// For each index directory we need to do the following:
	// - Read the settings.json file, filter out the settings we do not need,
	//   apply the settings and do the same for mappings and aliases. I believe
	//   we can do this on the same index creation request.
	// - Sleep a second
	// - Start batch inserting the documents onto the new index

	configs, err := id.getIndexConfigs()

	if err != nil {
		id.errorContainer.push(id.name, "_all", err)

		return id
	}

	pool := grpool.NewPool(id.options.Concurrency, len(configs))

	defer pool.Release()

	for _, config := range configs {
		pool.JobQueue <- func() {
			if err := id.importDump(config); err != nil {
				id.errorContainer.push(id.name, id.indexName(config.name), err)
			}

			pool.JobDone()
		}
	}

	pool.WaitCount(len(configs))

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

	if err := id.indexer.CreateIndex(id.indexName(config.name), schema.String()); err != nil {
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
