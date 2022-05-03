module example.com/monkey

go 1.16

replace (
  example.com/lexer => ./lexer
  example.com/token => ./token
)

require (
  "example.com/lexer" v0.0.0
  "example.com/token" v0.0.0
)
