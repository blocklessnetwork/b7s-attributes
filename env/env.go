package env

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/blocklessnetwork/b7s-attributes/attributes"
)

// ReadAttributes will read B7S node attributes from the environment.
// All environment variables with the specified `prefix` are taken into account.
// Resulting attribute set will have attribute names where the prefix is stripped.
// Environment variables listed in the `ignore` list are skipped. The `ignore` name should be the full name, including the prefix,
// i.e. if you want to skip `B7S_PRIVATE_IP`, specify the full name, not just the `PRIVATE_IP`.
func ReadAttributes(prefix string, ignore []string) ([]attributes.Attribute, error) {

	environ := os.Environ()

	ignoreMap := make(map[string]struct{})
	for _, ign := range ignore {
		ignoreMap[ign] = struct{}{}
	}

	out := make([]attributes.Attribute, 0)

	for _, env := range environ {
		fields := strings.SplitN(env, "=", 2)

		if len(fields) != 2 {
			return nil, fmt.Errorf("unexpected environment variable format: %s", env)
		}

		// TODO: Check - what if we have `key="value where value contains another = sign in the payload, do we have more fields"`?
		key := fields[0]
		value := fields[1]

		if !strings.HasPrefix(key, prefix) {
			continue
		}

		_, ignored := ignoreMap[key]
		if ignored {
			continue
		}

		name := strings.TrimPrefix(key, prefix)

		attr := attributes.Attribute{Name: name, Value: value}
		out = append(out, attr)
	}

	slices.SortStableFunc(out, func(a, b attributes.Attribute) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return out, nil
}
