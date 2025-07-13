package service

import (
	"context"
	"crypto/sha256"
	"encoding/binary"

	"feature_flag_config/db/cache"

	"github.com/sirupsen/logrus"
)

func (s *Service) checkFeature(
	ctx context.Context,
	featureName string,
	itemID *string,
) (
	enabled bool,
	reason string,
) {
	feature := cache.Feature{}
	if err := feature.GetByName(ctx, featureName, s.redisDB); err != nil {
		logrus.Errorf("GetByName error: %v", err)
	}

	if !feature.IsActive {
		return false, "feature is unavailable"
	}

	// there is no need to check feature if we don't get item id
	if itemID == nil {
		return true, ""
	}

	// if the item is not in available buckets
	if itemGroup := getItemGroup(*itemID, 100); itemGroup > feature.BucketsOpened {
		// check if the item is in the feature's whitelist
		if len(feature.Whitelist) > 0 {
			for i := 0; i < len(feature.Whitelist); i++ {
				if feature.Whitelist[i] == *itemID {
					enabled = true
				}
			}
		}

		// feature is unavailable if the item is not in availbale buckets and not in whitelist
		if !enabled {
			return false, "item is not in available bucket"
		}
	}

	// check the feature's blacklist
	if len(feature.Blacklist) > 0 {
		for i := 0; i < len(feature.Blacklist); i++ {
			if feature.Blacklist[i] == *itemID {
				return false, "item is in blacklist"
			}
		}
	}

	return true, ""
}

// get item's group using hash of the item and num of available groups
func getItemGroup(itemID string, numGroups int) int {
	hash := sha256.Sum256([]byte(itemID))
	hashInt := binary.BigEndian.Uint32(hash[len(hash)-4:])
	return int(hashInt % uint32(numGroups))
}
