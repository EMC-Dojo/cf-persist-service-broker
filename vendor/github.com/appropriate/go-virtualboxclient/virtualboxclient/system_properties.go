package virtualboxclient

import (
	"github.com/appropriate/go-virtualboxclient/vboxwebsrv"
)

type SystemProperties struct {
	virtualbox      *VirtualBox
	managedObjectId string
}

func (sp *SystemProperties) GetMaxNetworkAdapters(chipset *vboxwebsrv.ChipsetType) (uint32, error) {
	request := vboxwebsrv.ISystemPropertiesgetMaxNetworkAdapters{This: sp.managedObjectId, Chipset: chipset}

	response, err := sp.virtualbox.ISystemPropertiesgetMaxNetworkAdapters(&request)
	if err != nil {
		return 0, err // TODO: Wrap the error
	}

	return response.Returnval, nil
}
