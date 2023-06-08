package tomodoro

type MessageType string

const (
	Tick         MessageType = "tick"
	TimerStopped MessageType = "timerStopped"
	TimerStarted MessageType = "timerStarted"
	Connecting   MessageType = "connecting"  // Only Used for internal purposes
	Connected    MessageType = "connected"   // Only Used for internal purposes
	Error        MessageType = "error"       // Only Used for internal purposes
	Terminating  MessageType = "terminating" // Only Used for internal purposes
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload struct {
		Name          string `json:"name"`
		RemainingTime int64  `json:"remainingTime"`
		Team          string `json:"team"`
		Timestamp     int64  `json:"timestamp"`
	} `json:"payload"`
	Error string // Only Used for internal purposes
}

func (m *Message) IsTick() bool {
	return m.Type == Tick
}

func (m *Message) IsTimerStopped() bool {
	return m.Type == TimerStopped
}

func (m *Message) IsTimerStarted() bool {
	return m.Type == TimerStarted
}
func (m *Message) IsConnecting() bool {
	return m.Type == Connecting
}
func (m *Message) IsConnected() bool {
	return m.Type == Connected
}
func (m *Message) IsError() bool {
	return m.Type == Error
}
