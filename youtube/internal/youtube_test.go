package internal

import (
   "fmt"
   "github.com/89z/format"
   "github.com/89z/format/json"
   "net/http"
   "os"
   "testing"
   "time"
)

func TestContains(t *testing.T) {
   for name, version := range clients {
      _, err := newPlayer(name, version)
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}

func TestPlayer(t *testing.T) {
   for name, version := range clients {
      play, err := newPlayer(name, version)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(play.PlayabilityStatus.Status, name)
      time.Sleep(time.Second)
   }
}

func TestMweb(t *testing.T) {
   const name = "MWEB"
   version, err := newVersion("https://m.youtube.com", "iPad")
   if err != nil {
      t.Fatal(err)
   }
   if version != clients[name] {
      t.Fatal(name, version)
   }
}

func TestWebRemix(t *testing.T) {
   const name = "WEB_REMIX"
   version, err := newVersion("https://music.youtube.com", "Firefox/99")
   if err != nil {
      t.Fatal(err)
   }
   if version != clients[name] {
      t.Fatal(name, version)
   }
}

func TestWebCreator(t *testing.T) {
   const name = "WEB_CREATOR"
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   tok, err := format.Open[token](cache, "mech/youtube.json")
   if err != nil {
      t.Fatal(err)
   }
   req, err := http.NewRequest("GET", "https://studio.youtube.com", nil)
   if err != nil {
      t.Fatal(err)
   }
   req.URL.RawQuery = "approve_browser_access=true"
   req.Header.Set("Authorization", "Bearer " + tok.Access_Token)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   sep := []byte(`"client":`)
   var client struct {
      ClientVersion string
   }
   if err := json.Decode(res.Body, sep, &client); err != nil {
      t.Fatal(err)
   }
   if client.ClientVersion != clients[name] {
      t.Fatal(name, client.ClientVersion)
   }
}

func TestWebUnplugged(t *testing.T) {
   const name = "WEB_UNPLUGGED"
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   tok, err := format.Open[token](cache, "mech/youtube.json")
   if err != nil {
      t.Fatal(err)
   }
   req, err := http.NewRequest("GET", "https://tv.youtube.com", nil)
   if err != nil {
      t.Fatal(err)
   }
   req.Header.Set("Authorization", "Bearer " + tok.Access_Token)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   sep := []byte(`"client":`)
   var client struct {
      ClientVersion string
   }
   if err := json.Decode(res.Body, sep, &client); err != nil {
      t.Fatal(err)
   }
   if client.ClientVersion != clients[name] {
      t.Fatal(name, client.ClientVersion)
   }
}

func TestWeb(t *testing.T) {
   const name = "WEB"
   version, err := newVersion("https://www.youtube.com", "")
   if err != nil {
      t.Fatal(err)
   }
   if version != clients[name] {
      t.Fatal(name, version)
   }
}

func TestWebEmbeddedPlayer(t *testing.T) {
   const name = "WEB_EMBEDDED_PLAYER"
   version, err := newVersion("https://www.youtube.com/embed/MIchMEqVwvg", "")
   if err != nil {
      t.Fatal(err)
   }
   if version != clients[name] {
      t.Fatal(name, version)
   }
}

func TestTvhtml5(t *testing.T) {
   const name = "TVHTML5"
   version, err := newVersion(
      "https://www.youtube.com/tv",
      "Mozilla/5.0 (ChromiumStylePlatform) Cobalt/Version",
   )
   if err != nil {
      t.Fatal(err)
   }
   if version != clients[name] {
      t.Fatal(name, version)
   }
}

func TestWebKids(t *testing.T) {
   const name = "WEB_KIDS"
   version, err := newVersion("https://www.youtubekids.com", "Firefox/99")
   if err != nil {
      t.Fatal(err)
   }
   if version != clients[name] {
      t.Fatal(name, version)
   }
}
