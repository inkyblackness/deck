package events

// FileDropEventType is the name for events where one or more files are dropped into the user interface.
const FileDropEventType = EventType("env.filedrop")

// FileDropEvent is used to inform about dropped files.
type FileDropEvent struct {
	eventType EventType

	x, y      float32
	filePaths []string
}

// NewFileDropEvent initializes a basic clipboard event structure.
func NewFileDropEvent(x, y float32, filePaths []string) *FileDropEvent {
	event := &FileDropEvent{
		eventType: FileDropEventType,
		x:         x,
		y:         y,
		filePaths: filePaths}

	return event
}

// EventType implements the Event interface.
func (event *FileDropEvent) EventType() EventType {
	return event.eventType
}

// Position returns the coordinate of the event.
func (event *FileDropEvent) Position() (x, y float32) {
	return event.x, event.y
}

// FilePaths returns the list of file paths
func (event *FileDropEvent) FilePaths() []string {
	return event.filePaths
}
