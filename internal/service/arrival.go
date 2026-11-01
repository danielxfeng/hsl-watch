package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/danielxfeng/hsl-watch/internal/hsl"
	"github.com/danielxfeng/hsl-watch/internal/model"
)

type ArrivalService struct {
	client *hsl.HslClient
	cache  atomic.Pointer[model.Arrival]
	ids    []string
	debug  bool
}

func NewArrivalService(client *hsl.HslClient, ids []string, debug bool) *ArrivalService {
	return &ArrivalService{
		client: client,
		ids:    append([]string(nil), ids...),
		debug:  debug,
	}
}

func (s *ArrivalService) UpdateArrivals(ctx context.Context) error {
	if len(s.ids) == 0 {
		return errors.New("no stop ids configured")
	}

	arrivals, err := s.client.GetArrivals(ctx, s.ids)
	if err != nil {
		return fmt.Errorf("fetch arrivals: %w", err)
	}

	stops := arrivals.Data.Stops[:0]
	for _, stop := range arrivals.Data.Stops {
		if stop.GtfsID != "" {
			stops = append(stops, stop)
		}
	}
	arrivals.Data.Stops = stops

	if len(stops) == 0 {
		return errors.New("HSL returned no valid stops, check ARRIVAL_IDS")
	}
	if len(stops) < len(s.ids) {
		log.Printf("warning: %d of %d configured stops not found, check ARRIVAL_IDS", len(s.ids)-len(stops), len(s.ids))
	}
	arrivals.FetchedAt = time.Now().Unix()

	for _, stop := range arrivals.Data.Stops {
		log.Printf("arrivals: stop %s (%s): %d upcoming", stop.GtfsID, stop.Name, len(stop.Arrivals))
	}

	if s.debug {
		if b, err := json.MarshalIndent(arrivals, "", "  "); err == nil {
			log.Printf("arrivals: full payload:\n%s", b)
		}
	}

	s.cache.Store(arrivals)
	return nil
}

func (s *ArrivalService) GetArrivals() *model.Arrival {
	return s.cache.Load()
}
