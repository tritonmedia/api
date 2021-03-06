// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"github.com/facebook/ent/dialect/sql/schema"
	"github.com/facebook/ent/schema/field"
)

var (
	// MediaColumns holds the columns for the "media" table.
	MediaColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "title", Type: field.TypeString},
		{Name: "type", Type: field.TypeString},
		{Name: "source", Type: field.TypeString},
		{Name: "source_uri", Type: field.TypeString},
		{Name: "status", Type: field.TypeString, Default: "MediaStatusStageDownloadQueued"},
		{Name: "status_percent", Type: field.TypeFloat32},
	}
	// MediaTable holds the schema information for the "media" table.
	MediaTable = &schema.Table{
		Name:        "media",
		Columns:     MediaColumns,
		PrimaryKey:  []*schema.Column{MediaColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
		Indexes: []*schema.Index{
			{
				Name:    "media_source_source_uri",
				Unique:  true,
				Columns: []*schema.Column{MediaColumns[3], MediaColumns[4]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		MediaTable,
	}
)

func init() {
}
