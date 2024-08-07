package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/NikolosHGW/metric/internal/server/logger"
	"go.uber.org/zap"
)

const (
	applicationJSON = "application/json"
	textHTML        = "html/text"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func WithGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := strings.Join(r.Header.Values("Accept-Encoding"), ", ")
		contentType := strings.Join(r.Header.Values("Content-Type"), ", ")
		accept := strings.Join(r.Header.Values("Accept"), ", ")

		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		supportApplicationJSON := strings.Contains(contentType, applicationJSON)
		supportTextHTML := strings.Contains(accept, textHTML)

		if supportsGzip && (supportApplicationJSON || supportTextHTML) {
			cw := newCompressWriter(w)
			ow = cw

			defer func() {
				err := cw.Close()
				if err != nil {
					logger.Log.Info("err close compressWriter", zap.Error(err))
				}
			}()
		}

		contentEncoding := strings.Join(r.Header.Values("Content-Encoding"), ", ")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer func() {
				err := cr.Close()
				if err != nil {
					logger.Log.Info("err close compressReader", zap.Error(err))
				}
			}()
		}

		next.ServeHTTP(ow, r)
	})
}
