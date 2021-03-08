package local

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/subchen/go-trylock/v2"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// LockSegmentTimeout _
const LockSegmentTimeout = time.Duration(10 * time.Second)

// Dealer _
type Dealer struct {
	storageController models.IStorageController
	registry          models.IDealerRegistry
	freeSegmentLock   trylock.TryLocker
	logger            logrus.FieldLogger
	ctx               context.Context
}

// NewDealer _
func NewDealer(ctx context.Context, sc models.IStorageController, r models.IDealerRegistry) (*Dealer, error) {
	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, "fftb.dealer"); logger == nil {
		logger = ctxlog.New("fftb.dealer")
	}

	return &Dealer{
		storageController: sc,
		registry:          r,
		freeSegmentLock:   trylock.New(),
		logger:            logger,
		ctx:               ctx,
	}, nil
}

func (d *Dealer) getInputStorageClaim(segmentID string) (models.IStorageClaim, error) {
	segment, err := d.registry.FindSegmentByID(segmentID)

	if err != nil {
		return nil, err
	}

	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return nil, models.ErrUnknownSegmentType
	}

	if convertSegment.InputStorageClaimIdentity == "" {
		return nil, errors.Wrap(models.ErrMissingStorageClaim, "Getting input storage claim identity")
	}

	claim, err := d.storageController.BuildStorageClaim(convertSegment.InputStorageClaimIdentity)

	if err != nil {
		return nil, errors.Wrap(err, "Building storage claim from identity")
	}

	return claim, nil
}

func (d *Dealer) getOutputStorageClaim(segmentID string) (models.IStorageClaim, error) {
	segment, err := d.registry.FindSegmentByID(segmentID)

	if err != nil {
		return nil, err
	}

	convertSegment, ok := segment.(*models.ConvertSegment)

	if !ok {
		return nil, models.ErrUnknownSegmentType
	}

	if convertSegment.OutputStorageClaimIdentity == "" {
		return nil, errors.Wrap(models.ErrMissingStorageClaim, "Getting output storage claim identity")
	}

	claim, err := d.storageController.BuildStorageClaim(convertSegment.OutputStorageClaimIdentity)

	if err != nil {
		return nil, errors.Wrap(err, "Building storage claim from identity")
	}

	return claim, nil
}

func (d *Dealer) tryPurgeInputStorageClaim(segmentID string) {
	logger := d.logger.WithField(dlog.KeySegmentID, segmentID)

	inputClaim, err := d.getInputStorageClaim(segmentID)

	if err != nil {
		logger.WithError(err).
			Warn("Problem with getting input storage claim")
	}

	if inputClaim != nil {
		err = d.storageController.PurgeStorageClaim(inputClaim)

		if err != nil {
			logger.WithError(err).
				WithField(dlog.KeyStorageClaim, inputClaim.GetID()).
				Error("Purging input storage claim")
		}
	}
}

func (d *Dealer) tryPurgeOutputStorageClaim(segmentID string) {
	logger := d.logger.WithField(dlog.KeySegmentID, segmentID)

	outputClaim, err := d.getOutputStorageClaim(segmentID)

	if err != nil {
		logger.WithError(err).
			Warn("Problem with getting output storage claim")
	}

	if outputClaim != nil {
		err = d.storageController.PurgeStorageClaim(outputClaim)

		if err != nil {
			logger.WithError(err).
				WithField(dlog.KeyStorageClaim, outputClaim.GetID()).
				Error("Purging output storage claim")
		}
	}
}
