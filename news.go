package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	//"fmt"
	"html/template"
	"log"
	"net/http"

	//"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Info struct{
	Weather []struct{
	Main string `json:"main"`
	Description string `json:"description"`
	}`json:"weather"`
	Name string `json:"name"`
	Main struct{
		Kelvin float64 `json:"temp"`
		Celsius int64
	}`json:"main"`
}
type News struct{ 
    Articles []struct{
	Author string `json:"author"`
	Title string `json:"title"`
	Description string `json:"description"`
	Url string `json:"url"`
	UrlToImage string `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
}`json:"articles"`
}
type NewStruct struct{
	NewConfirmed uint `json:"NewConfirmed"`
	TotalConfirmed uint `json:"TotalConfirmed"`
	NewDeaths uint `json:"NewDeaths"`
	TotalDeaths uint `json:"TotalDeaths"`
	NewRecovered uint `json:"NewRecovered"`
	TotalRecovered uint `json:"TotalRecovered"`
}
type WorldWideCases struct{
	Global struct{
		NewConfirmed uint `json:"NewConfirmed"`
		TotalConfirmed uint `json:"TotalConfirmed"`
		NewDeaths uint `json:"NewDeaths"`
		TotalDeaths uint `json:"TotalDeaths"`
		NewRecovered uint `json:"NewRecovered"`
		TotalRecovered uint `json:"TotalRecovered"`
	}`json:"Global"`
}
type CountryCases struct{
	Confirmed uint `json:"Confirmed"`
	Deaths uint `json:"Deaths"`
	Recovered uint `json:"Recovered"`
	Active uint `json:"Active"`

}
type Mainstruct struct{
	Slice []Info
	Climate string
	Description string
	Day time.Weekday
	Date int
	Month time.Month
	Slice2 []News
	Widget1 []News
	Slice3 []NewStruct 
	Slice4 []CountryCases
}
var tpl *template.Template
var NKey,WKey string
func init(){
	err:=godotenv.Load()
	if err!=nil{
		log.Fatal(err)
	}
	NKey=os.Getenv("NEWS_KEY")
	WKey=os.Getenv("WEATHER_KEY")
tpl=template.Must(template.ParseGlob("templates/*.html"))
}
func main(){
	m:=mux.NewRouter().StrictSlash(true)
	m.PathPrefix("/images/").Handler(http.StripPrefix("/images/",http.FileServer(http.Dir("./images"))))
	m.PathPrefix("/css/").Handler(http.StripPrefix("/css/",http.FileServer(http.Dir("./css"))))
	m.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/",http.FileServer(http.Dir("./fonts"))))
	m.PathPrefix("/js/").Handler(http.StripPrefix("/js/",http.FileServer(http.Dir("./js"))))
	m.PathPrefix("/sass/").Handler(http.StripPrefix("/sass/",http.FileServer(http.Dir("./sass"))))
	m.HandleFunc("/",homehandler)
	m.HandleFunc("/weather",weatherhandler)
	m.HandleFunc("/weather2",weatherhandler2)
	m.HandleFunc("/news",newshandler)
	m.HandleFunc("/internationalnews",internationalnewshandler)
	m.HandleFunc("/category",categoryhandler)
	m.HandleFunc("/covid19",covidhandler)
	//m.HandleFunc("/covid19two",covidhandler2)
	m.HandleFunc("/article",articlehandler)
	http.ListenAndServe(":8080",m)
}
func homehandler(w http.ResponseWriter,r *http.Request){
	resp,err:=http.Get("http://newsapi.org/v2/top-headlines?country=in&apiKey="+NKey)
	if err!=nil{
		log.Fatal(err)
	}
	var m News
	defer resp.Body.Close()
	//fmt.Println(resp)
	var data Mainstruct
	err2:=json.NewDecoder(resp.Body).Decode(&m)
	if err2!=nil{
		log.Fatal(err2)
	}
	//fmt.Println(m.Articles[0])
	var n News
	n.Articles=nil
	for i:=0;i<3;i++{
		n.Articles=append(n.Articles,m.Articles[i])
	}
	data.Slice2=nil
	data.Slice2=append(data.Slice2,n)
	timeframe:=time.Now()
	data.Day=timeframe.Weekday()
	data.Date=timeframe.Day()
	data.Month=timeframe.Month()
	tpl.ExecuteTemplate(w,"home.html",data)
}
func weatherhandler(w http.ResponseWriter,r *http.Request){
	resp1,err1:=http.Get("http://newsapi.org/v2/top-headlines?country=in&apiKey="+NKey)
	if err1!=nil{
		log.Fatal(err1)
	}
	var m Mainstruct
	var m1 News
	defer resp1.Body.Close()
	//fmt.Println(resp)
	//var data1 Mainstruct
	err3:=json.NewDecoder(resp1.Body).Decode(&m1)
	if err3!=nil{
		log.Fatal(err3)
	}
	//fmt.Println(m.Articles[0])
	var n News
	n.Articles=nil
	for i:=0;i<3;i++{
		n.Articles=append(n.Articles,m1.Articles[i])
	}
	m.Slice2=nil
	m.Slice2=append(m.Slice2,n)
	timeframe:=time.Now()
	m.Day=timeframe.Weekday()
	m.Date=timeframe.Day()
	m.Month=timeframe.Month()
	tpl.ExecuteTemplate(w,"index.html",m)
}
func weatherhandler2(w http.ResponseWriter,r *http.Request){
	city:=r.FormValue("city")
	resp,err:=http.Get("http://api.openweathermap.org/data/2.5/weather?appid="+WKey+"?q="+city)
	if err!=nil{
		log.Fatal(err)
	}
	defer resp.Body.Close()
	//fmt.Println(resp)
	var data Info
	err2:=json.NewDecoder(resp.Body).Decode(&data)
	if err2!=nil{
		log.Fatal(err2)
	}
	
	data.Main.Celsius=int64(data.Main.Kelvin-273)
	//fmt.Println(data)
	var m Mainstruct
	m.Slice=nil
	m.Slice=append(m.Slice,data)
	m.Climate=data.Weather[0].Main
	m.Description=data.Weather[0].Description
	timeframe:=time.Now()
	m.Day=timeframe.Weekday()
	m.Date=timeframe.Day()
	m.Month=timeframe.Month()
	resp1,err1:=http.Get("http://newsapi.org/v2/top-headlines?country=in&apiKey="+NKey)
	if err1!=nil{
		log.Fatal(err1)
	}
	var m1 News
	defer resp1.Body.Close()
	//fmt.Println(resp)
	//var data1 Mainstruct
	err3:=json.NewDecoder(resp1.Body).Decode(&m1)
	if err3!=nil{
		log.Fatal(err3)
	}
	//fmt.Println(m.Articles[0])
	var n News
	n.Articles=nil
	m.Slice2=nil
	for i:=0;i<3;i++{
		n.Articles=append(n.Articles,m1.Articles[i])
	}
	m.Slice2=append(m.Slice2,n)
	
	tpl.ExecuteTemplate(w,"index.html",m)
}
func newshandler(w http.ResponseWriter,r *http.Request){
	resp,err:=http.Get("http://newsapi.org/v2/top-headlines?country=in&apiKey="+NKey)
	if err!=nil{
		log.Fatal(err)
	}
	var m News
	defer resp.Body.Close()
	//fmt.Println(resp)
	var data Mainstruct
	err2:=json.NewDecoder(resp.Body).Decode(&m)
	if err2!=nil{
		log.Fatal(err2)
	}
	//fmt.Println(m.Articles[0])
	var n News
	var Wid News
	for i:=0;i<10;i++{
		n.Articles=append(n.Articles,m.Articles[i])
		if i<=5{
			Wid.Articles=append(Wid.Articles,m.Articles[i])
		}
	}
	data.Slice2=nil
	data.Widget1=nil
	data.Slice2=append(data.Slice2,n)
	data.Widget1=append(data.Widget1,Wid)
	timeframe:=time.Now()
	data.Day=timeframe.Weekday()
	data.Date=timeframe.Day()
	data.Month=timeframe.Month()
	tpl.ExecuteTemplate(w,"news.html",data)
}
func internationalnewshandler(w http.ResponseWriter,r *http.Request){
	resp,err:=http.Get("http://newsapi.org/v2/top-headlines?country=in&apiKey="+NKey)
	if err!=nil{
		log.Fatal(err)
	}
	var m News
	defer resp.Body.Close()
	//fmt.Println(resp)
	var data Mainstruct
	err2:=json.NewDecoder(resp.Body).Decode(&m)
	if err2!=nil{
		log.Fatal(err2)
	}
	//fmt.Println(m.Articles[0])
	var n News
	var Wid News
	n.Articles=nil
	Wid.Articles=nil
	for i:=0;i<10;i++{
		n.Articles=append(n.Articles,m.Articles[i])
		if i<=5{
			Wid.Articles=append(Wid.Articles,m.Articles[i])
		}
	}
	data.Slice2=nil
	data.Widget1=nil
	data.Slice2=append(data.Slice2,n)
	data.Widget1=append(data.Widget1,Wid)
	timeframe:=time.Now()
	data.Day=timeframe.Weekday()
	data.Date=timeframe.Day()
	data.Month=timeframe.Month()
	tpl.ExecuteTemplate(w,"news.html",data)
}
func categoryhandler(w http.ResponseWriter,r *http.Request){
	category:=r.FormValue("q")
	resp,err:=http.Get("http://newsapi.org/v2/top-headlines?country=in&apiKey="+NKey+"&category="+category)
	if err!=nil{
		log.Fatal(err)
	}
	var m News
	defer resp.Body.Close()
	//fmt.Println(resp)
	var data Mainstruct
	err2:=json.NewDecoder(resp.Body).Decode(&m)
	if err2!=nil{
		log.Fatal(err2)
	}
	//fmt.Println(m.Articles[0])
	var n News
	var Wid News
	n.Articles=nil
	Wid.Articles=nil
	for i:=0;i<10;i++{
		n.Articles=append(n.Articles,m.Articles[i])
		if i<=5{
			Wid.Articles=append(Wid.Articles,m.Articles[i])
		}
	}
	data.Slice2=append(data.Slice2,n)
	data.Widget1=append(data.Widget1,Wid)
	timeframe:=time.Now()
	data.Day=timeframe.Weekday()
	data.Date=timeframe.Day()
	data.Month=timeframe.Month()
	tpl.ExecuteTemplate(w,"news.html",data)	
}
func covidhandler(w http.ResponseWriter,r *http.Request){
	resp,err:=http.Get("https://api.covid19api.com/summary")
	if err!=nil{
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var cases WorldWideCases
	err2:=json.NewDecoder(resp.Body).Decode(&cases)
	if err2!=nil{
		log.Fatal(err2)
	}
	//fmt.Println(cases)
	var n NewStruct
	var m Mainstruct
	n=cases.Global
	m.Slice3=nil
	m.Slice3=append(m.Slice3,n)
	timeframe:=time.Now()
	m.Day=timeframe.Weekday()
	m.Date=timeframe.Day()
	m.Month=timeframe.Month()
	tpl.ExecuteTemplate(w,"covid19.html",m)	
}
/*func covidhandler2(w http.ResponseWriter,r *http.Request){
	country:=r.FormValue("city")
	resp,err:=http.Get("https://api.covid19api.com/live/country/"+country)
	if err!=nil{
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var m Mainstruct
	err2:=json.NewDecoder(resp.Body).Decode(&m.Slice4)
	if err2!=nil{
		log.Fatal(err2)
	}
	fmt.Println(m.Slice4)
	//var n NewStruct
	//var m Mainstruct
	//n=cases.Global
	//m.Slice3=append(m.Slice3,n)
	timeframe:=time.Now()
	m.Day=timeframe.Weekday()
	m.Date=timeframe.Day()
	m.Month=timeframe.Month()
	tpl.ExecuteTemplate(w,"covid192.html",m)
}*/
func articlehandler(w http.ResponseWriter,r *http.Request){
	article:=r.FormValue("data")
	items:=strings.Fields(article)
	url:="http://newsapi.org/v2/everything?q="+items[0]
	resp,err:=http.Get(url+"&language=en&sortBy=publishedAt&apiKey="+NKey)
	if err!=nil{
		log.Fatal(err)
	}
	var m News
	defer resp.Body.Close()
	fmt.Println(resp)
	var data Mainstruct
	err2:=json.NewDecoder(resp.Body).Decode(&m)
	if err2!=nil{
		log.Fatal(err2)
	}
	fmt.Println(m.Articles)
	var n News
	var Wid News
	for i:=0;i<10;i++{
		n.Articles=append(n.Articles,m.Articles[i])
		if i<=5{
			Wid.Articles=append(Wid.Articles,m.Articles[i])
		}
	}
	data.Slice2=nil
	data.Widget1=nil
	data.Slice2=append(data.Slice2,n)
	data.Widget1=append(data.Widget1,Wid)
	timeframe:=time.Now()
	data.Day=timeframe.Weekday()
	data.Date=timeframe.Day()
	data.Month=timeframe.Month()
	tpl.ExecuteTemplate(w,"article.html",data)
}