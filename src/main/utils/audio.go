package utils

import (
  "github.com/timshannon/go-openal/openal"
  "io/ioutil"
  "time"
  //"math"
  "fmt"
)
func Play() {
  device := openal.OpenDevice("")
  defer device.CloseDevice()

  context := device.CreateContext()
  defer context.Destroy()
  context.Activate()

  source := openal.NewSource()
  defer source.Pause()
  source.SetLooping(false)

  buffer := openal.NewBuffer()

  if err := openal.Err(); err != nil {
    fmt.Println(err)
    return
  }

  data, err := ioutil.ReadFile("assets/audio/oborona.wav")
  if err != nil {
    fmt.Println(err)
    return
  }

  buffer.SetData(openal.FormatMono16, data, 44100)

  source.SetBuffer(buffer)
  source.Play()
  for source.State() == openal.Playing {
    // loop long enough to let the wave file finish
    time.Sleep(time.Millisecond * 100)
  }
  source.Delete()
  fmt.Println("sound played")
  // Output: sound played
}


func main() {

  Play()
}
