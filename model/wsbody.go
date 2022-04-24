package model

type WsBody struct {
	Code        string `json:"code"`
	CompileArgs string `json:"compile_args"`
	Input       string `json:"input"`
	Interrupt   bool   `json:"interrupt" default:"false"`
}
