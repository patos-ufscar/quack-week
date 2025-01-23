package schemas

type Id struct {
	Id string `json:"id" binding:"required"`
}

type Url struct {
	Url string `json:"url" binding:"required"`
}

type Name struct {
	Name string `json:"name" binding:"required"`
}

type UploadPicture struct {
	Content string `json:"content" binding:"required" example:"base64 encoded string"`
}
