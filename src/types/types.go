package types

// Request type
type Request struct {
	URL string
}

// Result type
type Result struct {
	URL string
}

// FetcherResult type
type FetcherResult struct {
	URL string
}

// Channels type
type Channels struct {
	URLs          chan string
	FetcherResult chan FetcherResult
	FetcherDone   chan int
	UploaderDone  chan int
	Req           chan Request
	Res           chan Result
	Quit          chan int
}

// NewChannels returns new ref
func NewChannels() *Channels {
	return &Channels{
		URLs:          make(chan string, 10),
		FetcherResult: make(chan FetcherResult, 10),
		FetcherDone:   make(chan int, 10),
		UploaderDone:  make(chan int, 10),
		Req:           make(chan Request, 10),
		Res:           make(chan Result, 10),
		Quit:          make(chan int, 10),
	}
}

// DetailFetcherResult type
type DetailFetcherResult struct {
	ID   string
	Item *Museum
}

// DetailChannels type
type DetailChannels struct {
	FetcherResult chan DetailFetcherResult
	FetcherDone   chan int
	UploaderDone  chan int
}

// NewDetailChannels returns new ref
func NewDetailChannels() *DetailChannels {
	return &DetailChannels{
		FetcherResult: make(chan DetailFetcherResult, 10),
		FetcherDone:   make(chan int, 10),
		UploaderDone:  make(chan int, 10),
	}
}
