package tomodoro

// MessageType is the type of message received from the websocket
type MessageType string

// Message Types
const (
	Tick         MessageType = "tick"
	TimerStopped MessageType = "timerStopped"
	TimerStarted MessageType = "timerStarted"
	Connecting   MessageType = "connecting"  // Only Used for internal purposes
	Connected    MessageType = "connected"   // Only Used for internal purposes
	Error        MessageType = "error"       // Only Used for internal purposes
	Terminating  MessageType = "terminating" // Only Used for internal purposes
)

// Message is the message received from the websocket
type Message struct {
	Type    MessageType `json:"type"`
	Payload struct {
		Name          string `json:"name"`
		RemainingTime int64  `json:"remainingTime"`
		Team          string `json:"team"`
		Timestamp     int64  `json:"timestamp"`
		Duration      int64  `json:"duration"`
	} `json:"payload"`
	Error error // Only Used for internal purposes
}

// IsTick returns true if the message is a tick message
func (m *Message) IsTick() bool {
	return m.Type == Tick
}

// IsTimerStopped returns true if the message is a timer stopped message
func (m *Message) IsTimerStopped() bool {
	return m.Type == TimerStopped
}

// IsTimerStarted returns true if the message is a timer started message
func (m *Message) IsTimerStarted() bool {
	return m.Type == TimerStarted
}

// IsConnecting returns true if the message is a connecting message
func (m *Message) IsConnecting() bool {
	return m.Type == Connecting
}

// IsConnected returns true if the message is a connected message
func (m *Message) IsConnected() bool {
	return m.Type == Connected
}

// IsError returns true if the message is an error
func (m *Message) IsError() bool {
	return m.Type == Error
}
