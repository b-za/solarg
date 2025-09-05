package tuya

type TokenResponse struct {
	Result  TokenResult `json:"result"`
	Success bool        `json:"success"`
	T       int64       `json:"t"`
	TID     string      `json:"tid"`
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
}

type TokenResult struct {
	AccessToken  string `json:"access_token"`
	ExpireTime   int    `json:"expire_time"`
	RefreshToken string `json:"refresh_token"`
	UID          string `json:"uid"`
}

type DeviceStatus struct {
	Code  string      `json:"code"`
	Value interface{} `json:"value"`
}

type DeviceStatusResponse struct {
	Result  []DeviceStatus `json:"result"`
	Success bool           `json:"success"`
	T       int64          `json:"t"`
	TID     string         `json:"tid"`
	Code    int            `json:"code"`
	Msg     string         `json:"msg"`
}

type SwitchStatusResponse struct {
	Success bool `json:"success"`
	Status  bool `json:"status"`
}

type DeviceSpecificationResponse struct {
	Result  DeviceSpecification `json:"result"`
	Success bool                `json:"success"`
	T       int64               `json:"t"`
	TID     string              `json:"tid"`
	Code    int                 `json:"code"`
	Msg     string              `json:"msg"`
}

type DeviceSpecification struct {
	Functions []FunctionSpec `json:"functions"`
	Status    []StatusSpec   `json:"status"`
}

type FunctionSpec struct {
	Code   string `json:"code"`
	Type   string `json:"type"`
	Values string `json:"values"`
}

type StatusSpec struct {
	Code   string `json:"code"`
	Type   string `json:"type"`
	Values string `json:"values"`
}
