package application

type EditorData struct {
	Time    int64         `json:"time"`
	Blocks  []EditorBlock `json:"blocks"`
	Version string        `json:"version"`
}

type EditorBlock struct {
	ID   string    `json:"id"`
	Type string    `json:"type"`
	Data BlockData `json:"data"`
}

type BlockData struct {
	// header / paragraph
	Text string `json:"text,omitempty"`

	// header
	Level int `json:"level,omitempty"`

	// list
	Style string   `json:"style,omitempty"`
	Items []string `json:"items,omitempty"`

	// image
	File    *ImageFile `json:"file,omitempty"`
	Caption string     `json:"caption,omitempty"`

	// linkTool
	Link string    `json:"link,omitempty"`
	Meta *LinkMeta `json:"meta,omitempty"`
}

type ImageFile struct {
	URL string `json:"url"`
}

type LinkMeta struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Image       LinkImage `json:"image,omitempty"`
}

type LinkImage struct {
	URL string `json:"url,omitempty"`
}
