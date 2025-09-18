package templates

import (
	"encoding/json"
	"fmt"
	"io"

	cpu_k3s "github.com/JarcauCristian/ztp-mcp/internal/server/templates/cpu_k3s_deployment"
	cpu_k8s "github.com/JarcauCristian/ztp-mcp/internal/server/templates/cpu_k8s_deployment"
)

type ZTPTemplate interface {
	Execute() (string, error)
}

func RetrieveModel(templateId string, body io.ReadCloser) (ZTPTemplate, error) {
	switch templateId {
	case "cpu_k8s_deployment":
		var ck8d cpu_k8s.CpuK8sDeployment

		if err := json.NewDecoder(body).Decode(&ck8d); err != nil {
			return nil, fmt.Errorf("Failed to parse body: %v", err)
		}

		return &ck8d, nil
	case "cpu_k3s_deployment":
		var ck3d cpu_k3s.CpuK3sDeployment

		if err := json.NewDecoder(body).Decode(&ck3d); err != nil {
			return nil, fmt.Errorf("Failed to parse body: %v", err)
		}

		return &ck3d, nil
	default:
		return nil, fmt.Errorf("Template id %v does not exist.", templateId)
	}
}
