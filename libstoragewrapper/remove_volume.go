package libstoragewrapper

import (
	"github.com/emccode/libstorage/api/types"
)

func RemoveVolume(c types.Client, ctx types.Context, volumeID string) error {
	return c.Storage().VolumeRemove(ctx, volumeID, nil)
}
