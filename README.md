# Radiance65 App

[![CI](https://github.com/FuLygon/sr65-app/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/FuLygon/sr65-app/actions/workflows/ci.yaml)
[![GoReportCard](https://goreportcard.com/badge/github.com/FuLygon/sr65-app)](https://goreportcard.com/report/github.com/FuLygon/sr65-app)

Unofficial media conversion tool for **SR Studio Radiance65** keyboard screen.

## Download

[![release](https://img.shields.io/github/release/FuLygon/sr65-app.svg?style=flat)](https://github.com/FuLygon/sr65-app/releases)

## Supported media formats
- Image: `png`, `jpg/jpeg`, `gif`
- Video: `mp4`

## Build from source
Make sure [go](https://go.dev) is installed.

- Clone the repository
```bash
git clone https://github.com/FuLygon/sr65-app.git
cd sr65-app
```

- Copy `ffmpeg` binary to `embed/bin` directory if you want to embed it when building (optional)

- Build the app
```shell
go build
```

## Note
- For configuring keyboard feature such as keymap, macros, layout, etc... refer to [Vial](https://get.vial.today) instead.
- Video conversion heavily rely on `ffmpeg`. GIF conversion can be converted using built-in function or `ffmpeg` if `ffmpeg` is available.
- Released binary were embedded with `ffmpeg` by default for all OS with `amd64` architecture. For `arm64` architecture, only `linux` version were embedded. For `windows` and `darwin` version with `arm64` architecture, make sure to have `ffmpeg` installed and added to system `PATH`. This is only necessary for converting video or GIF using
`ffmpeg`.
- Antivirus software may flag released binary as malicious. This is a **false positive**. You can [build the app from source](#build-from-source) as an alternative.