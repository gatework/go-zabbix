package zabbix

import "context"

// TaskService wraps the "task" API namespace.
type TaskService struct{ client *Client }

// Task returns the task service.
func (c *Client) Task() *TaskService { return &TaskService{c} }

// Task object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/task/object
type Task struct {
	TaskID  string `json:"taskid,omitempty"`
	Type    string `json:"type,omitempty"`
	Status  string `json:"status,omitempty"`
	Clock   string `json:"clock,omitempty"`
	TTL     string `json:"ttl,omitempty"`
	ProxyID string `json:"proxyid,omitempty"`
	Result  any    `json:"result,omitempty"`
}

// TaskGetParams are the parameters for task.get. Unlike most get methods, task.get
// does NOT accept the common query options (limit, sort, filter, search), so it
// exposes only its supported parameters.
// TaskGetParams queries the list form of task.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type TaskGetParams struct {
	Output  any      `json:"output,omitempty"`
	TaskIDs []string `json:"taskids,omitzero"`
}

// Create creates tasks and returns the new task ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/task/create
func (s *TaskService) Create(ctx context.Context, tasks ...TaskCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "task.create", "taskids", tasks)
}

// Get retrieves tasks matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/task/get
func (s *TaskService) Get(ctx context.Context, params TaskGetParams) ([]Task, error) {
	return getList[Task](ctx, s.client, "task.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *TaskService) GetMap(ctx context.Context, params TaskGetParams) (map[string]Task, error) {
	return getMap[Task](ctx, s.client, "task.get", params)
}
