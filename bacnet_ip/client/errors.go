package client

import "errors"

var ErrNotBACnetIP = errors.New("packet is not an IP payload")
var ErrSegmentationNotSupported = errors.New("segmentation is not supported for requested operation")
