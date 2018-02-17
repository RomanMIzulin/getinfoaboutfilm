package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"os"
	"log"
	"net/url"
	"io"
	"html/template"
	strings "strings"
)

const templ =`Number of films: {{.TotalNumber}} {{range .Items}} 
Title: {{.Title}}. Released: {{.Released}}. Runtime h.: {{.Runtime }}. Genre: {{.Genre}}. {{end}}`

const (
	DataRequest string ="http://www.omdbapi.com/?t="
)
var movieList =	template.Must(template.New("movielist").Parse(`
<h1>{{.TotalNumber}} Movies</h1>
<table>
<tr style='text-align: left'>
	<th>Title</th>
	<th>Released</th>
	<th>Runtime</th>
	<th>Genre</th>
</tr>
{{range .Items}}
<tr>
	<td><a href='{{.Poster}}'>{{.Title}}</td>
	<td>{{.Released}}</td>
	<td>{{.Runtime}}</td>
	<td>{{.Genre}}</td>
	<td></td>
</tr>
{{end}}
</table><hr>`))

type ResReqMovie struct {
	TotalNumber int
	Items *[]movie
}

type movie struct{
	Title string
	Poster string
	Released string
	Runtime string
	Genre string
}

//name []string because os.Args[1:] is slice
func GetMovieName(name []string) (*movie,error) {
	q := url.QueryEscape(strings.Join(name," "))
	resp,err := http.Get(DataRequest+q+"&apikey=8f4deded")
	if err != nil{
		return nil ,err
	}
	if resp.StatusCode != http.StatusOK{
		return nil,fmt.Errorf("failed with status %s", resp.Status)
	}
	var result movie
	if err := json.NewDecoder(resp.Body).Decode(&result);err != nil{
		return &result, nil //always return error - ignore
	}

	defer resp.Body.Close() //executes before jump out from function
	return &result, nil
}

func GetMoivePoster(name []string) (error){
	movie,err := GetMovieName(name)
	if err != nil {
		log.Println("failed to get movie from source")
	}

	data,err := http.Get(movie.Poster)
	if err != nil{
		return err
	}
	defer data.Body.Close()

	if data.StatusCode != http.StatusOK{
		log.Println("failed to get img with status ", data.StatusCode)
		return nil
	}

	out,err := os.Create("/where/to/storage/img.jpg")
	if err != nil {
		return err
	}
	defer out.Close()

	_,err= io.Copy(out,data.Body)
	if err != nil{
		return err
	}
	return nil
}

func result() ResReqMovie {
	var localResult []movie = []movie{}
	for i := range os.Args[1:] {

		tmp, err := GetMovieName(os.Args[i+1:i+2])
		if err != nil {
			log.Println("cant deal with args")
		}
		localResult = append(localResult, *tmp)
	}

	return ResReqMovie{len(localResult),&localResult}
}

var processedTmp =template.Must(template.New("films").Funcs(template.FuncMap{}).Parse(templ))

func main(){
	res:=result()
	if err := processedTmp.Execute(os.Stdout,res); err != nil{
		log.Fatal(err)
	}
	//"/home/romanmatveevskb161/page.html"
	w,errr := os.OpenFile("/home/romanmatveevskb161/page.html",os.O_CREATE|os.O_WRONLY, os.ModeDevice)
	if errr != nil {
		log.Fatal("cant open file")
	}

	defer w.Close()
	if err := movieList.Execute(w,res);err != nil{
		log.Fatal(err)
	}

	//if you want to download images
	//GetMoivePoster(os.Args[1:])
	//log.Println("well")
}