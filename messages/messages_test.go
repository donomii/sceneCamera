package messages

import "testing"

func TestRegisterAndUnregister(t *testing.T) {
	messageName := t.Name()
	Register(messageName, "first", nil)
	Register(messageName, "second", nil)

	handlers, found := MessageRegistry.Get(messageName)
	if !found {
		t.Fatalf("expected handlers for %q", messageName)
	}
	if _, found := handlers.Get("first"); !found {
		t.Error("expected first handler")
	}
	if _, found := handlers.Get("second"); !found {
		t.Error("expected second handler")
	}

	Unregister(messageName, "first")
	if _, found := handlers.Get("first"); found {
		t.Error("expected first handler to be removed")
	}
}

func TestSendMessageWithoutHandlers(t *testing.T) {
	SendMessage(t.Name(), 42)
}
