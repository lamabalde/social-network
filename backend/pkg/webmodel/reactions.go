package webmodel

type Reaction struct {
	MessageType string
	MessageID   string
	Reaction    bool
}

func (r *Reaction) Validate() string {
	if IsEmpty(r.MessageType) {
		return "type of message missing"
	}
	return ""
}
