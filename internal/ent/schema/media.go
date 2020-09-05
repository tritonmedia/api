package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
	"github.com/google/uuid"
)

// Media schema.
type Media struct {
	ent.Schema
}

// Fields of the Media.
func (Media) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique().Immutable().Default(uuid.New),
		field.String("title"),
		field.String("type"),
		field.String("source"),
		field.String("source_uri"),
		// NOTE: This is tied to the string value in the proto.
		field.String("status").Default("MediaStatusStageDownloadQueued"),
		field.Float32("status_percent").Default(float32(0)),
	}
}
