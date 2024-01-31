package models

import (
	"net/http"
	"time"
)

type Bucket struct {
	Name string `json:"name"`
}

type ObjectName struct {
	Name   string
	Bucket string
}

type Object struct {
	Name        string
	Data        string
	Size        int64
	ContentType string

	Bucket string
}

type ObjectStat struct {
	// An ETag is optionally set to md5sum of an object.  In case of multipart objects,
	// ETag is of the form MD5SUM-N where MD5SUM is md5sum of all individual md5sums of
	// each parts concatenated into one string.
	ETag string `json:"etag"`

	Key          string    `json:"name"`         // Name of the object
	LastModified time.Time `json:"lastModified"` // Date and time the object was last modified.
	Size         int64     `json:"size"`         // Size in bytes of the object.
	ContentType  string    `json:"contentType"`  // A standard MIME type describing the format of the object data.
	Expires      time.Time `json:"expires"`      // The date and time at which the object is no longer able to be cached.

	VersionID string `json:"version_id"`

	// Collection of additional metadata on the object.
	// eg: x-amz-meta-*, content-encoding etc.
	Metadata http.Header `json:"metadata" xml:"-"`
}

type RemoveObjectError struct {
	ObjectName string
	VersionID  string
	Err        error
}
