module example

go 1.23.6

replace github.com/k-tsurumaki/fuselage => ../

replace github.com/k-tsurumaki/fuselage/middleware => ../middleware

require (
	github.com/k-tsurumaki/fuselage v0.0.0
	github.com/k-tsurumaki/fuselage/middleware v0.0.0-00010101000000-000000000000
)