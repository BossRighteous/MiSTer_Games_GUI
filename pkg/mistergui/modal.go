package mistergui

type Modal interface {
	Label() string
	AcceptLabel() string
	RejectLabel() string
	OnAccept() bool
	OnReject() bool
	Render()
}

type BasicModal struct {
	label          string
	acceptLabel    string
	rejectLabel    string
	acceptCallback *func()
	rejectCallback *func()
	renderCallback *func()
}

func (modal *BasicModal) Label() string {
	return modal.label
}

func (modal *BasicModal) AcceptLabel() string {
	if modal.acceptCallback == nil {
		return ""
	}
	return modal.acceptLabel
}

func (modal *BasicModal) RejectLabel() string {
	if modal.rejectCallback == nil {
		return ""
	}
	return modal.rejectLabel
}

func (modal *BasicModal) OnAccept() {
	if modal.acceptCallback != nil {
		cb := *modal.acceptCallback
		cb()
	}
}

func (modal *BasicModal) OnReject() {
	if modal.rejectCallback != nil {
		cb := *modal.rejectCallback
		cb()
	}
}

func (modal *BasicModal) Render() {
	if modal.renderCallback != nil {
		cb := *modal.renderCallback
		cb()
	}
}
