package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Used to color text
const (
	White  = "\033[97m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Reset  = "\033[0m"
)

// Condition Icons
const rain string = `
                  ⢀⣤⣶⡿⣽⣻⣞⣶⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
           ⠀⠀⠀⡀⢀⣴⣾⣿⣿⣿⢿⣿⡿⣿⣿⣾⣯⢷⣂⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
        ⠀⠀⣤⣴⣾⣿⣿⣟⣾⣿⡿⣟⣿⣿⣾⣿⣟⣿⣿⣿⣿⣽⢣⠄⠀⠀⠀⠀⠀⠀⠀⠀
     ⠀⠀⢀⣰⣿⣿⣿⣿⣿⣿⢿⡿⣽⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣯⠻⡜⠀⠀⠀⠀⠀⠀⠀⠀
      ⣀⡶⣯⢻⢏⣿⣳⣾⣿⣿⣟⣿⣿⣿⣿⣿⣿⣿⣿⣿⡟⣏⢷⣺⣝⡾⣤⣄⡀⠀⠀⠀⠀⠀
⠀⠀⠀⢀⣴⣾⣿⣿⣿⣯⣾⣷⣿⣿⣿⣿⣿⢿⡻⣟⡻⢿⢿⡛⣝⣶⣿⣾⣿⣻⢿⣿⣿⣾⡝⡶⣄⠀⠀⠀
⠀⢀⡾⣿⣿⣿⢿⣽⣿⣿⣿⣿⣿⣿⣿⣹⣾⣯⣷⣿⣿⣞⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣾⡵⣩⠞⣀⠀
⠀⢾⣽⡿⣟⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⢿⣯⢿⣱⢻⡄⠠
⣘⢮⣟⣿⣽⡿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⣟⢯⢞⣯⢖⠠
⠌⠹⣾⡹⣞⣿⣽⣻⣿⣿⣿⣿⣿⣿⣿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢿⡽⣟⣯⢟⣾⡹⢎⡗⣎⠒
⠐⣌⢳⡽⣹⡞⣧⡟⣷⢫⣷⣻⢭⣻⡽⢯⣿⢿⣯⢿⣿⣟⡿⣿⢯⣟⣿⢾⡽⢯⢿⣹⢮⡟⡾⢡⣋⠼⡐⠂
⠀⠈⢆⠳⡌⣝⢲⡝⣯⢳⡞⣵⢫⣗⡻⣝⢾⣛⢾⣛⠾⣭⢻⡭⣟⢾⡹⢧⠻⣝⢮⡝⡾⠼⣍⠳⢌⠂⠁⠀
⠀⠀⠀⠑⠘⢌⠣⠞⡰⢣⠚⡌⢧⡘⡱⢊⠞⢨⠓⡌⠳⣡⠓⡜⢌⠲⡉⢏⡱⢎⠲⡙⡜⠣⠎⠑⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠁⠈⠈⠀⠀⠀⠁⠈⠀⠀⠀⠁⠀⠀⠀⠈⠀⠁⠈⠀⠈⠁⠈⠀⠁⠀⠀⠀⠀⠀⠀
` + Blue + `⠀⠀⠀⠀⠀⢀⣤⠆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣤⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀` + Reset + `
` + Blue + `⠀⠀⠀⠀⠀⠺⠙⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣴⡟⠀⠀⠀⠀⠀⠀⠀⠀⢐⠻⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀` + Reset + `
` + Blue + `⠀⠀⠀⠀⠀⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠐⠈⠄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠀⠀⠀⠀⠀ ⠀⠀⠀⠀⠀⠀` + Reset + `
` + Blue + `⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣶⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀ ⠀⠀⠀⠀⠀⠀` + Reset + `
` + Blue + `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠋⠆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠚⡽⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ⠀⠀⠀⠀⠀⠀` + Reset + `
` + Blue + `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ⠀⠀⠀⠀⠀⠀` + Reset + `
	`

const snow string = `
                  ⢀⣤⣶⡿⣽⣻⣞⣶⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
           ⠀⠀⠀⡀⢀⣴⣾⣿⣿⣿⢿⣿⡿⣿⣿⣾⣯⢷⣂⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
        ⠀⠀⣤⣴⣾⣿⣿⣟⣾⣿⡿⣟⣿⣿⣾⣿⣟⣿⣿⣿⣿⣽⢣⠄⠀⠀⠀⠀⠀⠀⠀⠀
     ⠀⠀⢀⣰⣿⣿⣿⣿⣿⣿⢿⡿⣽⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣯⠻⡜⠀⠀⠀⠀⠀⠀⠀⠀
      ⣀⡶⣯⢻⢏⣿⣳⣾⣿⣿⣟⣿⣿⣿⣿⣿⣿⣿⣿⣿⡟⣏⢷⣺⣝⡾⣤⣄⡀⠀⠀⠀⠀⠀
⠀⠀⠀⢀⣴⣾⣿⣿⣿⣯⣾⣷⣿⣿⣿⣿⣿⢿⡻⣟⡻⢿⢿⡛⣝⣶⣿⣾⣿⣻⢿⣿⣿⣾⡝⡶⣄⠀⠀⠀
⠀⢀⡾⣿⣿⣿⢿⣽⣿⣿⣿⣿⣿⣿⣿⣹⣾⣯⣷⣿⣿⣞⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣾⡵⣩⠞⣀⠀
⠀⢾⣽⡿⣟⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⢿⣯⢿⣱⢻⡄⠠
⣘⢮⣟⣿⣽⡿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⣟⢯⢞⣯⢖⠠
⠌⠹⣾⡹⣞⣿⣽⣻⣿⣿⣿⣿⣿⣿⣿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢿⡽⣟⣯⢟⣾⡹⢎⡗⣎⠒
⠐⣌⢳⡽⣹⡞⣧⡟⣷⢫⣷⣻⢭⣻⡽⢯⣿⢿⣯⢿⣿⣟⡿⣿⢯⣟⣿⢾⡽⢯⢿⣹⢮⡟⡾⢡⣋⠼⡐⠂
⠀⠈⢆⠳⡌⣝⢲⡝⣯⢳⡞⣵⢫⣗⡻⣝⢾⣛⢾⣛⠾⣭⢻⡭⣟⢾⡹⢧⠻⣝⢮⡝⡾⠼⣍⠳⢌⠂⠁⠀
⠀⠀⠀⠑⠘⢌⠣⠞⡰⢣⠚⡌⢧⡘⡱⢊⠞⢨⠓⡌⠳⣡⠓⡜⢌⠲⡉⢏⡱⢎⠲⡙⡜⠣⠎⠑⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠁⠈⠈⠀⠀⠀⠁⠈⠀⠀⠀⠁⠀⠀⠀⠈⠀⠁⠈⠀⠈⠁⠈⠀⠁⠀⠀⠀⠀⠀⠀
` + White + `        .                  .            ` + Reset + `
` + White + `         ❄️      .         ❄             ` + Reset + `
` + White + `             .                          ` + Reset + `
` + White + `      ❄️   .         .  ❄                ` + Reset + `
` + White + `               ❄️ .          .           ` + Reset + `
` + White + `       .  .             .               ` + Reset + `
` + White + `    .           .      ❄️                ` + Reset + `
` + White + `      ❄️    .                            ` + Reset + `
` + White + `       .                                ` + Reset + `
	`

const cloud string = `
                  ⢀⣤⣶⡿⣽⣻⣞⣶⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
           ⠀⠀⠀⡀⢀⣴⣾⣿⣿⣿⢿⣿⡿⣿⣿⣾⣯⢷⣂⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
        ⠀⠀⣤⣴⣾⣿⣿⣟⣾⣿⡿⣟⣿⣿⣾⣿⣟⣿⣿⣿⣿⣽⢣⠄⠀⠀⠀⠀⠀⠀⠀⠀
     ⠀⠀⢀⣰⣿⣿⣿⣿⣿⣿⢿⡿⣽⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣯⠻⡜⠀⠀⠀⠀⠀⠀⠀⠀
      ⣀⡶⣯⢻⢏⣿⣳⣾⣿⣿⣟⣿⣿⣿⣿⣿⣿⣿⣿⣿⡟⣏⢷⣺⣝⡾⣤⣄⡀⠀⠀⠀⠀⠀
⠀⠀⠀⢀⣴⣾⣿⣿⣿⣯⣾⣷⣿⣿⣿⣿⣿⢿⡻⣟⡻⢿⢿⡛⣝⣶⣿⣾⣿⣻⢿⣿⣿⣾⡝⡶⣄⠀⠀⠀
⠀⢀⡾⣿⣿⣿⢿⣽⣿⣿⣿⣿⣿⣿⣿⣹⣾⣯⣷⣿⣿⣞⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣾⡵⣩⠞⣀⠀
⠀⢾⣽⡿⣟⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⢿⣯⢿⣱⢻⡄⠠
⣘⢮⣟⣿⣽⡿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⣟⢯⢞⣯⢖⠠
⠌⠹⣾⡹⣞⣿⣽⣻⣿⣿⣿⣿⣿⣿⣿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢿⡽⣟⣯⢟⣾⡹⢎⡗⣎⠒
⠐⣌⢳⡽⣹⡞⣧⡟⣷⢫⣷⣻⢭⣻⡽⢯⣿⢿⣯⢿⣿⣟⡿⣿⢯⣟⣿⢾⡽⢯⢿⣹⢮⡟⡾⢡⣋⠼⡐⠂
⠀⠈⢆⠳⡌⣝⢲⡝⣯⢳⡞⣵⢫⣗⡻⣝⢾⣛⢾⣛⠾⣭⢻⡭⣟⢾⡹⢧⠻⣝⢮⡝⡾⠼⣍⠳⢌⠂⠁⠀
⠀⠀⠀⠑⠘⢌⠣⠞⡰⢣⠚⡌⢧⡘⡱⢊⠞⢨⠓⡌⠳⣡⠓⡜⢌⠲⡉⢏⡱⢎⠲⡙⡜⠣⠎⠑⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠁⠈⠈⠀⠀⠀⠁⠈⠀⠀⠀⠁⠀⠀⠀⠈⠀⠁⠈⠀⠈⠁⠈⠀⠁⠀⠀⠀⠀⠀⠀
	`

const storm string = `
                  ⢀⣤⣶⡿⣽⣻⣞⣶⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
           ⠀⠀⠀⡀⢀⣴⣾⣿⣿⣿⢿⣿⡿⣿⣿⣾⣯⢷⣂⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
        ⠀⠀⣤⣴⣾⣿⣿⣟⣾⣿⡿⣟⣿⣿⣾⣿⣟⣿⣿⣿⣿⣽⢣⠄⠀⠀⠀⠀⠀⠀⠀⠀
     ⠀⠀⢀⣰⣿⣿⣿⣿⣿⣿⢿⡿⣽⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣯⠻⡜⠀⠀⠀⠀⠀⠀⠀⠀
      ⣀⡶⣯⢻⢏⣿⣳⣾⣿⣿⣟⣿⣿⣿⣿⣿⣿⣿⣿⣿⡟⣏⢷⣺⣝⡾⣤⣄⡀⠀⠀⠀⠀⠀
⠀⠀⠀⢀⣴⣾⣿⣿⣿⣯⣾⣷⣿⣿⣿⣿⣿⢿⡻⣟⡻⢿⢿⡛⣝⣶⣿⣾⣿⣻⢿⣿⣿⣾⡝⡶⣄⠀⠀⠀
⠀⢀⡾⣿⣿⣿⢿⣽⣿⣿⣿⣿⣿⣿⣿⣹⣾⣯⣷⣿⣿⣞⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣾⡵⣩⠞⣀⠀
⠀⢾⣽⡿⣟⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⢿⣯⢿⣱⢻⡄⠠
⣘⢮⣟⣿⣽⡿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⣟⢯⢞⣯⢖⠠
⠌⠹⣾⡹⣞⣿⣽⣻⣿⣿⣿⣿⣿⣿⣿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢿⡽⣟⣯⢟⣾⡹⢎⡗⣎⠒
⠐⣌⢳⡽⣹⡞⣧⡟⣷⢫⣷⣻⢭⣻⡽⢯⣿⢿⣯⢿⣿⣟⡿⣿⢯⣟⣿⢾⡽⢯⢿⣹⢮⡟⡾⢡⣋⠼⡐⠂
⠀⠈⢆⠳⡌⣝⢲⡝⣯⢳⡞⣵⢫⣗⡻⣝⢾⣛⢾⣛⠾⣭⢻⡭⣟⢾⡹⢧⠻⣝⢮⡝⡾⠼⣍⠳⢌⠂⠁⠀
⠀⠀⠀⠑⠘⢌⠣⠞⡰⢣⠚⡌⢧⡘⡱⢊⠞⢨⠓⡌⠳⣡⠓⡜⢌⠲⡉⢏⡱⢎⠲⡙⡜⠣⠎⠑⠀⠀⠀⠀
` + Yellow + `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀  .                  ` + Reset + `
` + Yellow + `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣴⡏                   ` + Reset + `
` + Yellow + `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣾⡟                   ⠀` + Reset + `
` + Yellow + `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣿⡟⠀                   ⠀` + Reset + `
` + Yellow + `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣼⣿⡟⣼⣿⡟⡟⠋                 ` + Reset + `
` + Yellow + `                  ⠀⣼⣿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀  ⠀` + Reset + `
` + Yellow + `⠀                 ⣼⠟⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀  ⠀` + Reset + `
` + Yellow + `                 ⠴⠏                     ` + Reset + `
` + Yellow + `                ˙                       ` + Reset + `
`

const sunny string = `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
` + Yellow + `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡔⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀` + Reset + `⠀
` + Yellow + `⠀⠀⠀⠀⠀⠀⣀⠀⠀⠀⡼⠀⢰⣟⡇⠀⠀⢤⠂⠀⠀⠀⠀⠀⠀⠀` + Reset + `⠀
` + Yellow + `⠀⠀⠀⠀⠀⠀⣿⣦⣀⠀⢛⡵⣸⣟⡾⡆⡸⡹⠀⣴⡿⠟⠁⠀⠀⠀⠀` + Reset + `
` + Yellow + `⠀⠀⠀⣧⠀⢹⣾⡝⣿⢺⡖⣽⣯⣽⡏⣿⣥⣾⢳⣿⠀⠀⠀⠀⠀⠀⠀` + Reset + `
` + Yellow + `⡀⠀⠀⠈⠹⣢⠶⢿⠝⡫⣕⡱⡸⢔⡹⡌⣟⡹⣟⡮⡰⡍⠋⠁⠀⠀⠀` + Reset + `
` + Yellow + `⢻⡶⣶⢶⣖⡾⡯⢎⠵⡓⣤⠳⣱⠩⢖⡩⢦⠱⢎⣵⣯⣴⣾⡻⠟⠶⠄` + Reset + `
` + Yellow + `⠀⠉⠛⠚⠾⣿⢱⡩⢎⡵⣊⠗⣌⢳⡩⣝⢮⠹⡜⠬⣷⣟⠗⠁⠀⠀⠀` + Reset + `
` + Yellow + `⠈⠲⠬⠓⢮⣼⢣⡜⢥⡚⢬⡙⣬⠲⣑⠮⣌⠳⡜⢣⣯⣝⡬⠣⠖⠄⠀` + Reset + `
` + Yellow + `⠀⠀⣠⣾⣻⣿⡡⢞⢢⢝⢢⡝⢤⢳⢩⠖⣩⠎⡵⢣⡿⣵⢶⡦⣄⠀⠀` + Reset + `
` + Yellow + `⠲⢟⣟⠾⠓⡛⡟⣬⠣⣎⠷⣘⠇⣮⠱⣎⡱⢎⣱⣛⠻⠾⠋⠙⠫⢧⠀` + Reset + `
` + Yellow + `⠀⠀⢀⣄⡲⠝⣸⢯⡷⣌⡚⡥⢫⡕⢫⡔⣥⢿⡽⡊⠽⢱⠄⠀⠀⠀⠀` + Reset + `
` + Yellow + `⠀⠀⠀⠀⠀⠀⡿⣟⡻⣣⢛⣿⣧⡼⣟⣸⡛⣿⡻⣿⢄⠀⠣⠀⠀⠀⠀` + Reset + `
` + Yellow + `⠀⠀⠀⠀⢀⣴⣻⠝⠀⡷⠏⠹⡾⣽⡇⢺⣣⠈⠉⠻⣯⠀⠀⠀⠀⠀⠀` + Reset + `
` + Yellow + `⠀⠀⠀⠀⠍⠉⠀⠀⠔⠋⠀⠀⢹⣳⠇⠀⠍⠀⠀⠀⠍⠀⠀⠀⠀⠀⠀` + Reset + `
` + Yellow + `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠜⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀` + Reset + `
	`

const partialCloud string = `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
` + Yellow + `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡔⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀` + Reset + `⠀
` + Yellow + `⠀⠀⠀⠀⠀⠀⣀⠀⠀⠀⡼⠀⢰⣟⡇⠀⠀⢤⠂⠀⠀⠀⠀⠀⠀⠀⠀             ` + Reset + `
` + Yellow + `⠀⠀⠀⠀⠀⣿⣦⣀⠀⢛⡵⣸⣟⡾⡆⡸⡹⠀⣴⡿⠟⠁⠀⠀⠀⠀              ` + Reset + `
` + Yellow + `⠀⠀⠀⣧⠀⢹⣾⡝⣿⢺⡖⣽⣯⣽⡏⣿⣥⣾⢳⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀    ` + Reset + `
` + Yellow + `⡀⠀⠀⠈⠹⣢⠶⢿⠝⡫⣕⡱⡸⢔⡹⡌⣟⡹⣟⡮⡰⡍⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀    ` + Reset + `
` + Yellow + `⢻⡶⣶⢶⣖⡾⡯⢎⠵⡓⣤⠳⣱⠩⢖⡩⢦⠱⢎⣵⣯⣴⣾⡻⠟⠶⠄⠀⠀⠀⠀⠀⠀⠀⠀⠀  ⠀ ` + Reset + `
` + Yellow + `⠀⠀⠉⠛⠚⠾⣿⢱⡩⢎⡵⣊⠗⣌⢳⡩⣝⢮` + Reset + `⣴⣶⣶⡿⣽⣻⣞⣶⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀  ⠀
` + Yellow + `⠀⠈⠲⠬⠓⢮⣼⢣⡜⢥⡚⢬⡙⣬` + Reset + `⣴⣿⣿⣿⣿⣿⣿⢿⣿⡿⣿⣿⣾⣯⢷⣂⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
` + Yellow + `⠀⠀⠀⣠⣾⣻⣿⡡⢞⢢` + Reset + `⣤⣴⣾⣿⣿⣟⣾⣿⡿⣟⣿⣿⣾⣿⣟⣿⣿⣿⣿⣽⢣⠄⠀⠀⠀⠀⠀⠀⠀⠀
` + Yellow + `⠀⠲⢟⣟⠾⠓⡛⢀` + Reset + `⣰⣿⣿⣿⣿⣿⣿⢿⡿⣽⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣯⠻⡜⠀⠀⠀⠀⠀⠀⠀⠀
` + Yellow + `⠀⠀⠀⢀⣄⡲` + Reset + `⣀⡶⣯⢻⢏⣿⣳⣾⣿⣿⣟⣿⣿⣿⣿⣿⣿⣿⣿⣿⡟⣏⢷⣺⣝⡾⣤⣄⡀⠀⠀⠀⠀⠀
⠀⠀⠀⢀⣴⣾⣿⣿⣿⣯⣾⣷⣿⣿⣿⣿⣿⢿⡻⣟⡻⢿⢿⡛⣝⣶⣿⣾⣿⣻⢿⣿⣿⣾⡝⡶⣄⠀⠀⠀
⠀⢀⡾⣿⣿⣿⢿⣽⣿⣿⣿⣿⣿⣿⣿⣹⣾⣯⣷⣿⣿⣞⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣾⡵⣩⠞⣀⠀
⠀⢾⣽⡿⣟⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⢿⣯⢿⣱⢻⡄⠠
⣘⢮⣟⣿⣽⡿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣯⣟⢯⢞⣯⢖⠠
⠌⠹⣾⡹⣞⣿⣽⣻⣿⣿⣿⣿⣿⣿⣿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢿⡽⣟⣯⢟⣾⡹⢎⡗⣎⠒
⠐⣌⢳⡽⣹⡞⣧⡟⣷⢫⣷⣻⢭⣻⡽⢯⣿⢿⣯⢿⣿⣟⡿⣿⢯⣟⣿⢾⡽⢯⢿⣹⢮⡟⡾⢡⣋⠼⡐⠂
⠀⠈⢆⠳⡌⣝⢲⡝⣯⢳⡞⣵⢫⣗⡻⣝⢾⣛⢾⣛⠾⣭⢻⡭⣟⢾⡹⢧⠻⣝⢮⡝⡾⠼⣍⠳⢌⠂⠁⠀
⠀⠀⠀⠑⠘⢌⠣⠞⡰⢣⠚⡌⢧⡘⡱⢊⠞⢨⠓⡌⠳⣡⠓⡜⢌⠲⡉⢏⡱⢎⠲⡙⡜⠣⠎⠑⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠁⠈⠈⠀⠀⠀⠁⠈⠀⠀⠀⠁⠀⠀⠀⠈⠀⠁⠈⠀⠈⠁⠈⠀⠁⠀⠀⠀⠀⠀⠀
	`

// Gets the weather map from the API given an input address/location
func getWeather(address string) map[string]any {
	baseURL := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timelinepreview/%s?key=-&options=preview", address)

	resp, err := http.Get(baseURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll((resp.Body))
	if err != nil {
		panic(err)
	}
	stringBody := string(body)

	var weatherMap map[string]any
	err = json.Unmarshal([]byte(stringBody), &weatherMap)
	if err != nil {
		fmt.Println(stringBody)
		panic(err)
	}

	return weatherMap
}

// Receives user input for an address or location
func inputAddress() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter a location or address: ")
	inputAddress, _ := reader.ReadString('\n')
	address := strings.TrimSpace(inputAddress)

	return address
}

// Gets the address. Pulls from ~/.weather_config if it exists, otherwise requests it from the user
func getAddress(configFile string) string {
	var address string

	// Check if config file exists. If so, use address in config file. Else, request from user
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		address = inputAddress()
	} else if err != nil {
		fmt.Printf("An error occurred while checking file '%s': %v\n", configFile, err)
		panic(err)
	} else {
		content, err := os.ReadFile(configFile)
		if err != nil {
			panic(err)
		}

		configLines := strings.Split(string(content), "\n")
		foundAddress := false
		for _, line := range configLines {
			if strings.Contains(line, "=") {
				key := strings.TrimSpace(strings.Split(line, "=")[0])
				value := strings.TrimSpace(strings.Split(line, "=")[1])

				if key == "address" {
					address = value
					foundAddress = true
					break
				}
			}
		}

		// If no address found in config, request it from user
		if !foundAddress {
			address = inputAddress()
		}

	}

	return address
}

