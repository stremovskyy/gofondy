/*
 * Project: banker
 * File: helpers.go (2/20/24, 11:17 AM)
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
	"fmt"
)

func (r *redisRecorder) updateTagIndex(ctx context.Context, tags map[string]string, itemKey string) error {
	for key, value := range tags {
		tagKey := fmt.Sprintf("%s:%s:%s:%s", r.options.Prefix, TagsPrefix, key, value)
		tagValue := fmt.Sprintf("%s", itemKey)

		_, err := r.client.SAdd(ctx, tagKey, tagValue).Result()
		if err != nil {
			return fmt.Errorf("failed to add tag to index: %s", err.Error())
		}

		_, err = r.client.Expire(ctx, tagKey, r.options.DefaultTTL).Result()
		if err != nil {
			r.logger.Printf("[ERROR] failed to set expiration for tag: %s", err.Error())
		}
	}

	return nil
}
