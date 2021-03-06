package convert

import (
	"github.com/pkg/errors"
	ffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"
)

// H264Codec _
type H264Codec struct {
	task     Task
	metadata ffmpegModels.Metadata
}

// NewH264Codec _
func NewH264Codec(task Task, metadata ffmpegModels.Metadata) *H264Codec {
	return &H264Codec{
		task:     task,
		metadata: metadata,
	}
}

func (c *H264Codec) configure(mediaFile *ffmpegModels.Mediafile) error {
	var err error

	mediaFile.SetVideoCodec("libx264")
	mediaFile.SetPreset(c.task.Params.Preset)
	mediaFile.SetHideBanner(true)
	mediaFile.SetVsync(true)
	mediaFile.SetAudioCodec("copy")
	mediaFile.SetMaxMuxingQueueSize(102400)

	if c.task.Params.VideoQuality > 0 {
		mediaFile.SetConstantQuantization(c.task.Params.VideoQuality)
	} else {
		mediaFile.SetVideoBitRate(c.task.Params.VideoBitRate)
	}

	mediaFile.SetKeyframeInterval(c.task.Params.KeyframeInterval)

	hwaccel := chooseHwAccel(c.task, c.metadata)

	if err = hwaccel.configure(mediaFile); err != nil {
		return errors.Wrap(err, "Configuring hwaccel")
	}

	return nil
}

func (c *H264Codec) getType() string {
	return H264CodecType
}
