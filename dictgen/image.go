package dictgen

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"image"
	"io"
	"math"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/geek1011/dictutil/kobodict"
)

// ImageHandler transforms images referenced in a DictFile.
type ImageHandler interface {
	// Transform transforms an image read from ir, and returns a new value for
	// the img tag's src attribute. As a special case, if an empty string is
	// returned and the error is nil, the image tag is removed entirely. In
	// addition, custom CSS (which must not contain any double quotes) can be
	// returned to be set on the img tag.
	Transform(src string, ir io.Reader, dw *kobodict.Writer) (nsrc string, css string, err error)
}

// ImageHandlerRemove removes images from the dicthtml.
type ImageHandlerRemove struct{}

// Transform implements ImageHandler.
func (*ImageHandlerRemove) Transform(string, io.Reader, *kobodict.Writer) (string, string, error) {
	return "", "", nil
}

// ImageHandlerEmbed adds the images to the dictzip without any additional
// modifications. Usually, this would be the best choice, but unfortunately,
// it is too buggy as of firmware 4.19.14123.
type ImageHandlerEmbed struct{}

// Transform implements ImageHandler.
func (*ImageHandlerEmbed) Transform(src string, ir io.Reader, dw *kobodict.Writer) (string, string, error) {
	if !strings.HasSuffix(src, ".jpg") && !strings.HasSuffix(src, ".gif") {
		return "", "", fmt.Errorf("ImageHandlerEmbed: unsupported image file %s: extension must be .jpg or .gif when embedding", src)
	}

	// to generate a deterministic usually-unique filename
	fn := fmt.Sprintf("%x%s", sha1.Sum([]byte(src)), filepath.Ext(src))
	if !dw.Exists(fn) { // CreateFile will error if it already exists, and we're pretty confident the file is identical anyways
		if iw, err := dw.CreateFile(fn); err != nil {
			return "", "", fmt.Errorf("ImageHandlerEmbed: create dictfile entry %#v: %w", fn, err)
		} else if _, err := io.Copy(iw, ir); err != nil {
			return "", "", fmt.Errorf("ImageHandlerEmbed: copy image to dictfile: %w", err)
		}
	}
	return "dict:///" + fn, "", nil
}

// ImageHandlerBase64 optimizes the image and encodes it as base64. This is the
// most compatible option, but it comes at the expense of space and speed. In
// addition, if there are too many images, it can lead to nickel running out of
// memory when parsing the dictionary (and sickel should reboot it).
//
// In addition, it adds CSS to fix sizing issues (by default, images appear
// really small when rendered in the dictionary due to default styling).
//
// This is currently the recommended option for adding images.
//
// You must import image/* yourself for format support.
type ImageHandlerBase64 struct {
	// Images will be resized to fit within these dimensions, while preserving
	// aspect ratio. If not specified, the default is 1000x1000.
	MaxSize image.Point
	// NoGrayscale will prevent images from being grayscaled.
	NoGrayscale bool
	// JPEGQuality sets the JPEG quality for the encoded images. If not set, it
	// defaults to 60.
	JPEGQuality int
}

// Transform implements ImageHandler.
func (ih *ImageHandlerBase64) Transform(src string, ir io.Reader, dw *kobodict.Writer) (string, string, error) {
	img, err := imaging.Decode(ir)
	if err != nil {
		return "", "", fmt.Errorf("ImageHandlerBase64: decode image: %w", err)
	}

	// resize it
	mw, mh := float64(ih.MaxSize.X), float64(ih.MaxSize.Y)
	if mw < 1 {
		mw = 1000
	}
	if mh < 1 {
		mh = 1000
	}
	ow, oh := float64(img.Bounds().Dx()), float64(img.Bounds().Dy())
	sf := math.Min(mw/ow, mh/oh)
	nw, nh := ow*sf, oh*sf
	img = imaging.Resize(img, int(nw), int(nh), imaging.Lanczos)

	// make it grayscale
	if !ih.NoGrayscale {
		img = imaging.Grayscale(img)
	}

	// set the quality
	jq := ih.JPEGQuality
	if jq == 0 {
		jq = 60
	}

	// encode the image
	buf := bytes.NewBuffer(nil)
	bw := base64.NewEncoder(base64.StdEncoding, buf)
	if err := imaging.Encode(bw, img, imaging.JPEG, imaging.JPEGQuality(jq)); err != nil {
		return "", "", fmt.Errorf("ImageHandlerBase64: encode new image to dictfile: %w", err)
	}
	_ = bw.Close()

	// generate the css
	css := fmt.Sprintf("width:%dpx;height:%dpx;max-width:100%%;margin:1em auto;page-break-before:auto;object-fit:scale-down;object-position:center", img.Bounds().Dx(), img.Bounds().Dy())

	// build the URL
	return "data:image/jpeg;base64," + buf.String(), css, nil
}

var imgTagRe = regexp.MustCompile(`(<img)(\s+(?:[^>]*\s+)?src\s*=\s*['"]+)([^'"]+)(['"][^>]*>)`)

// transformHTMLImages transforms img tags in the specified HTML, using
// openImage to read the specified paths. If openImage implements io.Closer,
// it will be closed automatically. Img tags which reference have a data URL are
// skipped.
//
// The dictwriter may be used during this process, so callers should not rely on
// any entries opened before calling this.
func transformHTMLImages(ih ImageHandler, dw *kobodict.Writer, html []byte, openImage func(src string) (io.Reader, error)) ([]byte, error) {
	nhtml := html[:]
	for _, m := range imgTagRe.FindAllSubmatch(html, -1) {
		t, a, b, src, c := m[0], m[1], m[2], m[3], m[4]
		if bytes.HasPrefix(src, []byte("data:")) {
			continue
		}
		ir, err := openImage(string(src))
		if err != nil {
			return nil, fmt.Errorf("transform image %#v: open file: %w", string(src), err)
		}
		nsrc, css, err := ih.Transform(string(src), ir, dw)
		if err != nil {
			if c, ok := ir.(io.Closer); ok {
				c.Close()
			}
			return nil, fmt.Errorf("transform image %#v: transform image: %w", string(src), err)
		}
		if c, ok := ir.(io.Closer); ok {
			c.Close()
		}
		var nstyle string
		if len(css) != 0 {
			nstyle = " style=\"" + css + "\""
		}
		if len(nsrc) == 0 {
			nhtml = bytes.Replace(nhtml, t, nil, 1)
		} else {
			nhtml = bytes.Replace(nhtml, t, []byte(string(a)+nstyle+string(b)+nsrc+string(c)), 1)
		}
	}
	return nhtml, nil
}
