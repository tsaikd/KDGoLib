package httpRender

func (t *renderImpl) GetBody() []byte {
	return t.buffer.Bytes()
}

func (t *renderImpl) GetSize() int64 {
	return t.size
}
