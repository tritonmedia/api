// Code generated by entc, DO NOT EDIT.

package ent

import (
	"github.com/google/uuid"
	"github.com/tritonmedia/api/internal/ent/media"
	"github.com/tritonmedia/api/internal/ent/schema"
)

// The init function reads all schema descriptors with runtime
// code (default values, validators or hooks) and stitches it
// to their package variables.
func init() {
	mediaFields := schema.Media{}.Fields()
	_ = mediaFields
	// mediaDescStatus is the schema descriptor for status field.
	mediaDescStatus := mediaFields[5].Descriptor()
	// media.DefaultStatus holds the default value on creation for the status field.
	media.DefaultStatus = mediaDescStatus.Default.(string)
	// mediaDescStatusPercent is the schema descriptor for status_percent field.
	mediaDescStatusPercent := mediaFields[6].Descriptor()
	// media.DefaultStatusPercent holds the default value on creation for the status_percent field.
	media.DefaultStatusPercent = mediaDescStatusPercent.Default.(float32)
	// mediaDescID is the schema descriptor for id field.
	mediaDescID := mediaFields[0].Descriptor()
	// media.DefaultID holds the default value on creation for the id field.
	media.DefaultID = mediaDescID.Default.(func() uuid.UUID)
}
