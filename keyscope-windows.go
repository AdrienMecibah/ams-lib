// +build windows

package ams

import (
	"fmt"
	"syscall"
	"unsafe"
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
	reader *syscall.LazyProc
	handle uintptr
}

func (self *KeyScope) Init() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	self.reader = kernel32.NewProc("ReadConsoleInputW")
	self.handle, _, _ = kernel32.NewProc("GetStdHandle").Call(uintptr(-10&0xFFFFFFFF))
	var mode uint32
	kernel32.NewProc("GetConsoleMode").Call(self.handle, uintptr(unsafe.Pointer(&mode)))
	kernel32.NewProc("SetConsoleMode").Call(self.handle, uintptr(mode&^(0x0004|0x0002)|0x0001))
}

func (self *KeyScope) Test() {
	for {
		r, c, s := self.Read()
		fmt.Printf("rune = '%s'\ncode = 0x%X\nstate = %s\n\n", string(r), c, SprintStruct(s))
	}
}

func (self *KeyScope) Read()  (char rune, code uint16, state KeyState) {
	bytePair, keyState, _, keyVirtualCode, _ := self.ReadAll()
	bytes := []byte{bytePair[0], bytePair[1]}
	return ([]rune(string(bytes)))[0], keyVirtualCode, KeyStateFromMask(keyState)
}

func (self *KeyScope) ReadCombine() (char rune, code uint16, state KeyState) {
	for {
		char, code, state = self.Read()
		if int(char) > 0 || state.EnhancedKey{
			return char, code, state
		}
	}
	return 0, 0, KeyState{}
}

func (self *KeyScope) ReadAll() (bytes [2]byte, keyState uint32, repetitions uint16, keyVirtualCode uint16, keyHardCode uint16) {
	type InputRecord struct {
		EventType uint16
		_         uint16
		Event     [16]byte
	}
	type EventRecord struct {
		BKeyDown          int32
		BRepeatCount      uint16
		WVirtualKeyCode   uint16
		WVirtualScanCode  uint16
		UChar             [2]byte
		DwControlKeyState uint32
	}	
	var input InputRecord
	for {
		var numRead uint32
		self.reader.Call(self.handle, uintptr(unsafe.Pointer(&input)), 1, uintptr(unsafe.Pointer(&numRead)))
		if input.EventType == 0x0001 {
			ev := (*EventRecord)(unsafe.Pointer(&input.Event))
			if ev.BKeyDown != 0 { 
				return ev.UChar, ev.DwControlKeyState, ev.BRepeatCount, ev.WVirtualKeyCode, ev.WVirtualScanCode
			}
		}
	}
	return [2]byte{0, 0}, 0, 0, 0, 0
}

func NewKeyScope() *KeyScope {
	result := KeyScope{}
	result.Init()
	return &result
}