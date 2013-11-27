package util

import (
  "os"
  "crypto/md5"
  "fmt"
  "io"
)

func Md5sum(path string) (string, error) {
  file, err := os.Open(path)
  if err != nil {
    return "", err
  }
  defer file.Close()

  buf := make([]byte, 1024)
  hash := md5.New()
  for {
    n, err := file.Read(buf)
    if err != nil && err != io.EOF {
      panic(err)
    }
    if n == 0 {
      break
    }
    if _, err := io.WriteString(hash, string(buf[:n])); err != nil {
      panic(err)
    }
  }
  return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
