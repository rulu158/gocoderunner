package main

type CodePOST struct {
	Code string `json:"code"`
}

type Response struct {
	ID     string `json:"id,omitempty"`
	Error  bool   `json:"error"`
	Result string `json:"result"`
}
