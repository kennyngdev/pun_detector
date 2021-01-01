package main

/*
  #include <stdio.h>
  #include <unistd.h>
  #include <termios.h>
  char getch(){
      char ch = 0;
      struct termios old = {0};
      fflush(stdout);
      if( tcgetattr(0, &old) < 0 ) perror("tcsetattr()");
      old.c_lflag &= ~ICANON;
      old.c_lflag &= ~ECHO;
      old.c_cc[VMIN] = 1;
      old.c_cc[VTIME] = 0;
      if( tcsetattr(0, TCSANOW, &old) < 0 ) perror("tcsetattr ICANON");
      if( read(0, &ch,1) < 0 ) perror("read()");
      old.c_lflag |= ICANON;
      old.c_lflag |= ECHO;
      if(tcsetattr(0, TCSADRAIN, &old) < 0) perror("tcsetattr ~ICANON");
      return ch;
  }
*/
import "C"

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey"
	"github.com/brentnd/go-snowboy"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/gordonklaus/portaudio"
	wave "github.com/zenwerk/go-wave"
)

func errCheck(err error) {

	if err != nil {
		panic(err)
	}
}

func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func Train() {
	//Trainer Logo
	fmt.Println("")
	fmt.Println("")
	myFigure := figure.NewColorFigure(" TRAIN YOUR MODEL!  ", "", "white", true)
	myFigure.Print()
	fmt.Println("")
	fmt.Println("")
	//User Instruction
	blue2 := color.New(color.FgBlue)
	blue2.Add(color.Bold)
	whiteBackground2 := blue2.Add(color.BgWhite)
	whiteBackground2.Println("Arrows to Move")
	whiteBackground2.Println("Enter to Choose")
	//User Prompt Asking Model Name, Language, AgeGroup, Gender
	//Ask Model Name
	modelName := ""
	prompt := &survey.Input{
		Message: "What word/pun are you Training? This is going to be the model name (Enter English only)",
	}
	survey.AskOne(prompt, &modelName)
	//Language
	userLang := ""
	var modelLang snowboy.Language
	langPrompt := &survey.Select{
		Message: "Choose Language:",
		Options: []string{"English", "French", "Italian", "Chinese", "Japanese", "Korean", "Hindi", "Spanish"},
	}
	survey.AskOne(langPrompt, &userLang)

	if userLang == "English" {
		modelLang = snowboy.LanguageEnglish
	} else if userLang == "French" {
		modelLang = snowboy.LanguageFrench
	} else if userLang == "Chinese" {
		modelLang = snowboy.LanguageChinese
	} else if userLang == "Japanese" {
		modelLang = snowboy.LanguageJapanese
	} else if userLang == "Korean" {
		modelLang = snowboy.LanguageKorean
	} else if userLang == "Italian" {
		modelLang = snowboy.LanguageItalian
	} else if userLang == "Hindi" {
		modelLang = snowboy.LanguageHindi
	} else if userLang == "Spanish" {
		modelLang = snowboy.LanguageSpanish
	}
	//Age Group
	ageGroup := ""
	var modelAgeGroup snowboy.AgeGroup
	agePrompt := &survey.Select{
		Message: "What Age Group are you in? (for model training purpose only)",
		Options: []string{"0s", "10s", "20s", "30s", "40s", "50s", "60s or Above"},
	}
	survey.AskOne(agePrompt, &ageGroup)

	if ageGroup == "0s" {
		modelAgeGroup = snowboy.AgeGroup0s
	} else if ageGroup == "10s" {
		modelAgeGroup = snowboy.AgeGroup10s
	} else if ageGroup == "20s" {
		modelAgeGroup = snowboy.AgeGroup20s
	} else if ageGroup == "30s" {
		modelAgeGroup = snowboy.AgeGroup30s
	} else if ageGroup == "40s" {
		modelAgeGroup = snowboy.AgeGroup40s
	} else if ageGroup == "50s" {
		modelAgeGroup = snowboy.AgeGroup50s
	} else {
		modelAgeGroup = snowboy.AgeGroup60plus
	}
	//Gender
	gender := ""
	var modelGender snowboy.Gender
	genderPrompt := &survey.Select{
		Message: "Are u a:(for model training purpose only)",
		Options: []string{"Male", "Female"},
	}
	survey.AskOne(genderPrompt, &gender)
	if gender == "Male" {
		modelGender = snowboy.GenderMale
	} else {
		modelGender = snowboy.GenderFemale
	}

	voidInput := ""
	readyPrompt := &survey.Input{
		Message: "Let's start recording! When you are ready please press ENTER",
	}
	survey.AskOne(readyPrompt, &voidInput)
	//End of User Prompt
	//Record Wav File
	fmt.Println(modelAgeGroup, modelGender, modelLang)
	//Post request to snowboy server\
	ticker := []string{
		"-",
		"\\",
		"/",
		"|",
	}
	fmt.Println("Recording.Please say a word then press ENTER.")
	//create the three files
	waveFile1, err := os.Create("./user_recording/1.wav")
	errCheck(err)

	waveFile2, err := os.Create("./user_recording/2.wav")
	errCheck(err)

	waveFile3, err := os.Create("./user_recording/3.wav")
	errCheck(err)

	// record
	inputChannels := 1
	outputChannels := 0
	sampleRate := 44100
	framesPerBuffer := make([]byte, 64)

	// setup Wave file writer

	param1 := wave.WriterParam{
		Out:           waveFile1,
		Channel:       inputChannels,
		SampleRate:    sampleRate,
		BitsPerSample: 8, // if 16, change to WriteSample16()
	}

	param2 := wave.WriterParam{
		Out:           waveFile2,
		Channel:       inputChannels,
		SampleRate:    sampleRate,
		BitsPerSample: 8, // if 16, change to WriteSample16()
	}

	param3 := wave.WriterParam{
		Out:           waveFile3,
		Channel:       inputChannels,
		SampleRate:    sampleRate,
		BitsPerSample: 8, // if 16, change to WriteSample16()
	}

	//  PortAudio init and open stream
	portaudio.Initialize()

	//first run

	stream1, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(framesPerBuffer), framesPerBuffer)
	errCheck(err)

	waveWriter1, err := wave.NewWriter(param1)
	errCheck(err)

	go func() {
		key := C.getch()
		fmt.Println()
		fmt.Println("Cleaning up...")
		if key == 10 {
			waveWriter1.Close()
			stream1.Close()
			fmt.Println("Second Recording")
			stream2, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(framesPerBuffer), framesPerBuffer)
			errCheck(err)
			waveWriter2, err := wave.NewWriter(param2)
			errCheck(err)

			go func() {
				key := C.getch()
				fmt.Println()
				fmt.Println("Cleaning up...")
				if key == 10 {
					waveWriter2.Close()
					stream2.Close()
					fmt.Println("Third Recording")
					stream3, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(framesPerBuffer), framesPerBuffer)
					errCheck(err)
					waveWriter3, err := wave.NewWriter(param3)
					errCheck(err)
					//handler for ending recording
					go func() {
						key := C.getch()
						fmt.Println()
						fmt.Println("Cleaning up...")
						fmt.Println("Recording done!")
						fmt.Println("Downloading Model...")
						if key == 10 {
							waveWriter3.Close()
							stream3.Close()
							portaudio.Terminate()
							//POST request to snowboy API
							cmd := exec.Command("python", "training_service.py",
								"./user_recording/1_new.wav", "./user_recording/2_new.wav",
								"./user_recording/3_new.wav", modelName)
							stdout, err := cmd.StdoutPipe()
							if err != nil {
								panic(err)
							}
							stderr, err := cmd.StderrPipe()
							if err != nil {
								panic(err)
							}
							err = cmd.Start()
							if err != nil {
								panic(err)
							}

							go copyOutput(stdout)
							go copyOutput(stderr)
							cmd.Wait()

							os.Exit(0)
						}
					}()
					// start recording for 3rd file
					errCheck(stream3.Start())
					for {
						errCheck(stream3.Read())
						fmt.Printf("\rRecording...[%v]", ticker[rand.Intn(len(ticker)-1)])
						// write to wave file
						_, err := waveWriter3.Write([]byte(framesPerBuffer)) // WriteSample16 for 16 bits
						errCheck(err)
					}
					errCheck(stream3.Stop())
				}
			}()
			errCheck(stream2.Start())
			for {
				errCheck(stream2.Read())
				fmt.Printf("\rRecording...[%v]", ticker[rand.Intn(len(ticker)-1)])
				// write to wave file
				_, err := waveWriter2.Write([]byte(framesPerBuffer)) // WriteSample16 for 16 bits
				errCheck(err)
			}
			errCheck(stream2.Stop())
		}
	}()
	// start recording for 1st file
	errCheck(stream1.Start())
	for {
		errCheck(stream1.Read())
		fmt.Printf("\rRecording...[%v]", ticker[rand.Intn(len(ticker)-1)])
		// write to wave file
		_, err := waveWriter1.Write([]byte(framesPerBuffer)) // WriteSample16 for 16 bits
		errCheck(err)
	}
	errCheck(stream1.Stop())

}
