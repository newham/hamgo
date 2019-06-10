package hamgo

import "fmt"

const VERSION = "1.0.2"

func printLogo() {
	logo := `      ___           ___           ___           ___           ___     
     /\__\         /\  \         /\__\         /\  \         /\  \    
    /:/  /        /::\  \       /::|  |       /::\  \       /::\  \   
   /:/__/        /:/\:\  \     /:|:|  |      /:/\:\  \     /:/\:\  \  
  /::\  \ ___   /::\~\:\  \   /:/|:|__|__   /:/  \:\  \   /:/  \:\  \ 
 /:/\:\  /\__\ /:/\:\ \:\__\ /:/ |::::\__\ /:/__/_\:\__\ /:/__/ \:\__\
 \/__\:\/:/  / \/__\:\/:/  / \/__/~~/:/  / \:\  /\ \/__/ \:\  \ /:/  /
      \::/  /       \::/  /        /:/  /   \:\ \:\__\    \:\  /:/  / 
      /:/  /        /:/  /        /:/  /     \:\/:/  /     \:\/:/  /  
     /:/  /        /:/  /        /:/  /       \::/  /       \::/  /   
     \/__/         \/__/         \/__/         \/__/         \/__/
` + "hamgo V%s\n\n"
	fmt.Printf(logo, VERSION)
}