// Returns location of config file
func getConfigLocation() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return homeDir + "/.weather_config"
}

// Returns the correct condition icon depending on condition type and severity
func getIcon(condition string, severity float64) string {
	if severity >= 4 {
		return storm
	}
	switch condition {
	case "Clear":
		return sunny
	case "Snow":
		return snow
	case "Overcast":
		return cloud
	case "Cloudy":
		return cloud
	case "Partially cloudy":
		return partialCloud
	case "Rain":
		return rain
	default:
		return sunny
	}
}

// Converts military time to regular time
func convertTime(time any) any {
	stringTime := fmt.Sprintf("%v", time)
	timeParts := strings.Split(stringTime, ":")

	intHour, _ := strconv.Atoi(timeParts[0])
	if 12 > intHour {
		return fmt.Sprintf("%d:%s AM", intHour, timeParts[1])
	} else {
		intHour = intHour - 12
		return fmt.Sprintf("%d:%s PM", intHour, timeParts[1])
	}
}

// Prints the weather details with an icon and table of values
func printWeather(weatherMap map[string]any, address string) {
	currentConditionsInterface := weatherMap["currentConditions"]

	// Temp range and description aren't included in currentConditions, so grabs today from days array and gets them from there
	var description any
	var tempHigh any
	var tempLow any
	daysInterface := weatherMap["days"]
	if days, ok := daysInterface.([]any); ok {
		if today, ok := days[0].(map[string]any); ok {
			description = today["description"]
			tempHigh = today["tempmax"]
			tempLow = today["tempmin"]
		}
	}

	if currentConditions, ok := currentConditionsInterface.(map[string]any); ok {
		// Stores needed values into variables from weatherMap
		temp := currentConditions["temp"]
		precip := currentConditions["precip"]
		precipprob := currentConditions["precipprob"]
		preciptype := currentConditions["preciptype"]
		humidity := currentConditions["humidity"]
		conditions := currentConditions["conditions"]
		uvindex := currentConditions["uvindex"]
		windspeed := currentConditions["windspeed"]
		sunrise := currentConditions["sunrise"]
		sunset := currentConditions["sunset"]
		severity := currentConditions["severerisk"]

		// Converts severity interface into a float
		var severityVal float64
		if severityConvert, ok := severity.(float64); ok {
			severityVal = severityConvert
		} else {
			severityVal = 0
		}

		// Creates arrays that the table will be created from. Excludes precipitation if there is none, otherwise inserts it
		statsHeaders := []any{"", "Description", "Condition", "Temperature", "Precipitation Chance", "Humidity", "UV Index", "Wind Speed", "Sunrise", "Sunset", "Storm Risk", ""}
		statsList := []any{"┌────────────── " + Red + address + Reset + " " + extendHeader(description, address) + "──────────────┐", description, conditions, temp, precipprob, humidity, uvindex, windspeed, convertTime(sunrise), convertTime(sunset), severityVal, "└" + strings.Repeat("─", len(address)) + "──────────────────────────────┘"}
		if preciptype != nil {
			statsHeaders = slices.Insert(statsHeaders, 5, "Precipitation Type")
			statsList = slices.Insert(statsList, 5, preciptype)
			statsHeaders = slices.Insert(statsHeaders, 6, "Precipitation Depth")
			statsList = slices.Insert(statsList, 6, precip)
		}

		// Determines whether or not the stats list is odd or even. Needed to provide the correct offset compared to the icon
		var oddOrEven int
		if len(statsHeaders)%2 == 0 {
			oddOrEven = 1
		} else {
			oddOrEven = 0
		}

		firstCondition := strings.Split(fmt.Sprintf("%v", conditions), ", ")[0]
		iconSplit := strings.Split(getIcon(firstCondition, severityVal), "\n")

		// Prints main icon along with stats list with the correct offset and NerdFont icons
		offsetLeft := (len(iconSplit) / 2) - (len(statsList) / 2)
		offsetRight := (len(iconSplit) / 2) + (len(statsList) / 2)
		for i, iconLine := range iconSplit {
			fmt.Print(iconLine)
			if i >= offsetLeft && i <= offsetRight-oddOrEven {
				arrayOffset := i - offsetLeft
				stringPadding := "             "

				switch statsHeaders[arrayOffset] {
				case "":
					fmt.Print(stringPadding)
					fmt.Printf("%v\n", statsList[arrayOffset])
				case "Description":
					stringPadding += "│ "
					fmt.Printf("%s\uf405 "+Blue+"%v: "+Reset, stringPadding, statsHeaders[arrayOffset])
					fmt.Printf("%v\n", statsList[arrayOffset])
				case "Condition":
					stringPadding += "│ "
					fmt.Printf("%s\uebaa "+Blue+"%v: "+Reset, stringPadding, statsHeaders[arrayOffset])
					fmt.Printf("%v\n", statsList[arrayOffset])
				case "Temperature":
					stringPadding += "│ "
					fmt.Printf("%s\uf2c8 "+Blue+"%v: "+Reset, stringPadding, statsHeaders[arrayOffset])
					fmt.Printf(colorTemp(statsList[arrayOffset])+"%vºF "+Reset+"(%v - %v)\n", statsList[arrayOffset], tempLow, tempHigh)
				case "Precipitation Chance":
					stringPadding += "│ "
					fmt.Printf("%s\ue275 "+Blue+"%v: "+Reset, stringPadding, statsHeaders[arrayOffset])
					fmt.Printf("%v%%\n", statsList[arrayOffset])
				case "Precipitation Type":
					stringPadding += "│ "
					fmt.Printf("%s\ue318 "+Blue+"%v: "+Reset, stringPadding, statsHeaders[arrayOffset])
					fmt.Printf("%s\n", interfaceToArray(statsList[arrayOffset]))
				case "Precipitation Depth":
					stringPadding += "│ "
					fmt.Printf("%s\uef30 "+Blue+"%v: "+Reset, stringPadding, statsHeaders[arrayOffset])
					fmt.Printf("%v in.\n", statsList[arrayOffset])
				case "Humidity":
					stringPadding += "│ "
					fmt.Printf("%s\ue373 "+Blue+"%v: "+Reset, stringPadding, statsHeaders[arrayOffset])
					fmt.Printf("%v%%\n", statsList[arrayOffset])
				case "UV Index":
					stringPadding += "│ "
					fmt.Printf("%s\ue30d "+Blue+"%v: "+Reset, stringPadding, statsHeaders[arrayOffset])
					fmt.Printf("%v\n", statsList[arrayOffset])
				case "Wind Speed":
					stringPadding += "│ "
					fmt.Printf("%s\uef16 "+Blue+"%v: "+Reset, stringPadding, statsHeaders[arrayOffset])
					fmt.Printf("%v mph\n", statsList[arrayOffset])
				case "Sunrise":
					stringPadding += "│ "
					fmt.Printf("%s\ue34c "+Blue+"%v: "+Reset, stringPadding, statsHeaders[arrayOffset])
					fmt.Printf("%v\n", statsList[arrayOffset])
				case "Sunset":
					stringPadding += "│ "
					fmt.Printf("%s\ue34d "+Blue+"%v: "+Reset, stringPadding, statsHeaders[arrayOffset])
					fmt.Printf("%v\n", statsList[arrayOffset])
				case "Storm Risk":
					stringPadding += "│ "
					fmt.Printf("%s\ue315 "+Blue+"%v: "+Reset, stringPadding, statsHeaders[arrayOffset])
					fmt.Printf(colorSeverity(severityVal)+"%.0f"+Reset+"\n", statsList[arrayOffset])
				default:
					stringPadding += "│ "
					fmt.Printf("%s"+Blue+"%v: "+Reset, stringPadding, statsHeaders[arrayOffset])
					fmt.Printf("%v\n", statsList[arrayOffset])
				}

			} else {
				fmt.Println()
			}
		}
	}
}

