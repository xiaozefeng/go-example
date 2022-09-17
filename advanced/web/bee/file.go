package bee

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type FileUploader struct {
	FileField    string
	DestPathFunc func(fh *multipart.FileHeader) string
}

func (f *FileUploader) Handle() HandleFunc {
	return func(ctx *Context) {
		src, srcHeader, err := ctx.Request.FormFile(f.FileField)
		if err != nil {
			ctx.BadRequestErr("上传失败, 未找到数据")
			return
		}
		defer src.Close()

		dest, err := os.OpenFile(f.DestPathFunc(srcHeader), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			ctx.InternalServerErr("上传失败")
			return
		}
		defer dest.Close()

		_, err = io.CopyBuffer(dest, src, nil)
		if err != nil {
			ctx.InternalServerErr("上传失败")
			return
		}

		_ = ctx.WriteString("上传成功")
	}
}

type FileDownloader struct {
	Dir string
}

func (f *FileDownloader) Handle() HandleFunc {
	return func(c *Context) {
		reqFilename, _ := c.QueryValue("file")
		path := filepath.Join(f.Dir, filepath.Clean(reqFilename))
		fn := filepath.Base(path)
		header := c.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")
		http.ServeFile(c.Resp, c.Request, path)
	}
}
