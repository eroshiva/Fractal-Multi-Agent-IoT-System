// Package common is a placeholder for common structures to avoid circular dependencies
package common

// MapKey holds a complex key for accessing data of the measurement
type MapKey struct {
	Depth     int
	AppNumber int
	Instances int
}
