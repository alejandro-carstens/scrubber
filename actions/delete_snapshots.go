package actions

import (
	"errors"
	"time"

	"github.com/alejandro-carstens/golastic"
	"scrubber/actions/options"
)

type deleteSnapshots struct {
	filterAction
	options *options.DeleteSnapshotsOptions
}

// ApplyOptions implementation of the Actionable interface
func (ds *deleteSnapshots) ApplyOptions() Actionable {
	ds.options = ds.context.Options().(*options.DeleteSnapshotsOptions)

	ds.indexer.SetOptions(&golastic.IndexOptions{Timeout: ds.options.TimeoutInSeconds()})

	return ds
}

// Perform implementation of the Actionable interface
func (ds *deleteSnapshots) Perform() Actionable {
	ds.exec(func(snapshot string) error {
		count := 0

		for {
			response, err := ds.indexer.DeleteSnapshot(ds.options.Repository, snapshot)

			if err != nil && count < ds.options.RetryCount {
				time.Sleep(time.Duration(int64(ds.options.RetryInterval)) * time.Second)
				count++

				continue
			} else if err != nil && count == ds.options.RetryCount {
				return err
			}

			if acknowledged, _ := response.S("acknowledged").Data().(bool); !acknowledged && count < ds.options.RetryCount {
				time.Sleep(time.Duration(int64(ds.options.RetryInterval)) * time.Second)
				count++

				continue
			} else if !acknowledged && count == ds.options.RetryCount {
				return errors.New("delete_snapshot action was not acknowledged")
			}

			return nil
		}
	})

	return ds
}
