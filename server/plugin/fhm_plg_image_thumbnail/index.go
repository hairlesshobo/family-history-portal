package fhm_plg_image_thumbnail

import (
	_ "embed"
	"io"
	"net/http"
	"path/filepath"

	. "github.com/mickael-kerjean/filestash/server/common"
)

// go:embed dist/placeholder.png
var placeholder []byte

func init() {
	Hooks.Register.Thumbnailer("image/tiff", thumbnailBuilder{thumbnailTiff})
}

func thumbnailTiff(reader io.ReadCloser, ctx *App, res *http.ResponseWriter, req *http.Request) (io.ReadCloser, error) {
	h := (*res).Header()

	path := req.URL.Query().Get("path")
	tmb_path := filepath.Join(ctx.Session["path"], path+".jtmb")
	reader, err := ctx.Backend.Cat(tmb_path)
	if err != nil {
		h.Set("Content-Type", "image/png")
		h.Set("Cache-Control", "max-age=1")
		return NewReadCloserFromBytes(placeholder), nil
	}

	// Log.Debug("results: " + fmt.Sprintf("%v", len(fileinfo)))
	// Log.Debug("result[0].size:" + fmt.Sprintf("%v", fileinfo[0].Size()))

	return reader, nil
}

type thumbnailBuilder struct {
	fn func(reader io.ReadCloser, ctx *App, res *http.ResponseWriter, req *http.Request) (io.ReadCloser, error)
}

func (this thumbnailBuilder) Generate(reader io.ReadCloser, ctx *App, res *http.ResponseWriter, req *http.Request) (io.ReadCloser, error) {
	return this.fn(reader, ctx, res, req)
}
