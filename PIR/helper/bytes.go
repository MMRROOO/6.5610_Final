package main
import (
   "bufio"
   "fmt"
   "io"
   "os"
)

func fileNameToBytesArr(FileName string) []byte {
   file, err := os.Open(FileName)
   if err != nil {
	  fmt.Println(err)
	  return nil
   }
   defer file.Close()

   // Get the file size
   stat, err := file.Stat()
   if err != nil {
	  fmt.Println(err)
	  return nil
   }

   bs := make([]byte, stat.Size())
   _, err = bufio.NewReader(file).Read(bs)
   if err != nil && err != io.EOF {
	  fmt.Println(err)
	  return nil
   }
   return bs
}