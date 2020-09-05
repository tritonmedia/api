package api

import (
	"context"
	"fmt"

	"github.com/facebook/ent/dialect"
	entsql "github.com/facebook/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	// used by ent
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"

	apiv1 "github.com/tritonmedia/api/api/v1"
	"github.com/tritonmedia/api/internal/ent"
)

type GRPCServiceHandler struct {
	log logrus.FieldLogger

	db *ent.Client
}

// dbMediaToProto converts a media database object to
// it's protobuf representation
func dbMediaToProto(m *ent.Media) *apiv1.Media {
	return &apiv1.Media{
		Id:        m.ID.String(),
		Type:      apiv1.MediaType(apiv1.MediaType_value[m.Type]),
		Title:     m.Title,
		Source:    apiv1.MediaSource(apiv1.MediaSource_value[m.Source]),
		SourceURI: m.SourceURI,
		Status: &apiv1.MediaStatus{
			Status:  apiv1.MediaStatusStage(apiv1.MediaStatusStage_value[m.Status]),
			Percent: m.StatusPercent,
		},
	}
}

func NewServiceHandler(ctx context.Context, log logrus.FieldLogger) (*GRPCServiceHandler, error) {
	conf, err := pgx.ParseConfig("postgres://api:yeAUemR82sK82jcNjR0E8BqYejUUYtLM@127.0.0.1:5432/triton")
	if err != nil {
		return nil, errors.Wrap(err, "failed to build database config")
	}

	sdb := stdlib.OpenDB(*conf)
	db := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, sdb)))

	log.Info("running database migrations")
	if err := db.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return &GRPCServiceHandler{
		log,
		db,
	}, nil
}

// CreateMedia creates a new media object and sends it off for processing
func (h *GRPCServiceHandler) CreateMedia(ctx context.Context, r *apiv1.CreateMediaRequest) (*apiv1.Media, error) {
	if r.Media.Status != nil {
		return nil, fmt.Errorf("status can't be set here")
	}

	if r.Media.Id != "" {
		return nil, fmt.Errorf("id can't be set here")
	}

	if r.Media.Source == apiv1.MediaSource_MediaSourceUnset {
		return nil, fmt.Errorf("missing source")
	}

	if r.Media.SourceURI == "" {
		return nil, fmt.Errorf("missing sourceURI")
	}

	if r.Media.Type == apiv1.MediaType_MediaTypeUnset {
		return nil, fmt.Errorf("missing type")
	}

	m, err := h.db.Media.Create().SetTitle(r.Media.Title).
		SetSource(r.Media.Source.String()).SetSourceURI(r.Media.SourceURI).
		SetType(r.Media.Type.String()).Save(ctx)
	if err != nil {
		return nil, err
	}

	return dbMediaToProto(m), nil
}

// GetMedia gets a media using specific filters
func (h *GRPCServiceHandler) GetMedia(ctx context.Context, r *apiv1.GetMediaRequest) (*apiv1.Media, error) {
	if r.Id == "" {
		return nil, fmt.Errorf("missing id")
	}

	u, err := uuid.Parse(r.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse id")
	}

	m, err := h.db.Media.Get(ctx, u)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get media")
	}

	return dbMediaToProto(m), nil
}
