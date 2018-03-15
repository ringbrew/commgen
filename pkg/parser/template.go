package parser

const CommentTemplate = "/*\n" +
	"\t{{.FuncName}}\n" +
	"\tParams:\n" +
	"{{range $index,$data := .FuncParams}}" +
	"\t\t@{{$data.Name}}({{$data.Type}}):\n" +
	"{{end}}" +
	"\tReturns:\n" +
	"{{range $index,$data := .FuncReturn}}" +
	"\t\t@{{$data.Name}}({{$data.Type}}):\n" +
	"{{end}}" +
	"*/\n"
