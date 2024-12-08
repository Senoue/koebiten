module github.com/sago35/koebiten

go 1.22.5

require (
	tinygo.org/x/drivers v0.28.1-0.20240825183126-07216d3051aa
	tinygo.org/x/tinydraw v0.4.0
	tinygo.org/x/tinyfont v0.4.0
)

replace github.com/sago35/koebiten/games/10secGo/10secGo => ./games/10secGo/10secGo

require github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
