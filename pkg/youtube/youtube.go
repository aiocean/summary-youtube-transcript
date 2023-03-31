package youtube

type Transcript struct {
	VideoId  string    `json:"videoID"`
	Segments []Segment `json:"segments"`
}

type Segment struct {
	Time string `json:"time"`
	Text string `json:"text"`
}

func Summary(transcript *Transcript) (*string, error) {
	var summary string
	for _, segment := range transcript.Segments {
		summary += segment.Text + " "
	}

	return &summary, nil
}
