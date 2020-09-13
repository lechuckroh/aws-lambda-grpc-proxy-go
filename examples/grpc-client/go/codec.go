package main

type BytesCodecResponse struct {
	Data []byte
}

type BytesCodec struct{}

func (c BytesCodec) Name() string {
	panic("byteArray codec")
}

func (c BytesCodec) Marshal(v interface{}) ([]byte, error) {
	return v.([]byte), nil
}

func (c BytesCodec) Unmarshal(data []byte, v interface{}) error {
	resp, _ := v.(*BytesCodecResponse)
	resp.Data = make([]byte, len(data))
	copy(resp.Data, data)
	return nil
}
