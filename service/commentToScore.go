package service

func CommentToScoreInit(commentScore map[string]float64, winCoinComment map[int64]string, commentCoin map[string]int64, nodeInte map[int64]int64) {
	commentToScore = commentScore
	winCoinToComment = winCoinComment
	commentToCoin = commentCoin
	nodeIdToIntegral = nodeInte
}

// 查询评论是否在commentToScore中
func QueryCommentToScore(comment string) float64 {
	if commentToScore[comment] != 0 {
		return commentToScore[comment]
	}
	return 0
}
