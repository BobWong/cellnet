// Generated by github.com/bobwong89757/cellnet/protoc-gen-msg
// DO NOT EDIT!
// Source: pb.proto

package test

import (
	"github.com/bobwong89757/cellnet"
	"reflect"
	_ "github.com/bobwong89757/cellnet/codec/gogopb"
	"github.com/bobwong89757/cellnet/codec"
)

func init() {

	// pb.proto
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*ContentACK)(nil)).Elem(),
		ID:    60952,
	})
}
