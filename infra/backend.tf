terraform {
  backend "gcs" {
    bucket = "hamaguchi-family-linebot-tfstate"
    prefix = "default"
  }
}
