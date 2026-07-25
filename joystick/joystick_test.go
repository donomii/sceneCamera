package joystick

import "testing"

func TestSetupJoystick(t *testing.T) {
	joystick_type = "F310"
	Setup_joystick()
	if joystick_type != "F310" && joystick_type != "steamdeck" {
		t.Errorf("expected F310 or steamdeck joystick type, got %q", joystick_type)
	}
}
