package tomodoro

type MessageType string

const (
	Tick         MessageType = "tick"
	TimerStopped MessageType = "timerStopped"
	TimerStarted MessageType = "timerStarted"
	Connecting   MessageType = "connecting" // Only Used for internal purposes
	Listening    MessageType = "listening"  // Only Used for internal purposes
	Error        MessageType = "error"      // Only Used for internal purposes
	Terminated   MessageType = "terminated" // Only Used for internal purposes
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
