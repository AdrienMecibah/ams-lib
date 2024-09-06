// +build linux darwin

package ams

import (
	"fmt"
)

import "reflect"


func SprintStruct(s interface{}) string {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)
	if val.Kind() != reflect.Struct {
		panic(fmt.Errorf("Provided value is not a struct"))
	}
	result := "{"
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		result += fmt.Sprintf("%s: %v", fieldType.Name, field.Interface())
		if i != val.NumField()-1 {
			result += ", "
		} else {
			result += "}"
		}
	}
	return result
}

type KeyState struct {
	RightAltPressed   bool
	LeftAltPressed    bool
	RightCtrlPressed  bool
	LeftCtrlPressed   bool
	ShiftPressed      bool
	NumLockOn         bool
	ScrollLockOn      bool
	CapsLockOn        bool
	EnhancedKey       bool
}

func KeyStateFromMask(mask uint32) KeyState {
	return KeyState {
		RightAltPressed:  mask&0x0001 != 0,
		LeftAltPressed:   mask&0x0002 != 0,
		RightCtrlPressed: mask&0x0004 != 0,
		LeftCtrlPressed:  mask&0x0008 != 0,
		ShiftPressed:     mask&0x0010 != 0,
		NumLockOn:        mask&0x0020 != 0,
		ScrollLockOn:     mask&0x0040 != 0,
		CapsLockOn:       mask&0x0080 != 0,
		EnhancedKey:      mask&0x0100 != 0,
	}
}

func MaskFromKeyState(state KeyState) uint32 {
	var mask uint32
	if state.RightAltPressed  { mask |= 0x0001 }
	if state.LeftAltPressed   { mask |= 0x0002 }
	if state.RightCtrlPressed { mask |= 0x0004 }
	if state.LeftCtrlPressed  { mask |= 0x0008 }
	if state.ShiftPressed     { mask |= 0x0010 }
	if state.NumLockOn        { mask |= 0x0020 }
	if state.ScrollLockOn     { mask |= 0x0040 }
	if state.CapsLockOn       { mask |= 0x0080 }
	if state.EnhancedKey      { mask |= 0x0100 }
	return mask
}

type KeyScope struct {
}

func (self *KeyScope) Init() {
	panic(fmt.Errorf("KeyScope is not implemented for linux builds"))
}

func (self *KeyScope) Test() {
	panic(fmt.Errorf("KeyScope is not implemented for linux builds"))
}

func (self *KeyScope) Read()  (char rune, code uint16, state KeyState) {
	panic(fmt.Errorf("KeyScope is not implemented for linux builds"))
	return rune(0), uint16(0), KeyState{}
}

func (self *KeyScope) ReadCombine() (char rune, code uint16, state KeyState) {
	panic(fmt.Errorf("KeyScope is not implemented for linux builds"))
	return rune(0), uint16(0), KeyState{}
}

func (self *KeyScope) ReadAll() (bytes [2]byte, keyState uint32, repetitions uint16, keyVirtualCode uint16, keyHardCode uint16) {
	panic(fmt.Errorf("KeyScope is not implemented for linux builds"))
	return [2]byte{0, 0}, uint32(0), uint16(0), uint16(0), uint16(0) 
}

func NewKeyScope() *KeyScope {
	panic(fmt.Errorf("KeyScope is not implemented for linux builds"))
	result := KeyScope{}
	return &result
}