package tokenizer

// Token represents a single token from the config file
type Token struct {
	Type      TokenType
	Value     string
	Line      []byte
	LineNo    int
	CharStart int
	CharLen   int
}

func (t Token) String() string {
	buf := make([]byte, len(t.Line)+1+t.CharStart+t.CharLen)
	copy(buf, t.Line)
	bufI := len(t.Line)
	buf[bufI] = '\n'
	bufI++
	for i := 0; i < t.CharStart; i++ {
		buf[bufI] = ' '
		bufI++
	}
	for i := 0; i < t.CharLen; i++ {
		buf[bufI] = '^'
		bufI++
	}
	return string(buf)
}
