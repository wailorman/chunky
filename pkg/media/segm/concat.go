package segm

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/files"
	"github.com/wailorman/fftb/pkg/media/ff"
)

// ConcatOperation _
type ConcatOperation struct {
	ctx              context.Context
	ffctx            context.Context
	ffcancel         func()
	wg               *chwg.ChannelledWaitGroup
	outFile          files.Filer
	segments         []*Segment
	segmentsListFile files.Filer
	tmpPath          files.Pather
	ffworker         *ff.Instance
	initialized      bool
	started          bool
}

// ConcatRequest _
type ConcatRequest struct {
	OutFile  files.Filer
	Segments []*Segment
}

// NewConcatOperation _
func NewConcatOperation(ctx context.Context) *ConcatOperation {
	ffctx, ffcancel := context.WithCancel(ctx)

	return &ConcatOperation{
		ctx:      ctx,
		ffctx:    ffctx,
		ffcancel: ffcancel,
		wg:       chwg.New(),
	}
}

// Init _
func (co *ConcatOperation) Init(req ConcatRequest) error {
	if co.initialized {
		return ErrAlreadyInitialized
	}

	var err error

	co.outFile = req.OutFile
	co.segments = req.Segments

	co.tmpPath, err = createTmpSubdir(co.outFile.BuildPath())

	if err != nil {
		return errors.Wrap(err, "Create temp path for segments list file")
	}

	co.segmentsListFile = co.tmpPath.BuildFile("segments.txt")

	err = co.segmentsListFile.Create()

	if err != nil {
		return errors.Wrap(err, "Create temp segments list file")
	}

	writer, err := co.segmentsListFile.WriteContent()

	if err != nil {
		return errors.Wrap(err, "Building temp segments list file writer")
	}

	segmentsListContent := createSegmentsList(co.segments)

	_, err = io.WriteString(writer, segmentsListContent)

	if err != nil && err != io.EOF {
		return errors.Wrap(err, "Writing segments list")
	}

	co.ffworker = ff.New(co.ffctx)
	err = co.ffworker.Init(co.segmentsListFile, req.OutFile)

	if err != nil {
		return errors.Wrap(err, "Initializing ffworker")
	}

	mediaFile := co.ffworker.MediaFile()
	mediaFile.SetUnsafe(true)
	mediaFile.SetVideoCodec("copy")
	mediaFile.SetAudioCodec("copy")
	mediaFile.SetInputFormat("concat")

	co.initialized = true

	return nil
}

// Run _
func (co *ConcatOperation) Run() (progress chan ff.Progressable, failures chan error) {
	progress = make(chan ff.Progressable)
	failures = make(chan error)

	co.wg.Add(1)

	go func() {
		defer close(progress)
		defer close(failures)
		defer co.wg.Done()

		if !co.initialized {
			failures <- ErrNotInitialized
			return
		}

		if co.started {
			failures <- ErrAlreadyInitialized
			return
		}

		fProgress, fFailures := co.ffworker.Start()

		for {
			select {
			case failure, failed := <-fFailures:
				if !failed {
					return
				}

				co.ffcancel()
				<-co.ffworker.Closed()
				failures <- failure
				return

			case progressMessage, ok := <-fProgress:
				if ok {
					progress <- progressMessage
				}
			}
		}
	}()

	return progress, failures
}

// Prune _
func (co *ConcatOperation) Prune() error {
	if co.segmentsListFile != nil && co.segmentsListFile.IsExist() {
		return co.segmentsListFile.Remove()
	}

	return nil
}

// Closed _
func (co *ConcatOperation) Closed() <-chan struct{} {
	return co.wg.Closed()
}
