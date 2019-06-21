package actions

import (
	"errors"
	"scrubber/actions/options"
	"time"

	"github.com/alejandro-carstens/golastic"
)

type deleteSnapshots struct {
	filterAction
	options *options.DeleteSnapshotsOptions
}

func (ds *deleteSnapshots) ApplyOptions() Actionable {
	ds.options = ds.context.Options().(*options.DeleteSnapshotsOptions)

	ds.builder.SetOptions(&golastic.IndexOptions{Timeout: ds.options.TimeoutInSeconds()})

	return ds
}

func (ds *deleteSnapshots) Perform() Actionable {
	ds.exec(func(snapshot string) error {
		count := 0

		for {
			response, err := ds.builder.DeleteSnapshot(ds.options.Repository, snapshot)

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
