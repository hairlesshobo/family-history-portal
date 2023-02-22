package fhm_plg_image_preview

import (
	"io"
	"net/http"
	"path/filepath"

	. "github.com/mickael-kerjean/filestash/server/common"
)

func init() {
	Hooks.Register.ProcessFileContentBeforeSend(renderImages)
}

func renderImages(reader io.ReadCloser, ctx *App, res *http.ResponseWriter, req *http.Request) (io.ReadCloser, error) {
	query := req.URL.Query()
	if query.Get("thumbnail") == "true" {
		return reader, nil
	} else if query.Get("size") == "" {
		return reader, nil
	}

	var (
		out io.ReadCloser = nil
		err error         = nil
	)
	mType := GetMimeType(query.Get("path"))
	switch mType {
	case "image/tiff":
		path := req.URL.Query().Get("path")
		tmb_path := filepath.Join(ctx.Session["path"], path+".jprv")
		prvreader, prverr := ctx.Backend.Cat(tmb_path)
		if prverr == nil {
			out = prvreader
			mType = "image/jpeg"
		} else {
			Log.Warning("fhm_plg_image_preview::renderImages Transcoding " + path + " on the fly")
			out, mType, err = transcodeTiff(reader)
		}
	default:
		return reader, nil
	}
	reader.Close()
	if err == nil {
		(*res).Header().Set("Content-Type", mType)
	}
	if err != nil && err != ErrNotImplemented && err != ErrNotValid {
		Log.Debug("plg_image_transcode::err %s", err.Error())
		return nil, ErrNotValid
	}
	return out, err
}
