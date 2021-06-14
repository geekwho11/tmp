// RARBG
package rarbg

import (
   "encoding/json"
   "fmt"
   "github.com/89z/mech/ocr.space"
   "io"
   "net/http"
   "os"
   "path/filepath"
   "regexp"
   "time"
)

const (
   AJAXPHP = "/threat_defence_ajax.php"
   DefencePHP = "/threat_defence.php"
   Origin = "http://rarbg.to"
   Sleep = 4 * time.Second
   invert = "\x1b[7m"
   reset = "\x1b[m"
)

// This returns solution to the Captcha at the given path. After this, you will
// want to call IamHuman.
func Solve(php string) (solve string, err error) {
   fmt.Println(invert, "GET", reset, Origin + php)
   res, err := http.Get(Origin + php)
   if err != nil { return "", err }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return "", fmt.Errorf("status %v", res.Status)
   }
   capt := filepath.Join(os.TempDir(), "captcha.png")
   file, err := os.Create(capt)
   if err != nil { return "", err }
   file.ReadFrom(res.Body)
   // need to close before opening again, not after return
   file.Close()
   img, err := ocr.NewImage(capt)
   if err != nil { return "", err }
   return img.ParsedResults[0].ParsedText, nil
}

type Defence struct {
   CID string
   I string
   SK string
}

// This is the entrypoint into getting the SKT cookie, should you need to do
// that. After this you will want to call ThreatCaptcha.
func NewDefence() (*Defence, error) {
   fmt.Println(invert, "GET", reset, Origin + DefencePHP)
   res, err := http.Get(Origin + DefencePHP)
   if err != nil { return nil, err }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("status %v", res.Status)
   }
   body, err := io.ReadAll(res.Body)
   if err != nil { return nil, err }
   d := new(Defence)
   // CID
   re := regexp.MustCompile(" value_c = '([^']+)'")
   cid := re.FindSubmatch(body)
   if cid == nil {
      return nil, fmt.Errorf("findSubmatch %v", re)
   }
   d.CID = string(cid[1])
   // I
   re = regexp.MustCompile(" value_i = '([^']+)'")
   i := re.FindSubmatch(body)
   if i == nil {
      return nil, fmt.Errorf("findSubmatch %v", re)
   }
   d.I = string(i[1])
   // SK
   re = regexp.MustCompile(" value_sk = '([^']+)'")
   sk := re.FindSubmatch(body)
   if sk == nil {
      return nil, fmt.Errorf("findSubmatch %v", re)
   }
   d.SK = string(sk[1])
   return d, nil
}

// This saves the SKT cookie to the Cache folder for later use.
func (d Defence) IamHuman(id, solve string) error {
   req, err := http.NewRequest("GET", Origin + DefencePHP, nil)
   if err != nil { return err }
   val := req.URL.Query()
   val.Set("defence", "2")
   val.Set("cid", d.CID)
   val.Set("i", d.I)
   val.Set("sk", d.SK)
   val.Set("solve_string", solve)
   val.Set("captcha_id", id)
   req.URL.RawQuery = val.Encode()
   fmt.Println(invert, "GET", reset, req.URL)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil { return err }
   defer res.Body.Close()
   for _, c := range res.Cookies() {
      if c.Name != "skt" { continue }
      cache, err := os.UserCacheDir()
      if err != nil { return err }
      os.Mkdir(cache + "/mech", os.ModeDir)
      file, err := os.Create(cache + "/mech/skt.json")
      if err != nil { return err }
      defer file.Close()
      enc := json.NewEncoder(file)
      enc.SetIndent("", " ")
      return enc.Encode(c)
   }
   return http.ErrNoCookie
}

// This returns path to Captcha image, as well as Captcha ID. After this, you
// will want to call Solve.
func (d Defence) ThreatCaptcha() (php string, id string, err error) {
   if err := d.threatDefenceAJAX(); err != nil {
      return "", "", err
   }
   req, err := http.NewRequest("GET", Origin + DefencePHP, nil)
   if err != nil { return "", "", err }
   val := req.URL.Query()
   val.Set("defence", "2")
   val.Set("cid", d.CID)
   val.Set("i", d.I)
   val.Set("sk", d.SK)
   req.URL.RawQuery = val.Encode()
   time.Sleep(Sleep)
   fmt.Println(invert, "GET", reset, req.URL)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil { return "", "", err }
   defer res.Body.Close()
   body, err := io.ReadAll(res.Body)
   if err != nil { return "", "", err }
   // captchaPHP
   re := regexp.MustCompile(`/threat_captcha\.php\?[^"]+`)
   captchaPHP := re.Find(body)
   if captchaPHP == nil {
      return "", "", fmt.Errorf("find %v", re)
   }
   // captchaID
   re = regexp.MustCompile(`"captcha_id" value="([^"]+)"`)
   captchaID := re.FindSubmatch(body)
   if captchaID == nil {
      return "", "", fmt.Errorf("findSubmatch %v", re)
   }
   // return
   return string(captchaPHP), string(captchaID[1]), nil
}

func (d Defence) threatDefenceAJAX() error {
   req, err := http.NewRequest("GET", Origin + AJAXPHP, nil)
   if err != nil { return err }
   val := req.URL.Query()
   val.Set("cid", d.CID)
   val.Set("i", d.I)
   val.Set("sk", d.SK)
   req.URL.RawQuery = val.Encode()
   fmt.Println(invert, "GET", reset, req.URL)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil { return err }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return fmt.Errorf("status %v", res.Status)
   }
   return nil
}