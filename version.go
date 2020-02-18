package hamgo

import "fmt"

const VERSION = "1.1.0"

func printLogo() {
	logo := "##################################\n#          HamGo V%s          #\n##################################\n"
	fmt.Printf(logo, VERSION)
}
