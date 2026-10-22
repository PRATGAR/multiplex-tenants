package validate

import (
	"regexp"
	"strings"
)

var (
	slugPattern        = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{2,31}$`)
	searchPattern      = regexp.MustCompile(`^[A-Za-z0-9 _.-]{0,64}$`)
	inviteTokenPattern = regexp.MustCompile(`^[0-9a-f]{32}`)
)

func Slug(s string) bool {
	return slugPattern.MatchString(s)
}

func SearchTerm(s string) bool {
	return searchPattern.MatchString(strings.TrimSpace(s))
}

func Email(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) >= 3 && len(s) <= 254 && strings.Contains(s, "@")
}

func DisplayName(s string) bool {
	s = strings.TrimSpace(s)
	return s != "" && len(s) <= 80
}

func Hostname(s string) bool {
	if len(s) == 0 || len(s) > 253 {
		return false
	}
	if strings.ContainsAny(s, " /?#@" + string(rune(92))) {
		return false
	}
	return strings.Contains(s, ".")
}

func LocalPath(p string) bool {
	if p == "" {
		return false
	}
	return strings.HasPrefix(p, "/") && !strings.Contains(p, "://")
}

func InviteToken(s string) bool {
	return inviteTokenPattern.MatchString(s)
}

func GzipLevel(s string) bool {
	return s != "" && len(s) <= 12 && !strings.ContainsAny(s, ";|&`")
}

func LocalPathOrRoot(p string) string {
	if !LocalPath(p) {
		return "/"
	}
	return p
}
