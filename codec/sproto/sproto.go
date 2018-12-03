package sprotocodec

import (
	"reflect"

	"github.com/BobWong/cellnet"
	"github.com/BobWong/gosproto"
	"fmt"
	"path"
	"github.com/BobWong/cellnet/util"
	"github.com/BobWong/cellnet/codec"
)

type sprotoCodec struct {
}

func (self *sprotoCodec) Name() string {
	return "sproto"
}

func (self *sprotoCodec) MimeType() string {
	return "application/sproto"
}

func (self *sprotoCodec) Encode(msgObj interface{},ctx cellnet.ContextSet) (data interface{}, err error) {

	result, err := sproto.Encode(msgObj)
	if err != nil {
		return nil, err
	}

	return sproto.Pack(result), nil
}

func (self *sprotoCodec) Decode(data interface{}, msgObj interface{}) error {
	tmp := data.([]byte)
	// sproto要求必须有头, 但空包也是可以的
	if len(tmp) == 0 {
		return nil
	}

	raw, err := sproto.Unpack(tmp)
	if err != nil {
		return err
	}

	_, err2 := sproto.Decode(raw, msgObj)

	return err2
}

func AutoRegisterMessageMeta(msgTypes []reflect.Type) {

	for _, tp := range msgTypes {

		msgName := fmt.Sprintf("%s.%s", path.Base(tp.PkgPath()), tp.Name())

		//cellnet.RegisterMessageMeta("sproto", msgName, tp, util.StringHash(msgName))
		cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
			Codec:codec.MustGetCodec("sproto"),
			Type: tp,
			ID: int(util.StringHash(msgName)),
		})
	}

}

func init() {

	codec.RegisterCodec(new(sprotoCodec))
}