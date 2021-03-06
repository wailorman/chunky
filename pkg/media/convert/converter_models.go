package convert

import (
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"
)

const (
	// NvencHWAccelType _
	NvencHWAccelType = "nvenc"
	// VTBHWAccelType _
	VTBHWAccelType = "videotoolbox"
)

const (
	// HevcCodecType _
	HevcCodecType = "hevc"
	// H264CodecType _
	H264CodecType = "h264"
)

// VideoFileFilteringMessage _
type VideoFileFilteringMessage struct {
	File    files.Filer
	IsVideo bool
	Err     error
}

// BatchTask _
type BatchTask struct {
	Parallelism           int    `yaml:"parallelism"`
	StopConversionOnError bool   `yaml:"stop_conversion_on_error"`
	Tasks                 []Task `yaml:"tasks"`
}

// Task _
type Task struct {
	ID      string `yaml:"id"`
	InFile  string `yaml:"in_file"`
	OutFile string `yaml:"out_file"`
	Params  Params
}

// Params _
type Params struct {
	VideoCodec       string `yaml:"video_codec"`
	HWAccel          string `yaml:"hw_accel"`
	VideoBitRate     string `yaml:"video_bit_rate"`
	VideoQuality     int    `yaml:"video_quality"`
	Preset           string `yaml:"preset"`
	Scale            string `yaml:"scale"`
	KeyframeInterval int    `yaml:"keyframe_interval"`
}

// ErrFileIsNotVideo _
var ErrFileIsNotVideo = errors.New("File is not a video")

// ErrNoStreamsInFile _
var ErrNoStreamsInFile = errors.New("No streams in file")

// ErrCodecIsNotSupportedByEncoder _
var ErrCodecIsNotSupportedByEncoder = errors.New("Codec is not supported by encoder")

// ErrUnsupportedHWAccelType _
var ErrUnsupportedHWAccelType = errors.New("Unsupported hardware acceleration type")

// ErrUnsupportedScale _
var ErrUnsupportedScale = errors.New("Unsupported scale")

// ErrResolutionNotSupportScaling _
var ErrResolutionNotSupportScaling = errors.New("Resolution not support scaling")

// ErrOutputFileExistsOrIsDirectory _
var ErrOutputFileExistsOrIsDirectory = errors.New("Output file exists or is directory")

// ErrVtbQualityNotSupported _
var ErrVtbQualityNotSupported = errors.New("Video quality option is not supported by Apple VideoToolBox")
