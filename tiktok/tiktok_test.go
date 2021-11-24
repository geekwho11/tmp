package tiktok

import (
   "fmt"
   "github.com/89z/mech"
   "testing"
   "time"
)

const addr = "https://www.tiktok.com/@aamora_3mk/video/7028702876205632773"

func TestData(t *testing.T) {
   mech.Verbose = true
   for range [9]struct{}{} {
      vid, err := NewVideo(addr)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(vid.PlayAddr())
      time.Sleep(time.Second)
   }
}