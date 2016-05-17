package libstoragewrapper

import (
	"github.com/emccode/libstorage/api/types"
)

func CreateVolume(c types.Client, ctx types.Context, volumeName string, volumeCreateOpts types.VolumeCreateOpts) (*types.Volume, error) {
	return c.Storage().VolumeCreate(ctx, volumeName, &volumeCreateOpts)
}