/*
 * Project: banker
 * File: compressor.go (2/20/24, 11:13 AM)
 *
 * Copyright (C) Megakit Systems 2017-2024, Inc - All Rights Reserved
 * @link https://www.megakit.pro
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Anton (antonstremovskyy) Stremovskyy <stremovskyy@gmail.com>
 */

package redis_recorder

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"
)

type compressor struct {
	bufferPool sync.Pool
}

func newCompressor() *compressor {
	return &compressor{
		bufferPool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

func (c *compressor) compressData(data []byte, lvl int) ([]byte, error) {
	buf := c.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer c.bufferPool.Put(buf)

	gz, err := gzip.NewWriterLevel(buf, lvl)
	if err != nil {
		return nil, err
	}
	if _, err = gz.Write(data); err != nil {
		return nil, err
	}
	if err = gz.Close(); err != nil {
		return nil, err
	}

	compressedData := make([]byte, buf.Len())
	copy(compressedData, buf.Bytes())

	return compressedData, nil
}

func (c *compressor) decompressData(compressedData []byte) ([]byte, error) {
	buf := c.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer c.bufferPool.Put(buf)

	_, err := buf.Write(compressedData)
	if err != nil {
		return nil, err
	}

	gz, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	decompressedData, err := io.ReadAll(gz)
	if err != nil {
		return nil, err
	}

	return decompressedData, nil
}
