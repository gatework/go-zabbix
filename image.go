package zabbix

import "context"

// ImageService wraps the "image" API namespace.
type ImageService struct{ client *Client }

// Image returns the image service.
func (c *Client) Image() *ImageService { return &ImageService{c} }

// Image object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/image/object
type Image struct {
	ImageID   string `json:"imageid,omitempty"`
	ImageType string `json:"imagetype,omitempty"`
	Name      string `json:"name,omitempty"`
	Image     string `json:"image,omitempty"` // base64-encoded image data
}

// ImageGetParams are the parameters for image.get.
type ImageGetParams struct {
	GetParams
	ImageIDs    []string `json:"imageids,omitempty"`
	SysmapIDs   []string `json:"sysmapids,omitempty"`
	SelectImage bool     `json:"select_image,omitempty"`
}

// Get retrieves images matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/image/get
func (s *ImageService) Get(ctx context.Context, params ImageGetParams) ([]Image, error) {
	var out []Image
	err := s.client.Call(ctx, "image.get", params, &out)
	return out, err
}

// Create creates images and returns the new image ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/image/create
func (s *ImageService) Create(ctx context.Context, images ...Image) ([]string, error) {
	var res struct {
		IDs []string `json:"imageids"`
	}
	err := s.client.Call(ctx, "image.create", images, &res)
	return res.IDs, err
}

// Update updates images and returns the affected image ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/image/update
func (s *ImageService) Update(ctx context.Context, images ...Image) ([]string, error) {
	var res struct {
		IDs []string `json:"imageids"`
	}
	err := s.client.Call(ctx, "image.update", images, &res)
	return res.IDs, err
}

// Delete deletes images by id and returns the deleted image ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/image/delete
func (s *ImageService) Delete(ctx context.Context, imageIDs ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"imageids"`
	}
	err := s.client.Call(ctx, "image.delete", imageIDs, &res)
	return res.IDs, err
}
