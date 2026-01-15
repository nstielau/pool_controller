// Package pool provides high-level abstractions for pool devices.
package pool

// Device is the interface for all pool devices.
type Device interface {
	ID() interface{}
	Name() string
	State() interface{}
	HassType() string
	FriendlyState() string
}
