package common

// config data models
type config struct {
	App      app
	Database database
	Kafka    kafkaconfig
	Redis    redis
}

type redis struct {
	ConnectionString string
}

type database struct {
	Server   string
	Port     string
	Database string
	User     string
	Password string
	Debug    bool
}

type kafkaconfig struct {
	Ip    string
	Port  string
	Topic string
}

type app struct {
	Name        string
	Description string
	Host        string
	Proxy       string
	Environment string
}

// Common API models for using across different APIs

type MetaData struct {
	Filters    []Filter   `json:"filters"`
	Pagination Pagination `json:"pagination"`
	Order      Order      `json:"order"`
}

type Filter struct {
	Field    string `json:"field,required"`
	Value    string `json:"value,required"`
	Operator string `json:"operator,required"`
}

type Order struct {
	OrderBy   string `json:"orderBy,omitempty"`
	OrderType string `json:"orderType,omitempty"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
}

// GetAllRequest dataTables Request Types
type GetAllRequest struct {
	Filters []Filter `json:"filters,omitempty"`
	Order
	Limit  int    `json:"limit,required"`
	Offset int    `json:"offset,required"`
	Query  string `json:"query,omitempty"`
}
