package main

// PostalCode represents a single Indonesian postal code record.
type PostalCode struct {
	Province  string  `json:"province"`
	Regency   string  `json:"regency"`
	District  string  `json:"district"`
	Village   string  `json:"village"`
	Code      int     `json:"code"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"elevation"`
	Timezone  string  `json:"timezone"`
}

// PostalCodeResult is a PostalCode with optional computed fields.
type PostalCodeResult struct {
	Province  string  `json:"province"`
	Regency   string  `json:"regency"`
	District  string  `json:"district"`
	Village   string  `json:"village"`
	Code      int     `json:"code"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"elevation"`
	Timezone  string  `json:"timezone"`
	Distance  float64 `json:"distance,omitempty"`
}

// APIResponse is the standard API envelope.
type APIResponse struct {
	StatusCode int         `json:"statusCode"`
	Code       string      `json:"code"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

// newOK creates a 200 OK response.
func newOK(data interface{}) APIResponse {
	return APIResponse{StatusCode: 200, Code: "OK", Data: data}
}

// newBadRequest creates a 400 response.
func newBadRequest(msg string) APIResponse {
	return APIResponse{StatusCode: 400, Code: "BAD_REQUEST", Message: msg}
}

// newNotFound creates a 404 response.
func newNotFound(msg string) APIResponse {
	return APIResponse{StatusCode: 404, Code: "NOT_FOUND", Message: msg}
}

// newUnauthorized creates a 401 response.
func newUnauthorized(msg string) APIResponse {
	return APIResponse{StatusCode: 401, Code: "UNAUTHORIZED", Message: msg}
}

// newInternalError creates a 500 response.
func newInternalError() APIResponse {
	return APIResponse{StatusCode: 500, Code: "INTERNAL_SERVER_ERROR", Message: "Please contact the developer."}
}
