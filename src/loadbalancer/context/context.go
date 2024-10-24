package context

import "net/http"

const (
  Attempts int = iota
  Retry
)

func GetAttemptsFromContext(r *http.Request) int {
  if attempts, ok := r.Context().Value(Attempts).(int); ok {
    return attempts
  }

  return 1
}

func GetRetryFromContext(r *http.Request) int {
  if retries, ok := r.Context().Value(Retry).(int); ok {
    return retries
  }

  return 0
}
