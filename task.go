package zabbix

import "context"

// TaskService wraps the "task" API namespace.
type TaskService struct{ client *Client }

// Task returns the task service.
func (c *Client) Task() *TaskService { return &TaskService{c} }

// Task object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/task/object
type Task struct {
	TaskID  string `json:"taskid,omitempty"`
	Type    string `json:"type,omitempty"`
	Status  string `json:"status,omitempty"`
	Clock   string `json:"clock,omitempty"`
	TTL     string `json:"ttl,omitempty"`
	ProxyID string `json:"proxyid,omitempty"`
	Result  any    `json:"result,omitempty"`
}

// TaskCreate describes a task to create. Type selects the task kind and Request
// carries the type-specific request payload.
type TaskCreate struct {
	Type    int `json:"type"`
	Request any `json:"request"`
}

// TaskGetParams are the parameters for task.get. Unlike most get methods, task.get
// does NOT accept the common query options (limit, sort, filter, search), so it
// does not embed GetParams — only the parameters the server actually supports.
type TaskGetParams struct {
	Output           any      `json:"output,omitempty"`
	TaskIDs          []string `json:"taskids,omitempty"`
	Preservekeys     bool     `json:"preservekeys,omitempty"`
	SelectLastValues any      `json:"selectLastValues,omitempty"`
}

// Create creates tasks and returns the new task ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/task/create
func (s *TaskService) Create(ctx context.Context, tasks ...TaskCreate) ([]string, error) {
	var res struct {
		IDs []string `json:"taskids"`
	}
	err := s.client.Call(ctx, "task.create", tasks, &res)
	return res.IDs, err
}

// Get retrieves tasks matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/task/get
func (s *TaskService) Get(ctx context.Context, params TaskGetParams) ([]Task, error) {
	var out []Task
	err := s.client.Call(ctx, "task.get", params, &out)
	return out, err
}
