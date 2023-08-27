package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type PageData struct {
	Color, ColorBack string
	Text             []string
}

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("thetext")
	fileName := r.FormValue("chose")
	_, error := os.Stat(fileName + ".txt")
	userColor := strings.ToLower(r.FormValue("color"))
	backgroundColor := CheckBackgroundColor(userColor)
	indexTemplate, _ := template.ParseFiles("template/index.html")
	if (!CheckLetter(text) || text == "" || !CheckColor(userColor)) && r.Method == "POST" {
		w.WriteHeader(http.StatusBadRequest)
		http.ServeFile(w, r, "./template/400.html")
		return
	} else if (os.IsNotExist(error) || len(r.FormValue("thetext")) > 2000) && r.Method != "GET" {
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "template/500.html")
		return
	}
	textInASCII := serveIndex(text, fileName)
	pageData := PageData{
		Text:      textInASCII,
		Color:     userColor,
		ColorBack: backgroundColor,
	}
	if r.URL.Path == "/style.css" {
		http.ServeFile(w, r, "./template/style.css")
		return
	} else if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "./template/404.html")
		return
	} else if r.Method == "GET" || r.Method == "POST" {
		err := indexTemplate.Execute(w, pageData)
		if err != nil {
			fmt.Print(err)
		}
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		http.ServeFile(w, r, "./template/400.html") // should be 400
		return
	} 
}

func serveIndex(text, filename string) []string {
	var Text []string
	WordsInArr := strings.Split(text, "\r\n")
	var Words []string
	for l := 0; l < len(WordsInArr); l++ {
		var Words [][]string
		Text1 := strings.ReplaceAll(WordsInArr[l], "\\t", "   ")
		if Text1 != "" {
			for j := 0; j < len(Text1); j++ {
				Words = append(Words, ReadLetter(Text1[j], filename))
			}
			for x := 0; x < 8; x++ {
				Lines := ""
				for n := 0; n < len(Words); n++ {
					Lines += Words[n][x]
				}
				Text = append(Text, Lines)
			}
		} else {
			Text = append(Text, "\n")
		}
	}
	return append(Words, strings.Join(Text, "\n"))
}

func ReadLetter(Text1 byte, fileName string) []string {
	var Letter []string
	ReadFile, _ := os.Open(fileName + ".txt")
	FileScanner := bufio.NewScanner(ReadFile)
	stop := 1
	i := 0
	letterLength := (int(Text1)-32)*9 + 2
	for FileScanner.Scan() {
		i++
		if i >= letterLength {
			stop++
			Letter = append(Letter, FileScanner.Text())
			if stop > 8 {
				break
			}
		}
	}
	ReadFile.Close()
	return Letter
}

func CheckLetter(s string) bool {
	WordsInArr := strings.Split(s, "\r\n")
	for l := 0; l < len(WordsInArr); l++ {
		for g := 0; g < len(WordsInArr[l]); g++ {
			if WordsInArr[l][g] > 126 || WordsInArr[l][g] < 32 {
				return false
			}
		}
	}
	return true
}

func CheckColor(userValue string) bool {
	Colors := []string{"aliceblue", "antiquewhite", "aqua", "aquamarine", "azure", "beige", "bisque", "black", "blanchedalmond", "blue", "blueviolet", "brown",
		"burlywood", "cadetblue", "chartreuse", "chocolate", "coral", "cornflowerblue", "cornsilk", "crimson", "cyan",
		"darkblue", "darkcyan", "darkgoldenrod", "darkgray", "darkgreen", "darkkhaki", "darkmagenta", "darkolivegreen",
		"darkorange", "darkorchid", "darkred", "darksalmon", "darkseagreen", "darkslateblue", "darkslategray",
		"darkturquoise", "darkviolet", "deeppink", "deepskyblue", "dimgray", "dodgerblue", "firebrick", "floralwhite",
		"forestgreen", "fuchsia", "gainsboro", "ghostwhite", "gold", "goldenrod", "gray", "green", "greenyellow",
		"honeydew", "hotpink", "indianred", "indigo", "ivory", "khaki", "lavender", "lavenderblush", "lawngreen",
		"lemonchiffon", "lightblue", "lightcoral", "lightcyan", "lightgoldenrodyellow", "lightgray", "lightgreen",
		"lightpink", "lightsalmon", "lightseagreen", "lightskyblue", "lightslategray", "lightsteelblue", "lightyellow",
		"lime", "limegreen", "linen", "magenta", "maroon", "mediumaquamarine", "mediumblue", "mediumorchid", "mediumpurple",
		"mediumseagreen", "mediumslateblue", "mediumspringgreen", "mediumturquoise", "mediumvioletred", "midnightblue",
		"mintcream", "mistyrose", "moccasin", "navajowhite", "navy", "oldlace", "olive", "olivedrab", "orange", "orangered",
		"orchid", "palegoldenrod", "palegreen", "paleturquoise", "palevioletred", "papayawhip", "peachpuff", "peru", "pink", "plum",
		"powderblue", "purple", "red", "rosybrown", "royalblue", "saddlebrown", "salmon", "sandybrown", "seagreen", "seashell",
		"sienna", "silver", "skyblue", "slateblue", "slategray", "snow", "springgreen", "steelblue", "tan", "teal", "thistle",
		"tomato", "turquoise", "violet", "wheat", "white", "whitesmoke", "yellow", "yellowgreen"}
	for _, color := range Colors {
		if color == userValue {
			return true
		} else if strings.Index(userValue, "#") == 0 && len(userValue) == 7 {
			for i := 1; i <= 6; i++ {
				if (userValue[i] >= '0' && userValue[i] <= '9') || (userValue[i] >= 'a' && userValue[i] <= 'f') {
				} else {
					return false
				}
			}
			return true
		}
	}
	return false
}

func CheckBackgroundColor(userColor string) string {
	Colors := []string{"whitesmoke", "#f5f5f5", "seashell", "#fff5ee", "papayawhip", "#ffefd5",
		"oldlace", "#fdf5e6", "linen", "#faf0e6", "lightgoldenrodyellow", "#fafad2", "lemonchiffon", "#fffacd",
		"lavenderblush", "#fff0f5", "cornsilk", "#fff8dc", "blanchedalmond", "#ffebcd", "beige", "#f5f5dc", "antiquewhite",
		"#faebd7", "#f1f0e8"}
	for _, color := range Colors {
		if color == userColor {
			return "#f1aeb2"
		}
	}
	return "#f1f0e8"
}
