package main

type ElasticResponse struct {
	Took     int  `Json:"took"`
	TimedOut bool `Json:"timed_out"`
	Shards struct {
		Total      int `Json:"total"`
		Successful int `Json:"successful"`
		Skipped    int `Json:"skipped"`
		Failed     int `Json:"failed"`
	} `Json:"_shards"`
	Hits struct {
		Total    float64       `Json:"total"`
		MaxScore float64       `Json:"max_score"`
		Hits     []interface{} `Json:"hits"`
	} `Json:"hits"`
}

func (er ElasticResponse) HitCount() float64 {
	return er.Hits.Total
}
