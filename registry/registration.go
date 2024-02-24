/**
 * @Author: Keven5
 * @Description:
 * @File:  registration
 * @Version: 1.0.0
 * @Date: 2024/2/24 11:42
 */

package registry

type ServiceName string

type Registration struct {
	ServiceName       ServiceName
	ServiceUrl        string
	RequiredServices  []ServiceName
	ServiceUpdatedURL string
}

const (
	LogService     = ServiceName("Log_Service")
	GradingService = ServiceName("Grading_Service")
)

type patchEntry struct {
	Name ServiceName
	URL  string
}

type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
