package api

import "golang-distributed-parallel-image-processing/api/login"

func isTokenActive(token string) bool {
	if _, ok := login.ActiveTokens[token]; ok {
		return true
	}
	return false
}
