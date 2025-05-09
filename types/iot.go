package types

type (
	IoTReq struct {
		IotID       string  `json:"iot_id"`
		Temperature float32 `json:"temperature"`
	}
	IoTData struct {
		IoTReq
		Condition string `json:"condition"`
	}
)
