package bee

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"log"
	"path"

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

type StaticResourceHandlerOption func(*StaticResourceHandler)

type StaticResourceHandler struct {
	dir                     string
	pathPrefix              string
	extensionContentTypeMap map[string]string
	cache                   *lru.Cache
	maxFileSize             int
}

type cacheItem struct {
	filename    string
	filesize    int
	contentType string
	data        []byte
}

func NewStaticResourceHandler(dir, pathPrefix string, opts ...StaticResourceHandlerOption) *StaticResourceHandler {
	var handler = StaticResourceHandler{
		dir:        dir,
		pathPrefix: pathPrefix,
		extensionContentTypeMap: map[string]string{
			// 这里根据自己的需要不断添加
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"pdf":  "image/pdf",
		},
	}

	for _, opt := range opts {
		opt(&handler)
	}
	return &handler
}

func WithCacheItem(maxFileSizeThreshold int, maxCacheFileCount int) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		cache, err := lru.New(maxCacheFileCount)
		if err != nil {
			log.Println("初始化lru cache 失败")
			return
		}
		handler.maxFileSize = maxFileSizeThreshold
		handler.cache = cache
	}
}

func WithMoreContentTypeExtension(ext map[string]string) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.extensionContentTypeMap = ext
	}
}

func (h *StaticResourceHandler) readFromCache(filename string) (*cacheItem, bool) {
	if h.cache != nil {
		if item, ok := h.cache.Get(filename); ok {
			return item.(*cacheItem), ok
		}
	}
	return nil, false
}

func (h *StaticResourceHandler) writeCacheItem(item *cacheItem, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", item.contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", item.filesize))
	_, _ = w.Write(item.data)
}

func (h *StaticResourceHandler) Handle() HandleFunc {
	return func(c *Context) {
		filename, ok := c.PathParams["file"]
		if !ok {
			c.NotFoundErr("resource not found")
			return
		}
		item, ok := h.readFromCache(filename)
		if ok {
			log.Println("hit cache")
			h.writeCacheItem(item, c.Resp)
			return
		}
		file, err := os.Open(filepath.Join(h.dir, filename))
		if err != nil {

			c.InternalServerErr("服务器开小差了")
			return
		}

		ext := path.Ext(file.Name())
		t, ok := h.extensionContentTypeMap[ext[1:]]
		if !ok {
			c.BadRequestErr("resource not found")
			return
		}
		data, err := io.ReadAll(file)

		if err != nil {
			c.InternalServerErr("服务器开小差了")
			return
		}
		item = &cacheItem{
			filename:    filename,
			filesize:    len(data),
			contentType: t,
			data:        data,
		}
		h.cacheFile(item)
		h.writeCacheItem(item, c.Resp)
	}
}

func (h *StaticResourceHandler) cacheFile(item *cacheItem) {
	if h.cache != nil && item.filesize < h.maxFileSize {
		h.cache.Add(item.filename, item)
	}
}
