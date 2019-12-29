package entities

import "testing"

func TestEventID_New(t *testing.T) {
	uid, err := NewEventID("")
	if err != nil {
		t.Error(err)
	}

	if uid.String() == "" {
		t.Fail()
	}
}

// TestEventID_Equal_Different проверяет корректно ли работает сравнение разных значений
func TestEventID_Equal_Different(t *testing.T) {
	uid1, err := NewEventID("")
	if err != nil {
		t.Error(err)
	}

	uid2, err := NewEventID("")
	if err != nil {
		t.Error(err)
	}

	if uid1.Equal(uid2) == true {
		t.Fail()
	}
}

// TestEventID_Equal_Similar проверяет корректно ли работает сравнение одинаковых значений
func TestEventID_Equal_Similar(t *testing.T) {
	uid1, err := NewEventID("")
	if err != nil {
		t.Error(err)
	}

	// Создаст значение с тем же самым ID, что в uid1
	uid2, err := NewEventID(uid1.String())
	if err != nil {
		t.Error(err)
	}

	if uid1.Equal(uid2) == false {
		t.Fail()
	}
}
