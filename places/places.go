package places

import (
	"bytes"
	"io"
)

var (
	startDel = []byte("<@")
	endDel   = []byte("@>")
)

type Template struct {
	places   []int
	template []byte
}

func NewTemplate(t []byte) *Template {
	return &Template{template: t, places: Find(t)}
}

func (t *Template) Replace(bf Buffer, replacements map[string]io.ReadSeeker) {
	Replace(t.template, bf, t.places, replacements)
}

func (t *Template) ReplaceString(bf Buffer, replacements map[string]string) {
	ReplaceString(t.template, bf, t.places, replacements)
}

func (t *Template) ReplaceBytes(bf Buffer, replacements map[string][]byte) {
	ReplaceBytes(t.template, bf, t.places, replacements)
}

func (t *Template) ReplaceMapper(bf Buffer, mapper Mapper) {
	ReplaceMapper(t.template, bf, t.places, mapper)
}

// Find looks for placeholders written in the style "<@placeholdername@>" inside the given template.
// It returns a slice containing the positions of the placeholders that is meant to be passed to
// Replace or ReplaceString.
func Find(template []byte) (places []int) {

	places = make([]int, 0, 22)

	var (
		found  = -1
		start  int
		end    int
		length = len(template)
	)

	for i := 0; i < length; i++ {

		found = bytes.Index(template[i:], startDel)
		if -1 == found {
			break
		}
		start = found + i

		found = bytes.Index(template[start+2:], endDel)
		if -1 == found {
			break
		}
		end = found + start + 2 // two bytes for each delimiter

		places = append(places, start, end)
		i = end

	}

	return
}

// Replace replaces the placeholders at the given places inside the template with
// the replacements found inside the map and writes the result to the buffer.
// The given template must be the unchanged byte array that was passed to Find in order to get the
// places. For strings as replacements see the optimized ReplaceString function for bytes use ReplaceBytes.
func Replace(template []byte, bf Buffer, places []int, replacements map[string]io.ReadSeeker) {
	var (
		last        int
		first       int
		has         bool
		replacement io.ReadSeeker
		length      = len(places)
	)

	// we iterate over places always taking pairs of ints
	// where the first int is the starting and the last is the ending position
	// i.e.
	// "the quick <@colorOfFox@> fox"
	//            |           |
	//           places[i]   places[i+1]
	// so places[i] == 10 and places[i+1] == 22
	// so we can get:
	//   - everything before the place with
	//       template[:places[i]]
	//   - the placeholder name with
	//       template[places[i]+2:places[i+1]]
	//   - everything after the place with
	//       template[places[i+1]+2:]
	//
	// instead of going just from the beginning to the end, we iterate over a cursor (last)
	// from placeholder to placeholder until we are through the template
	for i := 0; i < length; i += 2 {
		// track the first position of the placeholder
		first = places[i]

		// take the bytes from the last position within the template up to the placeholder
		bf.Write(template[last:first])

		// lookup the placeholder name within the replacements and
		// write the replacement if we found one
		replacement, has = replacements[string(template[first+2:places[i+1]])]
		if has {
			replacement.Seek(0, 0)
			io.Copy(bf, replacement)
		}

		// track the last position for the next iteration
		last = places[i+1] + 2
	}

	bf.Write(template[last:]) // write any remaining parts of the template that don't have any placeholders
}

// The Buffer interface is fullfilled by *bytes.Buffer. However since for performance reasons
// the errors from writing to the buffer are always ignored, a buffer wrapper is needed to capture them.
type Buffer interface {
	io.Writer
	WriteString(string) (int, error)
}

// ReplaceString replaces the placeholders at the given places inside the template with
// the replacements found inside the map and writes the result to the buffer.
// The given template must be the unchanged byte array that was passed to Find in order to get the
// places.
func ReplaceString(template []byte, bf Buffer, places []int, replacements map[string]string) {
	var (
		last        int
		first       int
		has         bool
		replacement string
		length      = len(places)
	)

	// we iterate over places always taking pairs of ints
	// where the first int is the starting and the last is the ending position
	// i.e.
	// "the quick <@colorOfFox@> fox"
	//            |           |
	//           places[i]   places[i+1]
	// so places[i] == 10 and places[i+1] == 22
	// so we can get:
	//   - everything before the place with
	//       template[:places[i]]
	//   - the placeholder name with
	//       template[places[i]+2:places[i+1]]
	//   - everything after the place with
	//       template[places[i+1]+2:]
	//
	// instead of going just from the beginning to the end, we iterate over a cursor (last)
	// from placeholder to placeholder until we are through the template
	for i := 0; i < length; i += 2 {
		// track the first position of the placeholder
		first = places[i]

		// take the bytes from the last position within the template up to the placeholder
		bf.Write(template[last:first])

		// track the last position for the next iteration

		// lookup the placeholder name within the replacements and
		// write the replacement if we found one
		replacement, has = replacements[string(template[first+2:places[i+1]])]
		if has {
			bf.WriteString(replacement)
		}

		last = places[i+1] + 2
	}

	bf.Write(template[last:]) // write any remaining parts of the template that don't have any placeholders
}

