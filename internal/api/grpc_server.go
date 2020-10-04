package api

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	apiv1 "github.com/tritonmedia/api/api/v1"
	"github.com/tritonmedia/api/internal/ent"
	"github.com/tritonmedia/pkg/discovery"

	///StartBlock(imports)
	// goimports is fucking up these imports.
	"github.com/facebook/ent/dialect"
	entsql "github.com/facebook/ent/dialect/sql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	///EndBlock(imports)
)

var nonAlphaNumRegex = regexp.MustCompile("[^a-zA-Z0-9]+")

type GRPCServiceHandler struct {
	log logrus.FieldLogger

	///StartBlock(grpcConfig)
	db *ent.Client
	sc stan.Conn
	///EndBlock(grpcConfig)
}

///StartBlock(global)
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

///EndBlock(global)

func NewServiceHandler(ctx context.Context, log logrus.FieldLogger) (*GRPCServiceHandler, error) {
	///StartBlock(grpcInit)
	// TODO(jaredallard): we could make this easier to work with...
	log = log.WithField("service", "*api.GRPCServiceHandler")

	// TODO(jaredallard): when we add configuration, we need to change this
	conf, err := pgx.ParseConfig("postgres://api:yeAUemR82sK82jcNjR0E8BqYejUUYtLM@127.0.0.1:5432/triton")
	if err != nil {
		return nil, errors.Wrap(err, "failed to build database config")
	}

	sdb := stdlib.OpenDB(*conf)
	db := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, sdb)))

	log.Info("running database migrations")

	//nolint:govet
	if err := db.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	endpoint, err := discovery.Find("nats")
	if err != nil {
		return nil, errors.Wrap(err, "failed to find nats")
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hostname")
	}
	clientID := nonAlphaNumRegex.ReplaceAllString(hostname, "-")

	// TODO(jaredallard): handle connection loss
	sc, err := stan.Connect("test-cluster", clientID, stan.NatsURL(endpoint))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create nats client")
	}
	///EndBlock(grpcInit)

	return &GRPCServiceHandler{
		log,
		///StartBlock(grpcConfigInit)
		db,
		sc,
		///EndBlock(grpcConfigInit)
	}, nil
}

///StartBlock(grpcHandlers)
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

	h.log.WithField("media", r.Media).Info("creating media")

	m, err := h.db.Media.Create().SetTitle(r.Media.Title).
		SetSource(r.Media.Source.String()).SetSourceURI(r.Media.SourceURI).
		SetType(r.Media.Type.String()).Save(ctx)
	if err != nil {
		return nil, err
	}

	p := dbMediaToProto(m)
	b, err := proto.Marshal(p)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal into protobuf format")
	}
	if err := h.sc.Publish("v1.convert", b); err != nil {
		return nil, errors.Wrap(err, "failed to publish message")
	}

	return p, nil
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

///EndBlock(grpcHandlers)
