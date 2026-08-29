package decoder

type Product struct {
	Name      string `json:"name"`
	ModelYear int    `json:"model_year,omitempty"`
}