// ReplaceBytes replaces the placeholders at the given places inside the template with
// the replacements found inside the map and writes the result to the buffer.
// The given template must be the unchanged byte array that was passed to Find in order to get the
// places.
func ReplaceBytes(template []byte, bf Buffer, places []int, replacements map[string][]byte) {
	var (
		last        int
		first       int
		has         bool
		replacement []byte
		length      = len(places)
	)

	// we iterate over places always taking pairs of ints
	// where the first int is the starting and the last is the ending position
	// i.e.
	// "the quick <@colorOfFox@> fox"
	//            |           |
	//           places[i]   places[i+1]
	// so places[i] == 10 and places[i+1] == 22
	// so we can get:
	//   - everything before the place with
	//       template[:places[i]]
	//   - the placeholder name with
	//       template[places[i]+2:places[i+1]]
	//   - everything after the place with
	//       template[places[i+1]+2:]
	//
	// instead of going just from the beginning to the end, we iterate over a cursor (last)
	// from placeholder to placeholder until we are through the template
	for i := 0; i < length; i += 2 {
		// track the first position of the placeholder
		first = places[i]

		// take the bytes from the last position within the template up to the placeholder
		bf.Write(template[last:first])

		// track the last position for the next iteration

		// lookup the placeholder name within the replacements and
		// write the replacement if we found one
		replacement, has = replacements[string(template[first+2:places[i+1]])]
		if has {
			bf.Write(replacement)
		}

		last = places[i+1] + 2
	}

	bf.Write(template[last:]) // write any remaining parts of the template that don't have any placeholders
}

// Mapper maps strings.
type Mapper interface {
	// Map maps a placeholder to a value. It should return the empty string, if
	// a value could not be found
	Map(string) string
}

// ReplaceMapper replaces the placeholders at the given places inside the template with
// the replacements returned from the mapper and writes the result to the buffer.
// The given template must be the unchanged byte array that was passed to Find in order to get the
// places.
func ReplaceMapper(template []byte, bf Buffer, places []int, mapper Mapper) {
	var (
		last        int
		first       int
		replacement string
		length      = len(places)
	)

	// we iterate over places always taking pairs of ints
	// where the first int is the starting and the last is the ending position
	// i.e.
	// "the quick <@colorOfFox@> fox"
	//            |           |
	//           places[i]   places[i+1]
	// so places[i] == 10 and places[i+1] == 22
	// so we can get:
	//   - everything before the place with
	//       template[:places[i]]
	//   - the placeholder name with
	//       template[places[i]+2:places[i+1]]
	//   - everything after the place with
	//       template[places[i+1]+2:]
	//
	// instead of going just from the beginning to the end, we iterate over a cursor (last)
	// from placeholder to placeholder until we are through the template
	for i := 0; i < length; i += 2 {
		// track the first position of the placeholder
		first = places[i]

		// take the bytes from the last position within the template up to the placeholder
		bf.Write(template[last:first])

		// track the last position for the next iteration

		// lookup the placeholder name within the replacements and
		// write the replacement if we found one
		replacement = mapper.Map(string(template[first+2 : places[i+1]]))
		if len(replacement) > 0 {
			bf.WriteString(replacement)
		}

		last = places[i+1] + 2
	}

	bf.Write(template[last:]) // write any remaining parts of the template that don't have any placeholders
}

// FindAndReplace finds placeholders and replaces them in one go.
func FindAndReplace(template []byte, bf Buffer, replacements map[string]io.ReadSeeker) {
	Replace(template, bf, Find(template), replacements)
}

// FindAndReplaceString finds placeholders and replaces them in one go.
func FindAndReplaceString(template []byte, bf Buffer, replacements map[string]string) {
	ReplaceString(template, bf, Find(template), replacements)
}

// FindAndReplaceBytes finds placeholders and replaces them in one go.
func FindAndReplaceBytes(template []byte, bf Buffer, replacements map[string][]byte) {
	ReplaceBytes(template, bf, Find(template), replacements)
}

// FindAndReplaceMapper finds placeholders and replaces them in one go.
func FindAndReplaceMapper(template []byte, bf Buffer, mapper Mapper) {
	ReplaceMapper(template, bf, Find(template), mapper)
}
