package server

type analysisScoreRequest struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

type demoRequest struct {
	SGF string `json:"sgf" binding:"required"`
}
