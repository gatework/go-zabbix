package zabbix

import "context"

// ImageService wraps the "image" API namespace.
type ImageService struct{ client *Client }

// Image returns the image service.
func (c *Client) Image() *ImageService { return &ImageService{c} }

// Image object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/image/object
type Image struct {
	ImageID   string `json:"imageid,omitempty"`
	ImageType string `json:"imagetype,omitempty"`
	Name      string `json:"name,omitempty"`
	Image     string `json:"image,omitempty"` // base64-encoded image data
}

// ImageGetParams are the parameters for image.get.
// ImageGetParams queries the list form of image.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type ImageGetParams struct {
	Output                 any            `json:"output,omitempty"`
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	SortField              []string       `json:"sortfield,omitzero"`
	SortOrder              []string       `json:"sortorder,omitzero"`
	Limit                  int            `json:"limit,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	ImageIDs               []string       `json:"imageids,omitzero"`
	SysmapIDs              []string       `json:"sysmapids,omitzero"`
	SelectImage            bool           `json:"select_image,omitempty"`
}

// Get retrieves images matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/image/get
func (s *ImageService) Get(ctx context.Context, params ImageGetParams) ([]Image, error) {
	return getList[Image](ctx, s.client, "image.get", params)
}

// ImageCountParams contains the filters for image.get counts.
type ImageCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	ImageIDs               []string       `json:"imageids,omitzero"`
	SysmapIDs              []string       `json:"sysmapids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *ImageService) Count(ctx context.Context, params ImageCountParams) (uint64, error) {
	return getCount(ctx, s.client, "image.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *ImageService) GetMap(ctx context.Context, params ImageGetParams) (map[string]Image, error) {
	return getMap[Image](ctx, s.client, "image.get", params)
}

// Create creates images and returns the new image ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/image/create
func (s *ImageService) Create(ctx context.Context, images ...ImageCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "image.create", "imageids", images)
}

// Update updates images and returns the affected image ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/image/update
func (s *ImageService) Update(ctx context.Context, images ...ImageUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "image.update", "imageids", images)
}

// Delete deletes images by id and returns the deleted image ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/image/delete
func (s *ImageService) Delete(ctx context.Context, imageIDs ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "image.delete", "imageids", imageIDs)
}
