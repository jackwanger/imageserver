package native

import "github.com/pierrre/imageserver"

var _ imageserver.Server = &Server{}

var _ Processor = ProcessorFunc(nil)

var _ Encoder = EncoderFunc(nil)
