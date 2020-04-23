package filesystem

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// GCS represents the configuration for
// the Google Cloud Storage filesystem
type GCS struct {
	Bucket              string
	Context             context.Context
	CredentialsFilePath string
	Directory           string
}

// Validate implementation of the Configurable interface
func (g *GCS) Validate() error {
	if &g.Context == nil {
		return errors.New("context is a required field")
	}

	if len(g.Bucket) == 0 {
		return errors.New("bucket is a required field")
	}

	if len(g.CredentialsFilePath) == 0 {
		return errors.New("credentials_file_path is a required field")
	}

	return nil
}

// Name implementation of the Configurable interface
func (g *GCS) Name() string {
	return "gcs"
}

type gcs struct {
	client        *storage.Client
	writer        *storage.Writer
	bucket        string
	directory     string
	streamChannel chan string
	context       context.Context
}

// Init implementation of the Storeable interface
func (g *gcs) Init(configuration Configurable) (Storeable, error) {
	config := configuration.(*GCS)

	client, err := storage.NewClient(
		config.Context,
		option.WithCredentialsFile(config.CredentialsFilePath),
	)

	if err != nil {
		return nil, err
	}

	g.client = client
	g.bucket = config.Bucket
	g.context = config.Context
	g.directory = config.Directory

	return g, nil
}

// List implementation of the Storeable interface
func (g *gcs) List(name string) ([]string, error) {
	dirs := map[string]bool{}

	it := g.client.Bucket(g.bucket).Objects(g.context, &storage.Query{
		Prefix: name,
	})

	for {
		attributes, err := it.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		parts := strings.Split(
			strings.Replace(attributes.Name, fmt.Sprintf("%v/", name), "", -1),
			"/",
		)

		if len(parts) > 0 {
			dirs[parts[0]] = true
		}
	}

	list := []string{}

	for dir, _ := range dirs {
		list = append(list, dir)
	}

	return list, nil
}

// Get implementation of the Storeable interface
func (g *gcs) Open(name string) (io.Reader, error) {
	bucket := g.client.Bucket(g.bucket)

	if _, err := bucket.Attrs(g.context); err != nil {
		return nil, err
	}

	return bucket.Object(g.path(name)).NewReader(g.context)
}

// Put implementation of the Storeable interface
func (g *gcs) Put(name string, reader io.Reader) error {
	bucket := g.client.Bucket(g.bucket)

	if _, err := bucket.Attrs(g.context); err != nil {
		return err
	}

	writer := bucket.Object(g.path(name)).NewWriter(g.context)

	defer writer.Close()

	if _, err := io.Copy(writer, reader); err != nil {
		return err
	}

	return nil
}

// OpenStream implementation of the Storeable interface
func (g *gcs) OpenStream(name string) error {
	bucket := g.client.Bucket(g.bucket)

	if _, err := bucket.Attrs(g.context); err != nil {
		return err
	}

	g.writer = bucket.Object(g.path(name)).NewWriter(g.context)
	g.streamChannel = make(chan string)

	return nil
}

// Stream implementation of the Storeable interface
func (g *gcs) Stream() error {
	if g.writer == nil || g.streamChannel == nil {
		return errors.New("please open a stream")
	}

	for line := range g.streamChannel {
		n, err := g.writer.Write([]byte(line))

		if err != nil {
			return err
		}

		if n != len(line) {
			return errors.New("failed to write data")
		}
	}

	return nil
}

// Channel implementation of the Storeable interface
func (g *gcs) Channel(data string) {
	g.streamChannel <- data
}

// Close implementation of the Storeable interface
func (g *gcs) Close() error {
	if g.writer == nil || g.streamChannel == nil {
		return errors.New("please open a stream")
	}

	close(g.streamChannel)

	return g.writer.Close()
}

func (g *gcs) path(name string) string {
	return filepath.Join(g.directory, name)
}
