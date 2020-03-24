package actions

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/Jeffail/gabs"
	"github.com/ivpusic/grpool"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/options"
	"github.com/alejandro-carstens/scrubber/filesystem"
)

const CHUNK int = 20
const KEEP_ALIVE string = "10m"

type fetch func(index string) (*gabs.Container, error)

type counter int64

func (c *counter) increment() int64 {
	return atomic.AddInt64((*int64)(c), 1)
}

func (c *counter) reset() {
	atomic.StoreInt64((*int64)(c), 0)
}

func (c *counter) get() int64 {
	return atomic.LoadInt64((*int64)(c))
}

type semaphore int32

func (s *semaphore) turnOn() {
	atomic.StoreInt32((*int32)(s), 1)
}

func (s *semaphore) get() bool {
	return atomic.LoadInt32((*int32)(s)) == int32(1)
}

type dump struct {
	filterAction
	options   *options.DumpOptions
	counter   *counter
	semaphore *semaphore
}

// ApplyOptions implementation of the Actionable interface
func (d *dump) ApplyOptions() Actionable {
	d.counter = new(counter)
	d.semaphore = new(semaphore)
	d.options = d.context.Options().(*options.DumpOptions)

	d.indexer.SetOptions(&golastic.IndexOptions{Timeout: d.options.TimeoutInSeconds()})

	return d
}

// Perform implementation of the Actionable interface
func (d *dump) Perform() Actionable {
	d.exec(func(index string) error {
		fs, err := filesystem.Build(d.filesystemConfig())

		if err != nil {
			return err
		}

		done := make(chan bool)
		failed := make(chan error)

		go d.process(index, "mappings", fs, failed, done, func(index string) (*gabs.Container, error) {
			return d.indexer.Mappings(index)
		})

		go d.process(index, "aliases", fs, failed, done, func(index string) (*gabs.Container, error) {
			return d.indexer.Aliases(index)
		})

		go d.scroll(index, failed, done)

		go d.process(index, "settings", fs, failed, done, func(index string) (*gabs.Container, error) {
			response, err := d.indexer.Settings(index)

			if err != nil {
				return nil, err
			}

			settings, valid := response[index]

			if !valid {
				return nil, errors.New(fmt.Sprintf("no settings found for index: %v", index))
			}

			return settings, nil
		})

		select {
		case err := <-failed:
			d.semaphore.turnOn()
			d.close(failed, done)

			return err
		case <-done:
			d.close(failed, done)

			return nil
		}
	})

	return d
}

func (d *dump) process(index, name string, fs filesystem.Storeable, failed chan error, done chan bool, fn fetch) {
	container, err := fn(index)

	if err != nil {
		failed <- err

		return
	}

	if err := fs.Put(d.fileName(name), strings.NewReader(container.String())); err != nil {
		failed <- err

		return
	}

	if d.counter.increment() == 4 {
		done <- true
	}
}

func (d *dump) scroll(index string, failed chan error, done chan bool) {
	pool := grpool.NewPool(d.options.Concurrency, d.options.Concurrency)

	defer pool.Release()

	pool.WaitCount(d.options.Concurrency)

	for i := 0; i < d.options.Concurrency; i++ {
		var builder *golastic.Builder

		if d.options.Concurrency > 1 {
			builder = d.builder(index).InitSlicedScroller(i, d.options.Concurrency, CHUNK, KEEP_ALIVE)
		} else {
			builder = d.builder(index).InitScroller(CHUNK, KEEP_ALIVE)
		}

		stream, err := d.openStream(fmt.Sprintf("data_%v", i))

		if err != nil {
			failed <- err

			return
		}

		pool.JobQueue <- func() {
			defer pool.JobDone()

			go func() {
				if err := stream.Stream(); err != nil {
					failed <- err
				}
			}()

			for {
				response, err := builder.Scroll()

				if err == io.EOF {
					if err := stream.Close(); err != nil {
						failed <- err
					}

					return
				}

				if err != nil {
					failed <- err

					return
				}

				items, err := response.S("hits", "hits").Children()

				if err != nil {
					failed <- err

					return
				}

				for _, item := range items {
					if d.semaphore.get() {
						if err := stream.Close(); err != nil {
							d.reporter.logger.Errorf(err.Error())
						}

						return
					}

					stream.Channel(item.String() + "\n")
				}
			}
		}
	}

	pool.WaitAll()

	if d.counter.increment() == 4 {
		done <- true
	}
}

func (d *dump) openStream(fileName string) (filesystem.Storeable, error) {
	fs, err := filesystem.Build(d.filesystemConfig())

	if err != nil {
		return nil, err
	}

	err = fs.OpenStream(d.fileName(fileName))

	return fs, err
}

func (d *dump) filesystemConfig() filesystem.Configurable {
	return &filesystem.Local{
		Path: filepath.Join(d.options.Path, filepath.FromSlash(d.options.Name)),
	}
}

func (d *dump) fileName(name string) string {
	return fmt.Sprintf("%v.json", name)
}

func (d *dump) builder(index string) *golastic.Builder {
	builder := d.connection.Builder(index)

	buildQuery(builder, d.options.Criteria)

	return builder
}

func (d *dump) close(failed chan error, done chan bool) {
	close(failed)
	close(done)
}
