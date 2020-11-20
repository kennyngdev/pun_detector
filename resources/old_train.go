package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AlecAivazis/survey"
	"github.com/brentnd/go-snowboy"
	"github.com/common-nighthawk/go-figure"
)

func oldTrain() {
	//test

	//Trainer Logo
	fmt.Println("")
	fmt.Println("")
	myFigure := figure.NewColorFigure(" TRAIN YOUR MODEL!  ", "", "white", true)
	myFigure.Print()
	fmt.Println("")
	fmt.Println("")
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
	//End of User Prompt
	//Record Wav File

	//Post request to snowboy server
	t := snowboy.TrainRequest{
		Token:      "f2362919aedb3861e78ae1915cd35f9c0f84f8e5",
		Name:       modelName,
		Language:   modelLang,
		AgeGroup:   modelAgeGroup,
		Gender:     modelGender,
		Microphone: "standard USB mic",
	}
	t.AddWave(os.Args[3])
	t.AddWave(os.Args[4])
	t.AddWave(os.Args[5])
	pmdl, err := t.Train()
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(t.Name+".pmdl", pmdl, 0644)

}
