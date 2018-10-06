package libhac

type idResponse struct {
	IDPairs []idPair `json:"id_pair"`
}

type idPair struct {
	ID int `json:"id"`
}

type Title struct {
	ID           int           `json:"id"`
	Name         string        `json:"formal_name"`
	BannerURL    string        `json:"hero_banner_url"`
	ReleaseDate  string        `json:"release_date_on_eshop"`
	IsNew        bool          `json:"is_new"`
	Description  string        `json:"description"`
	Genre        string        `json:"genre"`
	Size         int           `json:"total_rom_size"`
	Screenshots  []Screenshot  `json:"screenshots"`
	Movies       []Movie       `json:"movies"`
	Publisher    Publisher     `json:"publisher"`
	Applications []Application `json:"applications"`
}

type Screenshot struct {
	Images []Image `json:"images"`
}

type Image struct {
	URL string `json:"url"`
}

type Movie struct {
	URL       string `json:"movie_url"`
	Thumbnail string `json:"thumbnail_url"`
}

type Publisher struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Application struct {
	ID    string `json:"id"`
	Image string `json:"image_url"`
}
