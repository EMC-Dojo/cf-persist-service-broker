package virtualboxclient

import (
	"github.com/appropriate/go-virtualboxclient/vboxwebsrv"
)

type Machine struct {
	virtualbox      *VirtualBox
	managedObjectId string
}

func (m *Machine) GetChipsetType() (*vboxwebsrv.ChipsetType, error) {
	request := vboxwebsrv.IMachinegetChipsetType{This: m.managedObjectId}

	response, err := m.virtualbox.IMachinegetChipsetType(&request)
	if err != nil {
		return nil, err // TODO: Wrap the error
	}

	return response.Returnval, nil
}

func (m *Machine) GetMediumAttachments() ([]*vboxwebsrv.IMediumAttachment, error) {
	request := vboxwebsrv.IMachinegetMediumAttachments{This: m.managedObjectId}

	response, err := m.virtualbox.IMachinegetMediumAttachments(&request)
	if err != nil {
		return nil, err // TODO: Wrap the error
	}

	return response.Returnval, nil
}

func (m *Machine) GetNetworkAdapter(slot uint32) (*NetworkAdapter, error) {
	request := vboxwebsrv.IMachinegetNetworkAdapter{This: m.managedObjectId, Slot: slot}

	response, err := m.virtualbox.IMachinegetNetworkAdapter(&request)
	if err != nil {
		return nil, err // TODO: Wrap the error
	}

	return &NetworkAdapter{m.virtualbox, response.Returnval}, nil
}

func (m *Machine) GetSettingsFilePath() (string, error) {
	request := vboxwebsrv.IMachinegetSettingsFilePath{This: m.managedObjectId}

	response, err := m.virtualbox.IMachinegetSettingsFilePath(&request)
	if err != nil {
		return "", err // TODO: Wrap the error
	}

	return response.Returnval, nil
}

func (m *Machine) GetStorageControllers() ([]*StorageController, error) {
	request := vboxwebsrv.IMachinegetStorageControllers{This: m.managedObjectId}

	response, err := m.virtualbox.IMachinegetStorageControllers(&request)
	if err != nil {
		return nil, err // TODO: Wrap the error
	}

	storageControllers := make([]*StorageController, len(response.Returnval))
	for i, oid := range response.Returnval {
		storageControllers[i] = &StorageController{m.virtualbox, oid}
	}

	return storageControllers, nil
}
