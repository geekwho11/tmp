package youtube

import (
   "io"
   "testing"
)

const desc = "Provided to YouTube by Epitaph\n\nSnowflake · Kate Bush\n\n" +
"50 Words For Snow\n\n" +
"℗ Noble & Brite Ltd. trading as Fish People, under exclusive license to Anti Inc.\n\n" +
"Released on: 2011-11-22\n\nMusic  Publisher: Noble and Brite Ltd.\n" +
"Composer  Lyricist: Kate Bush\n\nAuto-generated by YouTube."

func TestAndroid(t *testing.T) {
   p, err := NewPlayer("MeJVWBSsPAY", Key, Android)
   if err != nil {
      t.Fatal(err)
   }
   f := p.StreamingData.AdaptiveFormats[0]
   if err := f.Write(io.Discard); err != nil {
      t.Fatal(err)
   }
}

func TestEmbed(t *testing.T) {
   p, err := NewPlayer("QWlNyzzwgcc", Key, Embed)
   if err != nil {
      t.Fatal(err)
   }
   f := p.StreamingData.AdaptiveFormats[0]
   if err := f.Write(io.Discard); err != nil {
      t.Fatal(err)
   }
}

func TestMweb(t *testing.T) {
   p, err := NewPlayer("XeojXq6ySs4", Key, Mweb)
   if err != nil {
      t.Fatal(err)
   }
   if p.Date() != "2020-11-05" {
      t.Fatalf("%+v\n", p)
   }
   if p.Description() != desc {
      t.Fatalf("%+v\n", p)
   }
   if p.Title() != "Snowflake" {
      t.Fatalf("%+v\n", p)
   }
   if p.VideoDetails.ViewCount == 0 {
      t.Fatalf("%+v\n", p)
   }
}
