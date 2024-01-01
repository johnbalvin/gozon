package gozon

type Data struct {
	Asin            Asin
	URL             string
	Title           string
	MerchandID      string
	Rating          float32
	Available       bool
	Price           Price
	ShippingAddress string
	Variations      Variations
	MainPanelDesc   MainPanelDesciptions
	Images          []Imgs
}
type MainPanelDesciptions struct {
	Details    []LabelValue
	AboutThis  []string
	IconsPanel []IconPanel
}
type IconPanel struct {
	URL   string
	Label string
	Value string
}
type Variations struct {
	Labels       []LabelValues
	Combinations []Combination
}
type dataRaw struct {
	Asin                   Asin
	Title                  string
	URL                    string
	Alterts                string
	MerchandID             string
	Rating                 float32
	Price                  Price
	Available              bool
	Dimensions             []string
	ShippingAddress        string
	VariationDisplayLabels map[string][]string
	Variations             map[string][]string
	MainPanelDesc          MainPanelDesciptions
	Images                 []ImgsPage
}
type Asin struct {
	Me     string
	Parent string
}
type LabelValue struct {
	Label string
	Value string
}
type LabelValues struct {
	Label  string
	Values []string
}
type Combination struct {
	Asin   string
	Values []string
}
type Label struct {
	Size_name  string `json:"size_name"`
	Color_name string `json:"color_name"`
}
type Price struct {
	Low               float32
	High              float32
	Base              float32
	SavingsPercentage int
	CurrencySymbol    string
}
type ImgsPage struct {
	Variant string `json:"variant"`
	Large   string `json:"large"`
	Thumb   string `json:"thumb"`
	HiRes   string `json:"hiRes"`
}

type Imgs struct {
	Variant string
	Large   URLImg
	Thumb   URLImg
	HiRes   URLImg
}
type URLImg struct {
	ContentType string
	Extension   string
	URL         string
	Content     []byte `json:"-"`
}