// Extends the top header bar if the description is longer than it
func extendHeader(description any, address string) string {
	descriptionString := fmt.Sprintf("%v", description)
	lengthDescription := len(descriptionString) + 17
	lengthHeader := len(address) + 32
	if lengthDescription >= lengthHeader {
		return strings.Repeat("─", lengthDescription-lengthHeader+1)
	} else {
		return ""
	}
}

// Colors the temperature reading depending on heat levels
func colorTemp(temp any) string {
	if tempVal, ok := temp.(float64); ok {
		if tempVal < 50 {
			return Cyan
		} else if tempVal < 72 {
			return Green
		} else if tempVal < 85 {
			return Yellow
		} else {
			return Red
		}
	} else {
		return ""
	}
}

// Colors the severity score depending on severity level
func colorSeverity(severity float64) string {
	if severity < 4 {
		return Green
	} else if severity < 7 {
		return Yellow
	} else {
		return Red
	}
}

// Converts an interface to a joined string for the precipitation list
func interfaceToArray(i any) string {
	var s []string
	if arr, ok := i.([]any); ok {
		for _, j := range arr {
			if j, ok2 := j.(string); ok2 {
				s = append(s, capFirst(j))
			}
		}
	}
	return strings.Join(s, ", ")
}

// Capitalizes the first letter of a string
func capFirst(s string) string {
	if s == "" {
		return ""
	}

	r, _ := utf8.DecodeRuneInString(s)
	upperChar := unicode.ToUpper(r)
	return string(upperChar) + s[1:]
}

// Driver for the program
func main() {
	configFile := getConfigLocation()
	address := getAddress(configFile)

	weatherMap := getWeather(address)
	printWeather(weatherMap, address)
}
