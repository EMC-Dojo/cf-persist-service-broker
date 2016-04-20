package storage

type ScaleIODriver struct {
}

type VolumeAttachment struct  {
}

type StorageDriver interface{
  VolumeCreate(ctx Context, name string, opts *VolumeCreateOpts) (*Volume, error)
}

type Volume struct {
  // The volume's attachments.
  Attachments []*VolumeAttachment `json:"attachments,omitempty"`

  // The availability zone for which the volume is available.
  AvailabilityZone string `json:"availabilityZone,omitempty"`

  // The volume IOPs.
  IOPS int64 `json:"iops,omitempty"`

  // The name of the volume.
  Name string `json:"name"`

  // NetworkName is the name the device is known by in order to discover
  // locally.
  NetworkName string `json:"networkName,omitempty"`

  // The size of the volume.
  Size int64 `json:"size,omitempty"`

  // The volume status.
  Status string `json:"status,omitempty"`

  // The volume ID.
  ID string `json:"id"`

  // The volume type.
  Type string `json:"type"`

  // Fields are additional properties that can be defined for this type.
  Fields map[string]string `json:"fields,omitempty"`
}

type Context struct {
}

type VolumeCreateOpts struct {
}

func (s *ScaleIODriver) VolumeCreate(ctx Context, name string, opts *VolumeCreateOpts) (*Volume, error) {
  volume := Volume{
    Name: "volume1",
    Size: 8,
    ID:   "volume-id",
  }
  return &volume, nil
}