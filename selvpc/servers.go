package selvpc

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/selectel/go-selvpcclient/selvpcclient/resell/v2/servers"
)

// serversMapsFromStructs converts the provided license.Servers to
// the slice of maps correspondingly to the resource's schema.
func serversMapsFromStructs(serverStructs []servers.Server) []map[string]interface{} {
	associatedServers := make([]map[string]interface{}, len(serverStructs))

	if len(serverStructs) != 0 {
		for i, server := range serverStructs {
			associatedServers[i] = map[string]interface{}{
				"id":     server.ID,
				"name":   server.Name,
				"status": server.Status,
			}
		}
	}

	return associatedServers
}

// hashServers is a hash function to use with the "servers" set.
func hashServers(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["id"].(string)))
	return hashcode.String(buf.String())
}
