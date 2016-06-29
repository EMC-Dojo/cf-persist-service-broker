package virtualboxclient

import (
	"github.com/appropriate/go-virtualboxclient/vboxwebsrv"
)

type StorageController struct {
	virtualbox      *VirtualBox
	managedObjectId string
}

func (sc *StorageController) GetName() (string, error) {
	request := vboxwebsrv.IStorageControllergetName{This: sc.managedObjectId}

	response, err := sc.virtualbox.IStorageControllergetName(&request)
	if err != nil {
		return "", err // TODO: Wrap the error
	}

	return response.Returnval, nil
}

func (sc *StorageController) GetPortCount() (uint32, error) {
	request := vboxwebsrv.IStorageControllergetPortCount{This: sc.managedObjectId}

	response, err := sc.virtualbox.IStorageControllergetPortCount(&request)
	if err != nil {
		return 0, err // TODO: Wrap the error
	}

	return response.Returnval, nil
}
