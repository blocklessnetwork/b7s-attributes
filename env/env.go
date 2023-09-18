package env

import (
	"fmt"
	"os"
	"strings"
)

// ReadAttributes will read B7S node attributes from the environment.
// All environment variables with the specified `prefix` are taken into account.
// Resulting attribute set will have attribute names where the prefix is stripped.
// Environment variables listed in the `ignore` list are skipped. The `ignore` name should be the full name, including the prefix,
// i.e. if you want to skip `B7S_PRIVATE_IP`, specify the full name, not just the `PRIVATE_IP`.
// `limit` specifies how many attributes to include. In case there are more variables, the limited set is returned, alongside an error.
func ReadAttributes(prefix string, limit uint, ignore []string) (map[string]string, error) {

	environ := os.Environ()

	ignoreMap := make(map[string]struct{})
	for _, ign := range ignore {
		ignoreMap[ign] = struct{}{}
	}

	out := make(map[string]string)
	count := uint(0)

	for _, env := range environ {
		fields := strings.Split(env, "=")

		if len(fields) < 2 {
			return nil, fmt.Errorf("unexpected environment variable format: %s", env)
		}

		// TODO: Check - what if we have `key="value where value contains another = sign in the payload, do we have more fields"`?
		key := fields[0]
		value := fields[1]

		if !strings.HasPrefix(key, prefix) {
			continue
		}

		_, ok := ignoreMap[key]
		if ok {
			continue
		}

		name := strings.TrimPrefix(key, prefix)

		if count < limit {
			out[name] = value
			count++
			continue
		}

		return out, ErrTooManyAttributes
	}

	return out, nil
}
