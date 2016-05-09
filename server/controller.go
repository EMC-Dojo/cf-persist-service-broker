package server

import (
	"github.com/emccode/libstorage/api/types"
)

func CreateVolume(c types.Client, ctx types.Context, volumeName, availabilityZone, volumeType string, iops, size int64) (*types.Volume, error) {
	volumeCreateOpts := &types.VolumeCreateOpts{
		AvailabilityZone: &availabilityZone,
		IOPS:             &iops,
		Size:             &size,
		Type:             &volumeType,
	}

	return c.Storage().VolumeCreate(ctx, volumeName, volumeCreateOpts)
}

func RemoveVolume(c types.Client, ctx types.Context, volumeID string) error {
	return c.Storage().VolumeRemove(ctx, volumeID, nil)
}
