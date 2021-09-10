// SPDX-License-Identifier: Unlicense OR MIT

package text

import (
	"io"
	"strings"

	"golang.org/x/image/math/fixed"

	"gioui.org/op"
)

// Shaper implements layout and shaping of text.
type Shaper interface {
	// Layout a text according to a set of options.
	Layout(font Font, size fixed.Int26_6, maxWidth int, txt io.Reader) ([]Line, error)
	// LayoutString is Layout for strings.
	LayoutString(font Font, size fixed.Int26_6, maxWidth int, str string) []Line
	// Shape a line of text and return a clipping operation for its outline.
	Shape(font Font, size fixed.Int26_6, layout Layout) op.CallOp
}

// A FontFace is a Font and a matching Face.
type FontFace struct {
	Font Font
	Face Face
}

// Cache implements cached layout and shaping of text from a set of
// registered fonts.
//
// If a font matches no registered shape, Cache falls back to the
// first registered face.
//
// The LayoutString and ShapeString results are cached and re-used if
// possible.
type Cache struct {
	def   Typeface
	faces map[Font]*faceCache
}

type faceCache struct {
	face        Face
	layoutCache layoutCache
	pathCache   pathCache
}

func (c *Cache) lookup(font Font) *faceCache {
	f := c.faceForStyle(font)
	if f == nil {
		font.Typeface = c.def
		f = c.faceForStyle(font)
	}
	return f
}

func (c *Cache) faceForStyle(font Font) *faceCache {
	if closest, ok := c.closestFont(font); ok {
		return c.faces[closest]
	}
	font.Style = Regular
	if closest, ok := c.closestFont(font); ok {
		return c.faces[closest]
	}
	return nil
}

// closestFont returns the closest Font by weight, in case of equality the
// lighter weight will be returned.
func (c *Cache) closestFont(lookup Font) (Font, bool) {
	if c.faces[lookup] != nil {
		return lookup, true
	}
	found := false
	var match Font
	for cf := range c.faces {
		if cf.Typeface != lookup.Typeface || cf.Variant != lookup.Variant || cf.Style != lookup.Style {
			continue
		}
		if !found {
			found = true
			match = cf
			continue
		}
		cDist := weightDistance(lookup.Weight, cf.Weight)
		mDist := weightDistance(lookup.Weight, match.Weight)
		if cDist < mDist {
			match = cf
		} else if cDist == mDist && cf.Weight < match.Weight {
			match = cf
		}
	}
	return match, found
}

func NewCache(collection []FontFace) *Cache {
	c := &Cache{
		faces: make(map[Font]*faceCache),
	}
	for i, ff := range collection {
		if i == 0 {
			c.def = ff.Font.Typeface
		}
		c.faces[ff.Font] = &faceCache{face: ff.Face}
	}
	return c
}

// Layout implements the Shaper interface.
func (c *Cache) Layout(font Font, size fixed.Int26_6, maxWidth int, txt io.Reader) ([]Line, error) {
	cache := c.lookup(font)
	return cache.face.Layout(size, maxWidth, txt)
}

// LayoutString is a caching implementation of the Shaper interface.
func (c *Cache) LayoutString(font Font, size fixed.Int26_6, maxWidth int, str string) []Line {
	cache := c.lookup(font)
	return cache.layout(size, maxWidth, str)
}

// Shape is a caching implementation of the Shaper interface. Shape assumes that the layout
// argument is unchanged from a call to Layout or LayoutString.
func (c *Cache) Shape(font Font, size fixed.Int26_6, layout Layout) op.CallOp {
	cache := c.lookup(font)
	return cache.shape(size, layout)
}

func (f *faceCache) layout(ppem fixed.Int26_6, maxWidth int, str string) []Line {
	if f == nil {
		return nil
	}
	lk := layoutKey{
		ppem:     ppem,
		maxWidth: maxWidth,
		str:      str,
	}
	if l, ok := f.layoutCache.Get(lk); ok {
		return l
	}
	l, _ := f.face.Layout(ppem, maxWidth, strings.NewReader(str))
	f.layoutCache.Put(lk, l)
	return l
}

func (f *faceCache) shape(ppem fixed.Int26_6, layout Layout) op.CallOp {
	if f == nil {
		return op.CallOp{}
	}
	pk := pathKey{
		ppem: ppem,
		str:  layout.Text,
	}
	if clip, ok := f.pathCache.Get(pk); ok {
		return clip
	}
	clip := f.face.Shape(ppem, layout)
	f.pathCache.Put(pk, clip)
	return clip
}
