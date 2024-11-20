package main

import (
	"fmt"
	router "github.com/the-psyducks/metrics-service/src/router"
)

func main() {
	// Start the server
	rabbit, _ := router.NewRabbitRouter()
	rabbit.Run()

	fmt.Println("Creating web router")
	asciiArt := `⊂_ヽ
　 ＼＼
　　 ＼( ͡° ͜ʖ ͡°)
　　　 >　⌒ヽ
　　　/ 　 へ＼
　　 /　　/　＼＼
　　 ﾚ　ノ　　 ヽ_つ
　　/　/
　 /　/|
　(　(ヽ
　|　|、＼
　| 丿 ＼ ⌒)
　| |　　) /
ノ )　　Lﾉ
(_／`

	fmt.Println(asciiArt)
	router, err := router.CreateRouter()
	if err != nil {
		fmt.Println("Error creating router: ", err)
		return
	}

	if err := router.Run(); err != nil {
		fmt.Println("Error starting router: ", err)
		return
	}

	//rabbit.Run()
}
