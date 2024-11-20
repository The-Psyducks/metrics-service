package main

import (
	"fmt"
	router "github.com/the-psyducks/metrics-service/src/router"
)

func main() {
	// Start the server
	rabbit, err := router.NewRabbitRouter()
	if err != nil {
		fmt.Println("Error creating rabbit router: ", err)
		return
	}
	rabbit.Run()

	fmt.Println("Creating web router")
	happyJuanma := `⊂_ヽ
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

	fmt.Println(happyJuanma)
	webRouter, err := router.CreateRouter()
	if err != nil {
		fmt.Println("Error creating router: ", err)
		return
	}

	if err := webRouter.Run(); err != nil {
		fmt.Println("Error starting router: ", err)
		return
	}

	//rabbit.Run()
}
