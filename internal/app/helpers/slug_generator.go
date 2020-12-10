package helpers

import (
	"strings"
)

var transcription = make(map[string]string, 34)

func init() {
	transcription[" "] = "_"
	transcription["а"] = "a"
	transcription["б"] = "b"
	transcription["в"] = "v"
	transcription["г"] = "g"
	transcription["д"] = "d"
	transcription["е"] = "e"
	transcription["ё"] = "yo"
	transcription["ж"] = "zh"
	transcription["з"] = "z"
	transcription["и"] = "i"
	transcription["й"] = "y"
	transcription["к"] = "k"
	transcription["л"] = "l"
	transcription["м"] = "m"
	transcription["н"] = "n"
	transcription["о"] = "o"
	transcription["п"] = "p"
	transcription["р"] = "r"
	transcription["с"] = "s"
	transcription["т"] = "t"
	transcription["у"] = "u"
	transcription["ф"] = "f"
	transcription["х"] = "kh"
	transcription["ц"] = "ts"
	transcription["ч"] = "ch"
	transcription["ш"] = "sh"
	transcription["щ"] = "sch"
	transcription["ъ"] = "y"
	transcription["ы"] = "y"
	transcription["ь"] = ""
	transcription["э"] = "e"
	transcription["ю"] = "yu"
	transcription["я"] = "ya"
}

var toLat = strings.NewReplacer(
	" ", "_",
	"а", "a",
	"б", "b",
	"в", "v",
	"г", "g",
	"д", "d",
	"е", "e",
	"ё", "yo",
	"ж", "zh",
	"з", "z",
	"и", "i",
	"й", "y",
	"к", "k",
	"л", "l",
	"м", "m",
	"н", "n",
	"о", "o",
	"п", "p",
	"р", "r",
	"с", "s",
	"т", "t",
	"у", "u",
	"ф", "f",
	"х", "kh",
	"ц", "ts",
	"ч", "ch",
	"ш", "sh",
	"щ", "sch",
	"ъе", "ye",
	"ъ", "",
	"ый", "iy",
	"ий", "iy",
	"ы", "y",
	"ь", "",
	"э", "e",
	"ю", "yu",
	"я", "ya",
)

func GenerateSlug(s string) string {
	return toLat.Replace(strings.ToLower(s))
}
