module github.com/fumiama/sqlite3

go 1.17

require (
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab
	modernc.org/libc v1.21.5
)

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.4.0 // indirect
)

replace github.com/remyoudompheng/bigfft => github.com/fumiama/bigfft v0.0.0-20211011143303-6e0bfa3c836b

retract [v1.16.0, v1.17.2] // https://gitlab.com/cznic/sqlite/-/issues/100

retract v1.19.0 // module source tree too large (max size is 524288000 bytes)
