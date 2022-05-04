module github.com/kuromitsu0104/go-interpreter

go 1.16

replace (
	local.packages/lexer => ./lexer
	local.packages/repl => ./repl
	local.packages/token => ./token
)

require (
	local.packages/lexer v0.0.0-00010101000000-000000000000 // indirect
	local.packages/repl v0.0.0-00010101000000-000000000000
	local.packages/token v0.0.0-00010101000000-000000000000 // indirect
)
