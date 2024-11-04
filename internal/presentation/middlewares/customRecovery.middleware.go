package middlewares

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

func CustomRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// 在這裡添加您的自定義日誌記錄邏輯
				log.Error().Msgf("Custom recovery log: Panic recovered: %v", err)
				// 返回 500 狀態碼和錯誤消息
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
