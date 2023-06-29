package process

type KatagoAnalysisResult struct {
	Id             string `json:"id"`
	IsDuringSearch bool   `json:"isDuringSearch"`
	MoveInfos      []struct {
		Lcb           float64  `json:"lcb"`
		Move          string   `json:"move"`
		Order         int      `json:"order"`
		Prior         float64  `json:"prior"`
		Pv            []string `json:"pv"`
		ScoreLead     float64  `json:"scoreLead"`
		ScoreMean     float64  `json:"scoreMean"`
		ScoreSelfplay float64  `json:"scoreSelfplay"`
		ScoreStdev    float64  `json:"scoreStdev"`
		Utility       float64  `json:"utility"`
		UtilityLcb    float64  `json:"utilityLcb"`
		Visits        int      `json:"visits"`
		Weight        float64  `json:"weight"`
		Winrate       float64  `json:"winrate"`
	} `json:"moveInfos"`
	Ownership []float64 `json:"ownership"`
	RootInfo  struct {
		CurrentPlayer   string  `json:"currentPlayer"`
		RawStScoreError float64 `json:"rawStScoreError"`
		RawStWrError    float64 `json:"rawStWrError"`
		RawVarTimeLeft  float64 `json:"rawVarTimeLeft"`
		ScoreLead       float64 `json:"scoreLead"`
		ScoreSelfplay   float64 `json:"scoreSelfplay"`
		ScoreStdev      float64 `json:"scoreStdev"`
		SymHash         string  `json:"symHash"`
		ThisHash        string  `json:"thisHash"`
		Utility         float64 `json:"utility"`
		Visits          int     `json:"visits"`
		Weight          float64 `json:"weight"`
		Winrate         float64 `json:"winrate"`
	} `json:"rootInfo"`
	TurnNumber int `json:"turnNumber"`
}

type MapResult struct {
	Id   string               `json:"id"`
	Data KatagoAnalysisResult `json:"data"`
	Msg  string               `json:"msg"`
}
