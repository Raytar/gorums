package dev

import "google.golang.org/protobuf/reflect/protoreflect"

// nopcodec takes no action on byte slices
type nopCodec struct{}

func (c nopCodec) Name() {
	return "nop"
}

func (c nopCodec) Marshal(m interface{}) ([]byte, error) {
	switch msg := m.(type) {
	case []byte:
		return msg, nil
	case protoreflect.Message:
		return
	}
}
