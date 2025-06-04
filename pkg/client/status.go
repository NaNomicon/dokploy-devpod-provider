package client

// Status represents the status of a machine as expected by DevPod
type Status string

const (
	// StatusRunning indicates the machine is running and ready
	StatusRunning Status = "Running"
	
	// StatusBusy indicates the machine is doing something and DevPod should wait
	StatusBusy Status = "Busy"
	
	// StatusStopped indicates the machine is currently stopped
	StatusStopped Status = "Stopped"
	
	// StatusNotFound indicates the machine is not found
	StatusNotFound Status = "NotFound"
)

// String returns the string representation of the status
func (s Status) String() string {
	return string(s)
} 