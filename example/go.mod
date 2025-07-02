module example

go 1.23.6

replace (
    github.com/k-tsurumaki/fuselage => ../
    github.com/k-tsurumaki/fuselage/middleware => ../middleware/
)

require (
    github.com/k-tsurumaki/fuselage v1.0.0
    github.com/k-tsurumaki/fuselage/middleware v1.0.0
)