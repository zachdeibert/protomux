package common

import "fmt"

// Location represents the location in a config file at which a string resides
type Location struct {
	FileName  string
	Line      []byte
	LineNo    int
	CharStart int
	CharLen   int
}

// ShortString describes the location in a shorter manner than String()
func (t Location) ShortString() string {
	return fmt.Sprintf("%s:%d:%d", t.FileName, t.LineNo+1, t.CharStart+1)
}

func (t Location) String() string {
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

// Merge multiple locations into one
func Merge(locations []Location) Location {
	if len(locations) == 0 {
		return Location{
			FileName:  "/dev/null",
			Line:      []byte{},
			LineNo:    0,
			CharStart: 0,
			CharLen:   0,
		}
	}
	res := locations[0]
	loc := locations[len(locations)-1]
	if loc.FileName != res.FileName || loc.LineNo != res.LineNo {
		res.CharLen = len(res.Line) - res.CharStart
	} else {
		res.CharLen = loc.CharStart + loc.CharLen - res.CharStart
	}
	return res
}
