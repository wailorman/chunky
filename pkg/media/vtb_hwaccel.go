package media

import (
	ffmpegModels "github.com/wailorman/goffmpeg/models"
)

type vtbHWAccel struct {
	task     ConverterTask
	metadata ffmpegModels.Metadata
}

func (hw *vtbHWAccel) configure(mediaFile *ffmpegModels.Mediafile) error {
	if !isVideo(hw.metadata) {
		return ErrFileIsNotVideo
	}

	mediaFile.SetHardwareAcceleration("videotoolbox")
	mediaFile.SetPreset("")

	switch hw.task.VideoCodec {
	case HevcCodecType:
		mediaFile.SetVideoCodec("hevc_videotoolbox")
	case H264CodecType:
		mediaFile.SetVideoCodec("h264_videotoolbox")
	default:
		return ErrCodecIsNotSupportedByEncoder
	}

	return nil
}

func (hw *vtbHWAccel) getType() string {
	return VTBHWAccelType
}