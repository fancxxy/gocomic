package comic

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"sort"

	"github.com/fancxxy/gocomic/download/parser"
)

func resolvePath(paths []string) (string, error) {
	switch len(paths) {
	case 0:
		return "", nil
	case 1:
		path := paths[0]
		if len(path) == 0 || path[0] != '~' {
			return path, nil
		}

		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		return filepath.Join(usr.HomeDir, path[1:]), nil
	default:
		return "", fmt.Errorf("too many parameters")
	}
}

func createDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func clipCover(in io.Reader, out io.Writer, scale float64) error {
	img, format, err := image.Decode(in)
	if err != nil {
		return err
	}
	bounds := img.Bounds()
	var subBounds image.Rectangle
	width, height := bounds.Dx(), bounds.Dy()
	if s := float64(width) / float64(height); s == scale {
		subBounds = bounds
	} else if s > scale {
		redundant := int((float64(width) - float64(height)*scale) / 2)
		subBounds = image.Rect(bounds.Min.X+redundant, bounds.Min.Y, bounds.Max.X-redundant, bounds.Max.Y)
	} else {
		redundant := int((float64(height) - float64(width)/scale) / 2)
		subBounds = image.Rect(bounds.Min.X, bounds.Min.Y+redundant, bounds.Max.X, bounds.Max.Y-redundant)
	}

	switch format {
	case "jpeg":
		converted := img.(*image.YCbCr)
		subImg := converted.SubImage(subBounds).(*image.YCbCr)
		return jpeg.Encode(out, subImg, &jpeg.Options{Quality: 100})
	default:
		return fmt.Errorf("not support image format")
	}
}

type titles struct {
	data   []string
	parser parser.Parser
}

func (t titles) Less(i, j int) bool {
	return t.parser.Less(t.data[i], t.data[j])
}

func (t titles) Len() int {
	return len(t.data)
}

func (t titles) Swap(i, j int) {
	t.data[i], t.data[j] = t.data[j], t.data[i]
}

func sortTitle(data []string, parser parser.Parser) {
	ts := titles{parser: parser, data: data}
	sort.Sort(ts)
}
