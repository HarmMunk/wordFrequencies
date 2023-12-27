module example/wordfrequencies

go 1.21.1

replace example/passwords => ../passwords

require example/passwords v0.0.0-00010101000000-000000000000

require (
	github.com/pelletier/go-toml/v2 v2.1.1 // indirect
)
