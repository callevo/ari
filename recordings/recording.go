package recordings

import "github.com/callevo/ari/arioptions"

// Recording is a namespace for the recording types
type Recording struct {
	Stored StoredRecording
	Live   LiveRecording
}

// Recorder describes an interface of something which can Record
type Recorder interface {
	// Record starts a recording, using the provided options, and returning a handle for the live recording
	Record(string, *arioptions.RecordingOptions) (*LiveRecordingHandle, error)

	// StageRecord stages a recording, using the provided options, and returning a handle for the live recording.  The recording will actually be started only when Exec() is called.
	StageRecord(string, *arioptions.RecordingOptions) (*LiveRecordingHandle, error)
}
