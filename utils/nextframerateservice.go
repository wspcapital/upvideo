package utils

var (
	initialResolution = 1278
	frameRates        = []int{20, 21, 22, 23, 24, 25, 26, 27, 28, 29}
)

type NextFrameRateService struct {
	FrameRate      int
	Resolution     int
	frameRateIndex int
}

func NewNextFrameRateService() *NextFrameRateService {
	return &NextFrameRateService{FrameRate: frameRates[0], Resolution: initialResolution, frameRateIndex: 0}
}

func (this *NextFrameRateService) NextFrameRate() *NextFrameRateService {
	this.frameRateIndex++
	if this.frameRateIndex > len(frameRates)-1 {
		this.frameRateIndex = 0
		this.Resolution += 16
	}

	this.FrameRate = frameRates[this.frameRateIndex]

	return this
}
