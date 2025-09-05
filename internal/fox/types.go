package fox

type RealDataRequestBody struct {
	SN        string   `json:"sn"`
	Variables []string `json:"variables"`
}

type RealDataResponse struct {
	Errno  int              `json:"errno"`
	Msg    string           `json:"msg"`
	Result []RealDataResult `json:"result"`
}

type RealDataResult struct {
	Datas    []Datapoint `json:"datas"`
	DeviceSN string      `json:"deviceSN"`
	Time     string      `json:"time"`
}

type Datapoint struct {
	Name     string  `json:"name"`
	Unit     string  `json:"unit"`
	Value    float64 `json:"value"`
	Variable string  `json:"variable"`
}
