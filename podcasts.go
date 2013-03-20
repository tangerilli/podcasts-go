package main

import ("fmt"
        "flag"
        "os"
        "net/http"
        "log"
        "path/filepath"
        "encoding/xml"
        "strings"
        "strconv"
        "mime")

var movieExtensions []string = []string{"m4v", "avi", "mpg", "mpeg", "mkv"}

type Enclosure struct {
    XMLName xml.Name `xml:"enclosure"`
    Length int64 `xml:"length,attr"`
    Url string `xml:"url,attr"`
    Mimetype string `xml:"type,attr"`
}

type PodcastEntry struct {
    XMLName xml.Name `xml:"item"`
    Title string `xml:"title"`
    Link string `xml:"link"`
    Guid string `xml:"guid"`
    PubDate string `xml:"pubDate"`
    Size int64 `xml:"size"`
    Enclosure Enclosure `xml:"enclosure"`
}

type Channel struct {
    XMLName xml.Name `xml:"channel"`
    Title string `xml:"title"`
    Description string `xml:"description"`
    Link string `xml:"link"`
    Language string `xml:"language"`
    Entries []PodcastEntry
}

type Rss struct {
    XMLName xml.Name `xml:"rss"`
    Version string `xml:"version,attr"`
    Channel Channel
}

type PodcastHandler struct{
    dir string
    url string
    videoBaseUrl string
    title string
    description string
    language string
}

func FindFileTypes(path string, extensions []string) []string {
    var results []string
    for _, pattern := range extensions {
        r, err := filepath.Glob(filepath.Join(path, "*." + pattern))
        if err != nil {
            continue
        }
        results = append(results, r...)
    }
    return results
}

func (p PodcastHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "<?xml version=\"1.0\" encoding=\"UTF-8\" ?>\n")
    rss := Rss{Version:"2.0"}
    rss.Channel = Channel{Title:p.title, Description:p.description, Link:p.url, Language:p.language}
    movies := FindFileTypes(p.dir, movieExtensions)
    for _, movie := range(movies) {
        entry := PodcastEntry{Title:strings.Replace(filepath.Base(movie), filepath.Ext(movie), "", -1)}
        entry.Link = p.videoBaseUrl + filepath.Base(movie)
        entry.Guid = entry.Link
        finfo,err := os.Stat(movie)
        if err != nil {
            continue
        }
        entry.PubDate = finfo.ModTime().Format("Mon, 02 Jan 2006 15:04:05 +0000")
        entry.Size = finfo.Size()
        entry.Enclosure = Enclosure{Url:entry.Link, Length:finfo.Size(), Mimetype:mime.TypeByExtension(filepath.Ext(movie))}
        rss.Channel.Entries = append(rss.Channel.Entries, entry)
    }
    encoder := xml.NewEncoder(w)
    encoder.Encode(rss)
}

func usage() {
    fmt.Println("usage: podcasts <video directory> <podcast url>")
    flag.PrintDefaults()
    os.Exit(1)
}

func main() {
    flag.Usage = usage
    var port = flag.Int64("port", 4000, "HTTP port to listen on")
    var videoBaseUrl = flag.String("video_base_url", "", "Base URL that videos are served from")
    var title = flag.String("title", "Video Podcast", "The podcast title")
    var description = flag.String("description", "Autogenerated podcast", "The podcast description")
    var language = flag.String("language", "en-us", "The language code for the podcast (e.g. en-us)")
    flag.Parse()
    args := flag.Args()
    if len(args) < 2 {
        usage()
    }

    if *videoBaseUrl == "" {
        *videoBaseUrl = args[1]
    }
    if !strings.HasSuffix(*videoBaseUrl, "/") {
        *videoBaseUrl = *videoBaseUrl + "/"
    }

    fmt.Printf("Starting HTTP server on port %d\n", *port)
    handler := PodcastHandler{args[0], args[1], *videoBaseUrl, *title, *description, *language}
    log.Fatal(http.ListenAndServe(":" + strconv.FormatInt(*port, 10), handler))
}