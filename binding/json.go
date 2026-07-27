// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

var jsonBufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func (jsonBinding) Bind(req *http.Request, obj any) error {
	if req == nil || req.Body == nil {
		return nil
	}

	buf := jsonBufferPool.Get().(*bytes.Buffer)
	buf.Reset()

	_, err := buf.ReadFrom(req.Body)
	req.Body.Close()
	if err != nil {
		jsonBufferPool.Put(buf)
		return err
	}

	bodyBytes := buf.Bytes()
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	err = decodeJSON(req.Body, obj)
	if err != nil {
		failBytes := make([]byte, len(bodyBytes))
		copy(failBytes, bodyBytes)
		req.Body = io.NopCloser(bytes.NewReader(failBytes))
		jsonBufferPool.Put(buf)
		return err
	}

	req.Body = http.NoBody
	jsonBufferPool.Put(buf)
	return nil
}

func (jsonBinding) BindBody(body []byte, obj any) error {
	return decodeJSON(bytes.NewReader(body), obj)
}

func decodeJSON(r io.Reader, obj any) error {
	decoder := json.NewDecoder(r)
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if EnableDecoderDisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
