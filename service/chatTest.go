package service

// func TestChat() {
// 	t := time.NewTicker(3 * time.Second)
// 	for {
// 		<-t.C
// 		// fmt.Println("TestChat")
// 		if len(testChat) == 2 {
// 			pkId := "pk" + strconv.FormatInt(time.Now().UnixMilli(), 10)
// 			if err := twoConnect("start", testChat[0], testChat[0], testChat[1], pkId); err != nil {
// 				fmt.Println("twoConnect err: ", err)
// 			}
// 			go func(roomId1, roomId2, pkId string) {
// 				t1 := time.NewTicker(5 * time.Second)
// 				for {
// 					<-t1.C
// 					if err := twoConnect("heartbeat", roomId1, roomId1, roomId2, pkId); err != nil {
// 						fmt.Println("twoConnect err: ", err)
// 					}
// 				}
// 			}(testChat[0], testChat[1], pkId)
// 			testChat = make([]string, 0)
// 		}
// 	}
// }
