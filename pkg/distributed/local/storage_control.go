package local

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/files"
	// "github.com/machinebox/progress"
	// "github.com/machinebox/progress"
)

// ErrStorageClaimMissingFile _
var ErrStorageClaimMissingFile = errors.New("Storage claim missing file")

// ErrStorageClaimAlreadyAllocated _
var ErrStorageClaimAlreadyAllocated = errors.New("Storage claim already allocated")

// StorageControl _
type StorageControl struct {
	storagePath files.Pather
}

// NewStorageControl _
func NewStorageControl(path files.Pather) *StorageControl {
	return &StorageControl{
		storagePath: path,
	}
}

// AllocateStorageClaim _
func (sc *StorageControl) AllocateStorageClaim(ctx context.Context, identity string) (models.IStorageClaim, error) {
	file := sc.storagePath.BuildFile(identity)

	err := file.EnsureParentDirExists()

	if err != nil {
		return nil, errors.Wrap(err, "Creating directory for storage claim")
	}

	if file.IsExist() {
		return nil, ErrStorageClaimAlreadyAllocated
	}

	err = file.Create()

	if err != nil {
		return nil, errors.Wrap(err, "Creating file for storage claim")
	}

	claim := &StorageClaim{
		identity: identity,
		file:     file,
	}

	return claim, nil
}

// BuildStorageClaim _
func (sc *StorageControl) BuildStorageClaim(identity string) (models.IStorageClaim, error) {
	claimFile := sc.storagePath.BuildFile(identity)

	if claimFile.IsExist() == false {
		return nil, ErrStorageClaimMissingFile
	}

	size, err := claimFile.Size()

	if err != nil {
		return nil, errors.Wrap(err, "Getting claim file size")
	}

	return &StorageClaim{
		identity: identity,
		file:     claimFile,
		size:     size,
	}, nil
}

// PurgeStorageClaim _
func (sc *StorageControl) PurgeStorageClaim(ctx context.Context, claim models.IStorageClaim) error {
	localClaim, ok := claim.(*StorageClaim)

	if !ok {
		return models.ErrUnknownStorageClaimType
	}

	if localClaim.file == nil {
		return ErrStorageClaimMissingFile
	}

	err := localClaim.file.Remove()

	if err != nil {
		return errors.Wrap(err, "Removing file")
	}

	return nil
}