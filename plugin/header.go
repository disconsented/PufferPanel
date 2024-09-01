package plugin

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const FixedLength = 14

type header struct {
	// fixed header 14b
	sequence    uint32 // 4
	topicLength uint16 // 2
	bodyLength  uint32 // 4
	reserved    uint32 // 4
	// variable parts
	topic string
	body  []byte
}

func read(reader io.Reader) (*header, error) {
	header := header{}
	buf := make([]byte, 14)
	i, err := reader.Read(buf)
	if err != nil {
		return &header, err
	}
	if i != FixedLength {
		return nil, errors.New("invalid header length")
	}
	fmt.Printf("read: %+v\n", buf)
	header.sequence = binary.BigEndian.Uint32(buf[0:4])
	header.topicLength = binary.BigEndian.Uint16(buf[4:6])
	header.bodyLength = binary.BigEndian.Uint32(buf[6:10])
	header.reserved = binary.BigEndian.Uint32(buf[10:14])
	topic := make([]byte, uint32(header.topicLength))
	bytes, err := reader.Read(topic)
	if err != nil {
		return nil, err
	}
	if bytes != int(header.topicLength) {
		return nil, errors.New("invalid topic length")
	}
	header.topic = string(topic)
	header.body = make([]byte, header.bodyLength)
	bytes, err = reader.Read(header.body)
	if err != nil {
		return nil, err
	}
	if bytes != int(header.bodyLength) {
		return nil, errors.New("invalid body length")
	}
	return &header, err
}

func (h *header) write(writer io.Writer) error {
	buf := make([]byte, 14)
	binary.BigEndian.PutUint32(buf[0:4], h.sequence)
	binary.BigEndian.PutUint16(buf[4:6], h.topicLength)
	binary.BigEndian.PutUint32(buf[6:10], h.bodyLength)
	binary.BigEndian.PutUint32(buf[10:14], h.reserved)
	buf = append(buf, h.topic...)
	buf = append(buf, h.body...)
	fmt.Printf("%+v\n", buf)
	write, err := writer.Write(buf)
	if err != nil {
		return err
	}
	println("wrote bytes", write)
	return nil
}
