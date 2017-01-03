// this file is used to test white box cases
package apns2

import "time"

func GetMaxSize(mi ClientManager) int {
	if m, ok := mi.(*manager); ok {
		return m.maxSize
	}

	return 0
}

func GetMaxAge(mi ClientManager) time.Duration {
	if m, ok := mi.(*manager); ok {
		return m.maxAge
	}

	return time.Duration(0)
}

func GetFactory(mi ClientManager) ClientFactory {
	if m, ok := mi.(*manager); ok {
		return m.factory
	}

	return nil
}
