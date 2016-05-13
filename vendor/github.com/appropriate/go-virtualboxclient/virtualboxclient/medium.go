package virtualboxclient

import (
	"github.com/appropriate/go-virtualboxclient/vboxwebsrv"
)

type Medium struct {
	virtualbox      *VirtualBox
	managedObjectId string
}

func (m *Medium) CreateBaseStorage(logicalSize int64, variant []*vboxwebsrv.MediumVariant) (*Progress, error) {
	request := vboxwebsrv.IMediumcreateBaseStorage{This: m.managedObjectId, LogicalSize: logicalSize, Variant: variant}

	response, err := m.virtualbox.IMediumcreateBaseStorage(&request)
	if err != nil {
		return nil, err // TODO: Wrap the error
	}

	// TODO: See if we need to do anything with the response
	return &Progress{managedObjectId: response.Returnval}, nil
}

func (m *Medium) DeleteStorage() (*Progress, error) {
	request := vboxwebsrv.IMediumdeleteStorage{This: m.managedObjectId}

	response, err := m.virtualbox.IMediumdeleteStorage(&request)
	if err != nil {
		return nil, err // TODO: Wrap the error
	}

	// TODO: See if we need to do anything with the response
	return &Progress{managedObjectId: response.Returnval}, nil
}
