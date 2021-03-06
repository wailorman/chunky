package convert

import ffmpegModels "github.com/wailorman/fftb/pkg/goffmpeg/models"

// hwAccelerator _
type hwAccelerator interface {
	configure(mediaFile *ffmpegModels.Mediafile) error
	getType() string
}

func chooseHwAccel(task Task, metadata ffmpegModels.Metadata) hwAccelerator {
	switch task.Params.HWAccel {
	case NvencHWAccelType:
		return &nvencHWAccel{
			task:     task,
			metadata: metadata,
		}
	case VTBHWAccelType:
		return &vtbHWAccel{
			task:     task,
			metadata: metadata,
		}
	default:
		return &emptyHwAccel{}
	}
}

type emptyHwAccel struct {
}

func (n *emptyHwAccel) configure(mediaFile *ffmpegModels.Mediafile) error {
	return nil
}

func (n *emptyHwAccel) getType() string {
	return ""
}
