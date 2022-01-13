package gremlingo

import (
	"bytes"
	"encoding/binary"
)

type graphBinMessageSerializer struct {
	readerClass graphBinaryReader
	writerClass graphBinaryWriter
	memeType    string `default:"application/vnd.graphbinary-v1.0"`
}

const versionByte byte = 0x81

func (gs *graphBinMessageSerializer) getProcessor(processor string) string {
	return processor
}

func (gs *graphBinMessageSerializer) serializeMessage(request *request) []byte {
	finalMessage, _ := gs.buildMessage(request, 0x20, gs.memeType)
	return finalMessage
}

func (gs *graphBinMessageSerializer) buildMessage(message *request, mimeLen byte, mimeType string) ([]byte, error) {
	byteBuffer := bytes.Buffer{}

	// header
	byteBuffer.WriteByte(mimeLen)
	byteBuffer.WriteString(mimeType)
	byteBuffer.WriteByte(versionByte)

	// RequestID
	binary.Write(&byteBuffer, binary.BigEndian, message.RequestId)

	// Op
	gs.writerClass.writeValue(message.Op, &byteBuffer, false)
	//binary.Write(&byteBuffer, binary.BigEndian, len(message.Op))
	//binary.Write(&byteBuffer, binary.BigEndian, message.Op)

	// Processor
	gs.writerClass.writeValue(message.Processor, &byteBuffer, false)
	//binary.Write(&byteBuffer, binary.BigEndian, len(message.Processor))
	//binary.Write(&byteBuffer, binary.BigEndian, message.Processor)

	// Args
	args := message.Args
	binary.Write(&byteBuffer, binary.BigEndian, len(args))
	//binary.Write(&byteBuffer, binary.BigEndian, args)
	for k, v := range args {
		gs.writerClass.writeObject(k, &byteBuffer)
		gs.writerClass.writeObject(v, &byteBuffer)
	}

	return byteBuffer.Bytes(), nil
}

func (gs *graphBinMessageSerializer) deserializeMessage(message interface{}) map[string]interface{} {
	//msg := map[string]interface{}{
	//	"requestId": request_id,
	//	"status": map[string]interface{}{"code": status_code,
	//		"message":    status_msg,
	//		"attributes": status_attrs},
	//	"result": map[string]interface{}{"meta": meta_attrs,
	//		"data": result}}
	return nil
}
