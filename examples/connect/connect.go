package main

import (
	"log"

	"github.com/EtienneBruines/gobci"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

var (
	timePeriod = float32(5) // seconds
)

func main() {
	// Connecting
	log.Println("Trying to connect to localhost:1972 ...")

	conn, err := gobci.Connect("localhost:1972")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Getting header information
	log.Println("Requesting header data ...")

	header, err := conn.GetHeader()
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Getting samples
		log.Println("Requesting sample data ...")

		var amountOfSamples uint32 = 100
		if header.NSamples < amountOfSamples {
			log.Fatal("Not enough samples avialable")
		}

		samples, err := conn.GetData(0, 0)
		if err != nil {
			log.Fatal(err)
		}

		// Visualizing the channels
		channels := make([]channelXYer, header.NChannels)
		for _, sample := range samples {
			for i := uint32(0); i < header.NChannels; i++ {
				channels[i].Values = append(channels[i].Values, sample[i])
			}
		}
		for chIndex := range channels {
			channels[chIndex].freq = header.SamplingFrequency
		}

		log.Println("Plotting samples ...")
		plt, err := plot.New()
		plotutil.AddLinePoints(plt,
			"CH0", plotter.XYer(channels[0]))
		//"CH1", plotter.XYer(channels[1]),
		//"CH2", plotter.XYer(channels[2]))

		log.Println("Saving plot to output.jpg ...")
		plt.Save(10*vg.Inch, 5*vg.Inch, "output.jpg")

	}

	log.Println("Done")
}

type channelXYer struct {
	Values []float64
	freq   float32
}

func (c channelXYer) Len() int {
	if max := c.freq * timePeriod; float32(len(c.Values)) < max {
		return len(c.Values)
	} else {
		return int(max)
	}
}

func (c channelXYer) XY(index int) (x, y float64) {
	if max := c.freq * timePeriod; float32(len(c.Values)) < max {
		return float64(index), c.Values[index]
	} else {
		return float64((float32(index) - max) / c.freq), c.Values[len(c.Values)-int(max)+index]
	}
}
