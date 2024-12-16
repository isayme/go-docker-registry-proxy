package src

import (
	"fmt"
	"strings"
)

type WwwAuthenticate struct {
	Realm   string
	Service string
	Scope   string
	Extras  map[string]string
}

func ParseWwwAuthenticate(v string) (*WwwAuthenticate, bool) {
	result := &WwwAuthenticate{
		Extras: make(map[string]string),
	}

	headerParts := strings.Split(v, " ")
	if len(headerParts) != 2 || headerParts[0] != BEARER {
		return nil, false
	}

	params := headerParts[1]
	for _, param := range strings.Split(params, ",") {
		keyValue := strings.SplitN(strings.TrimSpace(param), "=", 2)
		if len(keyValue) == 2 {
			key := keyValue[0]
			value := strings.Trim(keyValue[1], `"`)

			if key == SCOPE {
				result.Scope = value
			} else if key == REALM {
				result.Realm = value
			} else if key == SERVICE {
				result.Service = value
			} else {
				result.Extras[key] = value
			}
		}
	}

	return result, true
}

func (info *WwwAuthenticate) String() string {
	keyValus := []string{}

	if info.Realm != "" {
		keyValus = append(keyValus, fmt.Sprintf(`%s="%s"`, REALM, info.Realm))
	} else if info.Service != "" {
		keyValus = append(keyValus, fmt.Sprintf(`%s="%s"`, SERVICE, info.Service))
	} else if info.Scope != "" {
		keyValus = append(keyValus, fmt.Sprintf(`%s="%s"`, SCOPE, info.Scope))
	} else {
		for key, value := range info.Extras {
			keyValus = append(keyValus, fmt.Sprintf(`%s="%s"`, key, value))
		}
	}

	return fmt.Sprintf("%s %s", BEARER, strings.Join(keyValus, ","))
}
