package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/AlecAivazis/survey"
	"github.com/brentnd/go-snowboy"
	"github.com/common-nighthawk/go-figure"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/fatih/color"
	"github.com/gordonklaus/portaudio"
)

type Sound struct {
	stream *portaudio.Stream
	data   []int16
}

func (s *Sound) Init() {
	inputChannels := 1
	outputChannels := 0
	sampleRate := 16000
	s.data = make([]int16, 1024)

	// initialize the audio recording interface
	err := portaudio.Initialize()
	if err != nil {
		fmt.Errorf("Error initialize audio interface: %s", err)
		return
	}

	// open the sound input stream for the microphone
	stream, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(s.data), s.data)
	if err != nil {
		fmt.Errorf("Error open default audio stream: %s", err)
		return
	}

	err = stream.Start()
	if err != nil {
		fmt.Errorf("Error on stream start: %s", err)
		return
	}

	s.stream = stream
}

func (s *Sound) Close() {
	s.stream.Close()
	portaudio.Terminate()
}

func (s *Sound) Read(p []byte) (int, error) {
	s.stream.Read()

	buf := &bytes.Buffer{}
	for _, v := range s.data {
		binary.Write(buf, binary.LittleEndian, v)
	}

	copy(p, buf.Bytes())
	return len(p), nil
}

func center(s string, w int) string {
	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(s))/2, s))
}

func main() {
	//App Logo
	fmt.Println("")
	fmt.Println("")
	myFigure := figure.NewColorFigure("PUN DETECTOR", "", "white", true)
	myFigure.Print()
	fmt.Println("")
	fmt.Println("")

	//App Title
	blue := color.New(color.FgBlue)
	blue.Add(color.Bold)
	whiteBackground := blue.Add(color.BgWhite)
	whiteBackground.Println("🔥🔥🔥Kenny's Pun Detector🔥🔥🔥")

	//User chooses to enter training mode or play mode
	userMode := ""
	modePrompt := &survey.Select{
		Message: "Please choose mode",
		Options: []string{"Listener", "Trainer"},
	}
	survey.AskOne(modePrompt, &userMode)
	if userMode == "Trainer" {
		Train()
	}

	//User Instruction
	whiteBackground.Println("Arrows to Move")
	whiteBackground.Println("Space to Choose Words")
	whiteBackground.Println("Enter to Submit")

	//list all alerts in folder
	var listOfAlerts []string
	files, err1 := ioutil.ReadDir("./resources/alert")
	if err1 != nil {
		log.Fatal(err1)
	}
	for _, file := range files {
		listOfAlerts = append(listOfAlerts, file.Name())
	}
	//list all models in folder
	var listOfFiles []string
	files, err := ioutil.ReadDir("./resources/models")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		listOfFiles = append(listOfFiles, f.Name())
	}
	//removing .DS_Store d
	listOfFiles = listOfFiles[1:]

	//User Prompt for Alert Sounds
	userAlert := ""
	alertPrompt := &survey.Select{
		Message: "Choose Your Alert Sound:",
		Options: listOfAlerts,
	}
	survey.AskOne(alertPrompt, &userAlert)
	//Multi-select User Prompt
	userChoice := []string{}
	selectPrompt := &survey.MultiSelect{
		Message: "🍕Select Puns/Words to Ban:",
		Options: listOfFiles,
	}
	survey.AskOne(selectPrompt, &userChoice)
	fmt.Print(userChoice)

	// open the mic
	mic := &Sound{}
	mic.Init()
	defer mic.Close()

	// open the snowboy detector
	d := snowboy.NewDetector("./resources/common.res")
	defer d.Close()
	//set up the Alert Sound
	f, err := os.Open("resources/alert/" + userAlert)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	//initialize speaker to play horn sound
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()
	//handlers for each hot words
	for i := 0; i < len(userChoice); i++ {
		d.HandleFunc(snowboy.NewHotword("./resources/models/"+userChoice[i], 0.46), func(string) {
			fmt.Println("You said the banned word!")
			shot := buffer.Streamer(0, buffer.Len())
			speaker.Play(shot)
		})
	}

	// display the detector's expected audio format
	sr, nc, bd := d.AudioFormat()
	fmt.Printf("sample rate=%d, num channels=%d, bit depth=%d\n", sr, nc, bd)

	// start detecting using the microphone
	d.ReadAndDetect(mic)

}
