/*
 * Project: banker
 * File: recorder.go (2/20/24, 10:57 AM)
 *
 * Copyright (C) Megakit Systems 2017-2024, Inc - All Rights Reserved
 * @link https://www.megakit.pro
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Anton (antonstremovskyy) Stremovskyy <stremovskyy@gmail.com>
 */

package redis_recorder

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/stremovskyy/gofondy/recorder"
)

const (
	RequestPrefix  = "request"
	ResponsePrefix = "response"
	ErrorPrefix    = "error"
	MetricsPrefix  = "metrics"
	TagsPrefix     = "tag"
)

type redisRecorder struct {
	client     *redis.Client
	options    *Options
	compressor *compressor
	logger     *log.Logger
}

// NewRedisRecorder creates a new instance of redisRecorder.
func NewRedisRecorder(options *Options) recorder.Client {
	client := redis.NewClient(
		&redis.Options{
			Addr:       options.Addr,
			Password:   options.Password,
			DB:         options.DB,
			ClientName: "RedisRecorder",
		},
	)

	statusCmd := client.Ping(context.Background())
	if statusCmd.Err() != nil {
		panic("failed to connect to redis server: " + statusCmd.Err().Error())
	}

	return &redisRecorder{
		client:     client,
		options:    options,
		compressor: newCompressor(),
		logger:     log.New(log.Writer(), "[Redis Recorder]: ", log.LstdFlags),
	}
}

func (r *redisRecorder) RecordRequest(ctx context.Context, orderID *string, requestID string, request []byte, tags map[string]string) error {
	compressedRequest, err := r.compressor.compressData(request, r.options.CompressionLvl)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s:%s", r.options.Prefix, RequestPrefix, requestID)

	if orderID != nil {
		key = fmt.Sprintf("%s:%s:%s:%s", r.options.Prefix, RequestPrefix, *orderID, requestID)
	}

	tags["request_id"] = requestID

	if r.options.Debug {
		r.logger.Printf("request key: %s", key)
	}

	err = r.client.Set(ctx, key, compressedRequest, r.options.DefaultTTL).Err()
	if err != nil {
		return err
	}

	return r.updateTagIndex(ctx, tags, key)
}

func (r *redisRecorder) RecordResponse(ctx context.Context, orderID *string, requestID string, response []byte, tags map[string]string) error {
	compressedResponse, err := r.compressor.compressData(response, r.options.CompressionLvl)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s:%s", r.options.Prefix, ResponsePrefix, requestID)

	if orderID != nil {
		key = fmt.Sprintf("%s:%s:%s:%s", r.options.Prefix, ResponsePrefix, *orderID, requestID)
	}

	tags["request_id"] = requestID

	if r.options.Debug {
		r.logger.Printf("response key: %s", key)
	}

	err = r.client.Set(ctx, key, compressedResponse, r.options.DefaultTTL).Err()
	if err != nil {
		return err
	}

	return r.updateTagIndex(ctx, tags, key)
}

func (r *redisRecorder) RecordError(ctx context.Context, id *string, requestID string, err error, tags map[string]string) error {
	compressedResponse, err := r.compressor.compressData([]byte(err.Error()), r.options.CompressionLvl)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s:%s", r.options.Prefix, ErrorPrefix, requestID)

	if id != nil {
		key = fmt.Sprintf("%s:%s:%s:%s", r.options.Prefix, ErrorPrefix, *id, requestID)
	}

	tags["request_id"] = requestID

	if r.options.Debug {
		r.logger.Printf("error key: %s", key)
	}

	err = r.client.Set(ctx, key, compressedResponse, r.options.DefaultTTL).Err()
	if err != nil {
		return err
	}

	return r.updateTagIndex(ctx, tags, key)
}

func (r *redisRecorder) RecordMetrics(ctx context.Context, orderID *string, requestID string, metrics map[string]string, tags map[string]string) error {
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("cannot marshal metrics: %w", err)
	}

	key := fmt.Sprintf("%s:%s:%s", r.options.Prefix, MetricsPrefix, requestID)

	if orderID != nil {
		key = fmt.Sprintf("%s:%s:%s:%s", r.options.Prefix, MetricsPrefix, *orderID, requestID)
	}

	if r.options.Debug {
		r.logger.Printf("metrics key: %s", key)
	}

	compressedMetrics, err := r.compressor.compressData(jsonData, r.options.CompressionLvl)
	if err != nil {
		return fmt.Errorf("cannot compress metrics: %w", err)
	}

	err = r.client.Set(ctx, key, compressedMetrics, r.options.DefaultTTL).Err()
	if err != nil {
		return fmt.Errorf("cannot set metrics: %w", err)
	}

	return r.updateTagIndex(ctx, tags, key)
}

func (r *redisRecorder) GetRequest(ctx context.Context, requestID string) ([]byte, error) {
	key := fmt.Sprintf("%s:%s:%s", r.options.Prefix, RequestPrefix, requestID)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return r.compressor.decompressData([]byte(data))
}

func (r *redisRecorder) GetResponse(ctx context.Context, requestID string) ([]byte, error) {
	key := fmt.Sprintf("%s:%s:%s", r.options.Prefix, ResponsePrefix, requestID)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return r.compressor.decompressData([]byte(data))
}

func (r *redisRecorder) FindByTag(ctx context.Context, tag string) ([]string, error) {
	tagKey := fmt.Sprintf("%s:%s:%s", r.options.Prefix, TagsPrefix, tag)

	return r.client.SMembers(ctx, tagKey).Result()
}
