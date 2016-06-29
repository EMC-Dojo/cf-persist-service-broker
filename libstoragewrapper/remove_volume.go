package libstoragewrapper

import (
	"github.com/emccode/libstorage/api/context"
	"github.com/emccode/libstorage/api/types"
)

// RemoveVolume : Removes Volume from storage array using Libstorage API
func RemoveVolume(c types.APIClient, serviceName, volumeID string) error {
	return c.VolumeRemove(context.Background(), serviceName, volumeID)
}
