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
	"fmt"
	"math/rand"
	"os"

	"github.com/gordonklaus/portaudio"
	wave "github.com/zenwerk/go-wave"
)

func errCheck(err error) {

	if err != nil {
		panic(err)
	}
}

func record() {
	// recording in progress ticker. From good old DOS days.
	ticker := []string{
		"-",
		"\\",
		"/",
		"|",
	}
	fmt.Println("Recording.Please say a word then press ESC to exit.")
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

	//first run !!

	stream1, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(framesPerBuffer), framesPerBuffer)
	errCheck(err)

	waveWriter1, err := wave.NewWriter(param1)
	errCheck(err)

	// rand.Seed(time.Now().UnixNano())
	go func() {
		key := C.getch()
		fmt.Println()
		fmt.Println("Cleaning up...1")
		if key == 10 {
			waveWriter1.Close()
			stream1.Close()
			fmt.Println("do i run?1")
			stream2, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(framesPerBuffer), framesPerBuffer)
			errCheck(err)
			waveWriter2, err := wave.NewWriter(param2)
			errCheck(err)

			go func() {
				key := C.getch()
				fmt.Println()
				fmt.Println("Cleaning up...2")
				if key == 10 {
					waveWriter2.Close()
					stream2.Close()
					fmt.Println("do i run2?")
					stream3, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(framesPerBuffer), framesPerBuffer)
					errCheck(err)
					waveWriter3, err := wave.NewWriter(param3)
					errCheck(err)

					//handler for ending recording
					go func() {
						key := C.getch()
						fmt.Println()
						fmt.Println("Cleaning up...")
						fmt.Println("do i run3?")
						fmt.Println("Recording done!")
						if key == 10 {
							waveWriter3.Close()
							stream3.Close()
							portaudio.Terminate()
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

// sox -t wav ./user_recording/1.wav -t wav -r 16000 -b 16 -e signed-integer -c 1 ./user_recording/1_new.wav
// sox -t wav ./user_recording/2.wav -t wav -r 16000 -b 16 -e signed-integer -c 1 ./user_recording/2_new.wav
// sox -t wav ./user_recording/3.wav -t wav -r 16000 -b 16 -e signed-integer -c 1 ./user_recording/3_new.wav

// python training_service.py ./user_recording/1_new.wav ./user_recording/2_new.wav ./user_recording/3_new.wav saved_model.pmdl
