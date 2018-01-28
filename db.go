package main

import (
	"log"
	"time"

	"github.com/go-pg/pg"
)

var db *pg.DB

// Movie all values
type Movie struct {
	ID          int64       `sql:"id,pk"        json:"id"`
	Section     string      `sql:"section"      json:"section"`
	Name        string      `sql:"name"         json:"name"`
	EngName     string      `sql:"eng_name"     json:"eng_name"`
	Year        int         `sql:"year"         json:"year"`
	Genre       []string    `sql:"genre"        json:"genre"        pg:",array"`
	Country     []string    `sql:"country"      json:"country"      pg:",array"`
	RawCountry  string      `sql:"raw_country"  json:"raw_country"`
	Director    []string    `sql:"director"     json:"director"     pg:",array"`
	Producer    []string    `sql:"producer"     json:"producer"     pg:",array"`
	Actor       []string    `sql:"actor"        json:"actor"        pg:",array"`
	Description string      `sql:"description"  json:"description"`
	Age         string      `sql:"age"          json:"age"`
	ReleaseDate string      `sql:"release_date" json:"release_date"`
	RussianDate string      `sql:"russian_date" json:"russian_date"`
	Duration    string      `sql:"duration"     json:"duration"`
	Kinopoisk   float64     `sql:"kinopoisk"    json:"kinopoisk"`
	IMDb        float64     `sql:"imdb"         json:"imdb"`
	Poster      string      `sql:"poster"       json:"poster"`
	PosterURL   string      `sql:"poster_url"   json:"poster_url"`
	NNM         float64     `sql:"-"            json:"nnm"`
	Torrent     []Torrent   `sql:"-"            json:"torrent"`
	CreatedAt   pg.NullTime `sql:"created_at"`
	UpdatedAt   pg.NullTime `sql:"updated_at"`
}

// Torrent all values
type Torrent struct {
	ID            int64       `sql:"id,pk"          json:"id"`
	MovieID       int64       `sql:"movie_id"       json:"movie_id"`
	DateCreate    string      `sql:"date_create"    json:"date_create"`
	Href          string      `sql:"href"           json:"href"`
	Torrent       string      `sql:"torrent"        json:"torrent"`
	Magnet        string      `sql:"magnet"         json:"magnet"`
	NNM           float64     `sql:"nnm"            json:"nnm"`
	SubtitlesType string      `sql:"subtitles_type" json:"subtitles_type"`
	Subtitles     string      `sql:"subtitles"      json:"subtitles"`
	Video         string      `sql:"video"          json:"video"`
	Quality       string      `sql:"quality"        json:"quality"`
	Resolution    string      `sql:"resolution"     json:"resolution"`
	Audio1        string      `sql:"audio1"         json:"audio1"`
	Audio2        string      `sql:"audio2"         json:"audio2"`
	Audio3        string      `sql:"audio3"         json:"audio3"`
	Translation   string      `sql:"translation"    json:"translation"`
	Size          int         `sql:"size"           json:"size"`
	Seeders       int         `sql:"seeders"        json:"seeders"`
	Leechers      int         `sql:"leechers"       json:"leechers"`
	CreatedAt     pg.NullTime `sql:"created_at"`
	UpdatedAt     pg.NullTime `sql:"updated_at"`
}

type search struct {
	ID      int64 `sql:"max"      json:"id"`
	MovieID int64 `sql:"movie_id" json:"movie_id"`
}

// initDB initialize database
func initDB(dbname string, user string, password string, logsql bool) {
	opt := pg.Options{
		User:     user,
		Password: password,
		Database: dbname,
	}
	db = pg.Connect(&opt)
	if logsql {
		db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, err := event.FormattedQuery()
			if err != nil {
				panic(err)
			}

			log.Printf("%s %s", time.Since(event.StartTime), query)
		})
	}
}

func getMovies(page int64) ([]Movie, int64, error) {
	var (
		movies   []Movie
		searches []search
	)

	count, err := db.Model(&Movie{}).Count()
	if err != nil {
		return nil, 0, err
	}
	_, err = db.Query(&searches, `SELECT max(t.id), t.movie_id FROM torrents AS t GROUP BY movie_id ORDER BY max(id) desc LIMIT ? OFFSET ?;`, 100, (page-1)*100)
	if err != nil {
		log.Println("Query search ", err)
		return nil, 0, err
	}
	for _, s := range searches {
		movie := getMovieByID(s.MovieID)
		torrents := getMovieTorrents(movie.ID)
		if len(torrents) > 0 {
			var i float64
			for _, t := range torrents {
				i = i + t.NNM
			}
			movie.Torrent = torrents
			movie.NNM = round(i/float64(len(torrents)), 1)
			movies = append(movies, movie)
		}
	}
	return movies, int64(count), nil
}

func getMovieByID(id int64) Movie {
	var movie Movie
	err := db.Model(&movie).Where("id = ?", id).Select()
	errchkmsg("getMovieByID", err)
	return movie
}

func getMovieTorrents(id int64) []Torrent {
	var torrents []Torrent
	err := db.Model(&torrents).Where("movie_id = ?", id).OrderExpr("seeders DESC").Select()
	errchkmsg("getMovieTorrents", err)
	return torrents
}
