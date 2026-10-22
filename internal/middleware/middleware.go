package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"

	"tenantly/internal/httpx"
	"tenantly/internal/model"
)

type sessionKey struct{}

const SessionPrefix = "sess:"

type Guard struct {
	rdb *redis.Client
}

func NewGuard(rdb *redis.Client) *Guard {
	return &Guard{rdb: rdb}
}

func SessionFrom(ctx context.Context) (model.Session, bool) {
	s, ok := ctx.Value(sessionKey{}).(model.Session)
	return s, ok
}

func (g *Guard) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			httpx.Error(w, http.StatusUnauthorized, "missing credentials")
			return
		}
		raw, err := g.rdb.Get(r.Context(), SessionPrefix+token).Result()
		if err != nil {
			httpx.Error(w, http.StatusUnauthorized, "session expired")
			return
		}
		var sess model.Session
		if err := json.Unmarshal([]byte(raw), &sess); err != nil {
			httpx.Error(w, http.StatusUnauthorized, "session invalid")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), sessionKey{}, sess)))
	})
}

func TenantScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, ok := SessionFrom(r.Context())
		if !ok || mux.Vars(r)["tid"] != strconv.FormatInt(sess.TenantID, 10) {
			httpx.Error(w, http.StatusNotFound, "tenant not found")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, role := range roles {
		allowed[role] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, ok := SessionFrom(r.Context())
			if !ok || !allowed[sess.Role] {
				httpx.Error(w, http.StatusForbidden, "insufficient role")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func PathCheck(name string, valid func(string) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !valid(mux.Vars(r)[name]) {
				httpx.Error(w, http.StatusNotFound, "not found")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func QueryCheck(name string, valid func(string) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if v := r.URL.Query().Get(name); v != "" && !valid(v) {
				httpx.Error(w, http.StatusBadRequest, "unsupported "+name)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func QueryRewrite(name string, fn func(string) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			q.Set(name, fn(q.Get(name)))
			r.URL.RawQuery = q.Encode()
			next.ServeHTTP(w, r)
		})
	}
}

func CSVAttachment(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=export-%s.csv", mux.Vars(r)["eid"]))
		next.ServeHTTP(w, r)
	})
}
