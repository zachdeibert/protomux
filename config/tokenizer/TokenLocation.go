package tokenizer

import "fmt"

// TokenLocation represents the location in a config file at which a token resides
type TokenLocation struct {
	FileName  string
	Line      []byte
	LineNo    int
	CharStart int
	CharLen   int
}

// ShortString describes the location in a shorter manner than String()
func (t TokenLocation) ShortString() string {
	return fmt.Sprintf("%s:%d:%d", t.FileName, t.LineNo+1, t.CharStart+1)
}

func (t TokenLocation) String() string {
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
