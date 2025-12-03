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
	"time"

	"github.com/go-viper/mapstructure/v2"
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
   ` + White + `         ❄️      .         ❄            ` + Reset + `
   ` + White + `             .                          ` + Reset + `
   ` + White + `      ❄️   .         .  ❄               ` + Reset + `
   ` + White + `               ❄️ .          .          ` + Reset + `
   ` + White + `       .  .             .               ` + Reset + `
   ` + White + `    .           .      ❄️               ` + Reset + `
   ` + White + `      ❄️    .                           ` + Reset + `
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
          ` + Yellow + `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡔⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀      ⠀⠀` + Reset + `⠀
          ` + Yellow + `⠀⠀⠀⠀⠀⠀⣀⠀⠀⠀⡼⠀⢰⣟⡇⠀⠀⢤⠂⠀⠀⠀⠀⠀⠀⠀      ` + Reset + `⠀
          ` + Yellow + `⠀⠀⠀⠀⠀⠀⣿⣦⣀⠀⢛⡵⣸⣟⡾⡆⡸⡹⠀⣴⡿⠟⠁⠀⠀⠀⠀      ` + Reset + `
          ` + Yellow + `⠀⠀⠀⣧⠀⢹⣾⡝⣿⢺⡖⣽⣯⣽⡏⣿⣥⣾⢳⣿⠀⠀⠀⠀⠀⠀⠀      ` + Reset + `
          ` + Yellow + `⡀⠀⠀⠈⠹⣢⠶⢿⠝⡫⣕⡱⡸⢔⡹⡌⣟⡹⣟⡮⡰⡍⠋⠁⠀⠀⠀      ` + Reset + `
          ` + Yellow + `⢻⡶⣶⢶⣖⡾⡯⢎⠵⡓⣤⠳⣱⠩⢖⡩⢦⠱⢎⣵⣯⣴⣾⡻⠟⠶⠄      ` + Reset + `
          ` + Yellow + `⠀⠉⠛⠚⠾⣿⢱⡩⢎⡵⣊⠗⣌⢳⡩⣝⢮⠹⡜⠬⣷⣟⠗⠁⠀⠀⠀      ` + Reset + `
          ` + Yellow + `⠈⠲⠬⠓⢮⣼⢣⡜⢥⡚⢬⡙⣬⠲⣑⠮⣌⠳⡜⢣⣯⣝⡬⠣⠖⠄⠀      ` + Reset + `
          ` + Yellow + `⠀⠀⣠⣾⣻⣿⡡⢞⢢⢝⢢⡝⢤⢳⢩⠖⣩⠎⡵⢣⡿⣵⢶⡦⣄⠀⠀      ` + Reset + `
          ` + Yellow + `⠲⢟⣟⠾⠓⡛⡟⣬⠣⣎⠷⣘⠇⣮⠱⣎⡱⢎⣱⣛⠻⠾⠋⠙⠫⢧⠀      ` + Reset + `
          ` + Yellow + `⠀⠀⢀⣄⡲⠝⣸⢯⡷⣌⡚⡥⢫⡕⢫⡔⣥⢿⡽⡊⠽⢱⠄⠀⠀⠀⠀      ` + Reset + `
          ` + Yellow + `⠀⠀⠀⠀⠀⠀⡿⣟⡻⣣⢛⣿⣧⡼⣟⣸⡛⣿⡻⣿⢄⠀⠣⠀⠀⠀⠀      ` + Reset + `
          ` + Yellow + `⠀⠀⠀⠀⢀⣴⣻⠝⠀⡷⠏⠹⡾⣽⡇⢺⣣⠈⠉⠻⣯⠀⠀⠀⠀⠀⠀      ` + Reset + `
          ` + Yellow + `⠀⠀⠀⠀⠍⠉⠀⠀⠔⠋⠀⠀⢹⣳⠇⠀⠍⠀⠀⠀⠍⠀⠀⠀⠀⠀⠀      ` + Reset + `
          ` + Yellow + `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠜⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀      ` + Reset + `
                                           `

const moon string = `
          ⠀⠀⠀⠀⠀⠀⠀⣀⣤⡴⠶⠿⠛⢏⡿⠖⠂⠀⠀⠀⠀⠀⠀⠀⠀⠀      ⠀
          ⠀⠀⠀⠀⢀⡴⣞⠯⠉⠈⠀⣠⡶⠊⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀      ⠀
          ⠀⠀⣠⣴⠟⠙⠈⣎⡹⠂⣴⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀      ⠀
          ⠀⢠⣯⡟⢚⣀⠀⠀⠀⡰⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀      ⠀
          ⠀⡿⡃⢋⣌⠂⠈⠆⡠⡎⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀      ⠀
          ⣸⢿⠎⢰⡈⠀⠈⠀⣹⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀      ⠀
          ⣸⣿⡄⠷⣠⠀⠀⠀⡸⡗⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀      ⠀
          ⢹⣾⣿⣷⣛⠀⠀⠀⣜⣷⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀      ⠀
          ⠸⣿⢻⣏⠸⡃⠀⠀⠈⢹⢧⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀      ⠀
          ⠀⢻⣿⡮⣠⣗⠒⠤⠀⠀⠹⣳⣤⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀      ⠀
          ⠀⠀⣼⢿⣗⣧⣦⢤⠇⢀⣤⣄⠙⣿⡦⣄⠀⠀⠀⠀⠀⠀⠀⠀⢀⣤⠄      
          ⠀⠀⠈⠛⢏⣷⣶⣜⣤⣈⠂⠜⠀⢀⣀⡉⡭⠯⠖⠲⠒⢶⢖⣯⠟⠁⠀      
          ⠀⠀⠀⠀⠀⠙⠻⣿⣿⣷⣷⣯⣔⡿⣃⠦⡵⣠⠠⢤⣤⠿⠋⠀⠀⠀⠀      
          ⠀⠀⠀⠀⠀⠀⠀⠀⠉⠛⠓⠿⠽⣷⣿⣾⡿⠞⠛⠉⠀⠀⠀⠀⠀⠀⠀      
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

const partialCloudNight string = `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
   ⠀⠀⠀⠀⠀⠀⠀⠀⣀⣤⡴⠶⠿⠛⢏⡿⠖⠂⠀⠀⠀⠀⠀⠀⠀⠀              
   ⠀⠀⠀⠀⠀⢀⡴⣞⠯⠉⠈⠀⣠⡶⠊⠀⠀⠀⠀⠀⠀⠀                  
   ⠀⠀⠀⣠⣴⠟⠙⠈⣎⡹⠂⣴⠋⠀⠀⠀⠀                       
   ⠀⠀⢠⣯⡟⢚⣀⠀⠀⠀⡰⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀            
   ⡀⠀⡿⡃⢋⣌⠂⠈⠆⡠⡎⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀                 
    ⣸⢿⠎⢰⡈⠀⠈⠀⣹⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀  ⠀                 
   ⠀⣸⣿⡄⠷⣠⠀⠀⠀⡸⡗        ⢀⣤⣶⡿⣽⣻⣞⣶⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀  ⠀
   ⠀⢹⣾⣿⣷⣛⠀⠀⠀⣜⣷    ⡀⢀⣴⣾⣿⣿⣿⢿⣿⡿⣿⣿⣾⣯⢷⣂⠀⠀⠀⠀⠀   ⠀
   ⠀⠸⣿⢻⣏⠸⡃⠀⠀⠈⢹⣤⣴⣾⣿⣿⣟⣾⣿⡿⣟⣿⣿⣾⣿⣟⣿⣿⣿⣿⣽⢣⠄⠀⠀⠀⠀  ⠀
   ⠀ ⠲⢟⣟⠾⠓⡛⢀⣰⣿⣿⣿⣿⣿⣿⢿⡿⣽⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣯⠻⡜⠀⠀⠀⠀ ⠀⠀
   ⠀⠀ ⣼⢿⣗⣀⡶⣯⢻⢏⣿⣳⣾⣿⣿⣟⣿⣿⣿⣿⣿⣿⣿⣿⣿⡟⣏⢷⣺⣝⡾⣤⣄⡀⠀⠀⠀⠀⠀
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

// Weather structs that store data received from the weather API
type WeatherConditions struct {
	Address           string
	Alerts            []Alert
	CurrentConditions Condition
	Description       string
	Days              []Day
	ResolvedAddress   string
}

type Alert struct {
	Description string
	Event       string
	Headline    string
}

type Condition struct {
	Conditions string
	DateTime   string
	Humidity   float64
	Precip     float64
	PrecipProb int
	PrecipType []string
	Sunrise    string
	Sunset     string
	Temp       float64
	UVIndex    int
	SevereRisk float64
	WindSpeed  float64
}

type Day struct {
	Description string
	TempMax     float64
	TempMin     float64
}

// Gets the weather data from the API given an input address/location
func getWeather(address string) WeatherConditions {
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

	// Converts the received JSON data and converts it into a map
	var weatherData map[string]any
	err = json.Unmarshal([]byte(stringBody), &weatherData)
	if err != nil {
		fmt.Println(stringBody)
		panic(err)
	}

	// Loads the JSON data into the weather structs using mapstructure
	var weatherMap WeatherConditions
	err = mapstructure.Decode(weatherData, &weatherMap)
	if err != nil {
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
func getIcon(condition string, severity float64, currentTime string, sunrise string, sunset string) string {
	if severity >= 4 && condition == "Rain" {
		return storm
	}

	ct, _ := time.Parse("15:04:05", currentTime)
	sr, _ := time.Parse("15:04:05", sunrise)
	ss, _ := time.Parse("15:04:05", sunset)
	switch condition {
	case "Clear":
		if ct.Before(sr) || ct.After(ss) {
			return moon
		} else {
			return sunny
		}
	case "Snow":
		return snow
	case "Overcast":
		return cloud
	case "Cloudy":
		return cloud
	case "Partially cloudy":
		if ct.Before(sr) || ct.After(ss) {
			return partialCloudNight
		} else {
			return partialCloud
		}
	case "Rain":
		return rain
	default:
		return sunny
	}
}

// Converts military time to regular time
func convertTime(time string) string {
	timeParts := strings.Split(time, ":")

	intHour, _ := strconv.Atoi(timeParts[0])
	if 12 > intHour {
		return fmt.Sprintf("%d:%s AM", intHour, timeParts[1])
	} else {
		intHour = intHour - 12
		return fmt.Sprintf("%d:%s PM", intHour, timeParts[1])
	}
}

// Prints the weather details with an icon and table of values
func printWeather(weatherMap WeatherConditions) {
	currentConditions := weatherMap.CurrentConditions

	// Collects any alerts if they exist
	var alertEvents []string
	var joinedAlerts string
	for _, alert := range weatherMap.Alerts {
		alertEvents = append(alertEvents, alert.Event)
	}
	joinedAlerts = strings.Join(alertEvents, ", ")

	// Temp range and description aren't included in currentConditions, so grabs today from days array and gets them from there
	description := weatherMap.Days[0].Description
	tempHigh := weatherMap.Days[0].TempMax
	tempLow := weatherMap.Days[0].TempMin

	// Creates array of headers that the table will be created from
	statsHeaders := []string{"Header 1", "Description", "Condition", "Temperature", "Precipitation Chance", "Humidity", "UV Index", "Wind Speed", "Sunrise", "Sunset", "Storm Risk", "Header 2"}
	topHeader := "┌────────────── " + Red + weatherMap.ResolvedAddress + Reset + " " + extendUpperHeader(description, weatherMap.ResolvedAddress) + "──────────────┐"
	bottomHeader := "└" + strings.Repeat("─", len(weatherMap.ResolvedAddress)) + "──────────────────────────────┘"

	// Adds precipitation type and depth if preciptype exists
	if currentConditions.PrecipType != nil {
		statsHeaders = slices.Insert(statsHeaders, 5, "Precipitation Type")
		statsHeaders = slices.Insert(statsHeaders, 6, "Precipitation Depth")
	}
	// Adds alerts to headers and extends the lower header if the alerts list go past the current length
	if len(alertEvents) > 0 {
		statsHeaders = slices.Insert(statsHeaders, len(statsHeaders)-1, "Alerts")
		bottomHeader = bottomHeader[:3] + extendLowerHeader(weatherMap.ResolvedAddress, joinedAlerts) + bottomHeader[3:]
	}

	// Pulls first condition from the list to determine the weather icon used
	firstCondition := strings.Split(currentConditions.Conditions, ", ")[0]
	iconSplit := strings.Split(getIcon(firstCondition, currentConditions.SevereRisk, currentConditions.DateTime, currentConditions.Sunrise, currentConditions.Sunset), "\n")

	// Prints main icon along with stats list with the correct offset and NerdFont icons
	offsetLeft := (len(iconSplit) / 2) - (len(statsHeaders) / 2)
	offsetRight := (len(iconSplit) / 2) + (len(statsHeaders) / 2)
	oddOffset := 1 - (len(statsHeaders) % 2)
	for i, iconLine := range iconSplit {
		fmt.Print(iconLine)
		if i >= offsetLeft && i <= offsetRight-oddOffset {
			arrayOffset := i - offsetLeft
			stringPadding := "             "

			switch statsHeaders[arrayOffset] {
			case "Header 1":
				fmt.Print(stringPadding)
				fmt.Printf("%s\n", topHeader)
			case "Header 2":
				fmt.Print(stringPadding)
				fmt.Printf("%s\n", bottomHeader)
			case "Description":
				stringPadding += "│ "
				fmt.Printf("%s\uf405 "+Blue+"%s: "+Reset, stringPadding, statsHeaders[arrayOffset])
				fmt.Printf("%s\n", description)
			case "Condition":
				stringPadding += "│ "
				fmt.Printf("%s\uebaa "+Blue+"%s: "+Reset, stringPadding, statsHeaders[arrayOffset])
				fmt.Printf("%s\n", currentConditions.Conditions)
			case "Temperature":
				stringPadding += "│ "
				fmt.Printf("%s\uf2c8 "+Blue+"%s: "+Reset, stringPadding, statsHeaders[arrayOffset])
				fmt.Printf(colorTemp(currentConditions.Temp)+"%0.fºF "+Reset+"(%0.f - %0.f)\n", currentConditions.Temp, tempLow, tempHigh)
			case "Precipitation Chance":
				stringPadding += "│ "
				fmt.Printf("%s\ue275 "+Blue+"%s: "+Reset, stringPadding, statsHeaders[arrayOffset])
				fmt.Printf("%d%%\n", currentConditions.PrecipProb)
			case "Precipitation Type":
				stringPadding += "│ "
				fmt.Printf("%s\ue318 "+Blue+"%s: "+Reset, stringPadding, statsHeaders[arrayOffset])
				fmt.Printf("%s\n", strings.Join(currentConditions.PrecipType, ", "))
			case "Precipitation Depth":
				stringPadding += "│ "
				fmt.Printf("%s\uef30 "+Blue+"%s: "+Reset, stringPadding, statsHeaders[arrayOffset])
				fmt.Printf("%0.f in.\n", currentConditions.Precip)
			case "Humidity":
				stringPadding += "│ "
				fmt.Printf("%s\ue373 "+Blue+"%s: "+Reset, stringPadding, statsHeaders[arrayOffset])
				fmt.Printf("%0.f%%\n", currentConditions.Humidity)
			case "UV Index":
				stringPadding += "│ "
				fmt.Printf("%s\ue30d "+Blue+"%s: "+Reset, stringPadding, statsHeaders[arrayOffset])
				fmt.Printf("%d\n", currentConditions.UVIndex)
			case "Wind Speed":
				stringPadding += "│ "
				fmt.Printf("%s\uef16 "+Blue+"%s: "+Reset, stringPadding, statsHeaders[arrayOffset])
				fmt.Printf("%0.f mph\n", currentConditions.WindSpeed)
			case "Sunrise":
				stringPadding += "│ "
				fmt.Printf("%s\ue34c "+Blue+"%s: "+Reset, stringPadding, statsHeaders[arrayOffset])
				fmt.Printf("%s\n", convertTime(currentConditions.Sunrise))
			case "Sunset":
				stringPadding += "│ "
				fmt.Printf("%s\ue34d "+Blue+"%s: "+Reset, stringPadding, statsHeaders[arrayOffset])
				fmt.Printf("%s\n", convertTime(currentConditions.Sunset))
			case "Storm Risk":
				stringPadding += "│ "
				fmt.Printf("%s\ue315 "+Blue+"%s: "+Reset, stringPadding, statsHeaders[arrayOffset])
				fmt.Printf(colorSeverity(currentConditions.SevereRisk)+"%.0f"+Reset+"\n", currentConditions.SevereRisk)
			case "Alerts":
				stringPadding += "│ "
				fmt.Printf("%s\uf421 "+Blue+"%s: "+Reset, stringPadding, statsHeaders[arrayOffset])
				fmt.Printf("%s\n", joinedAlerts)
			default:
				stringPadding += "│ "
				fmt.Printf("%s"+Blue+"%s: "+Reset, stringPadding, statsHeaders[arrayOffset])
			}

		} else {
			fmt.Println()
		}
	}
}

// Extends the upper header bar if the description is longer than it
func extendUpperHeader(description string, address string) string {
	lengthDescription := len(description) + 17
	lengthHeader := len(address) + 32
	if lengthDescription >= lengthHeader {
		return strings.Repeat("─", lengthDescription-lengthHeader+1)
	} else {
		return ""
	}
}

// Extends the lower header bar if the alerts list is longer than it
func extendLowerHeader(address string, alerts string) string {
	lengthAlerts := len(alerts) + 12
	lengthHeader := len(address) + 32
	if lengthAlerts >= lengthHeader {
		return strings.Repeat("─", lengthAlerts-lengthHeader+1)
	} else {
		return ""
	}
}

// Colors the temperature reading depending on heat levels
func colorTemp(temp float64) string {
	if temp < 45 {
		return Cyan
	} else if temp < 60 {
		return Blue
	} else if temp < 78 {
		return Green
	} else if temp < 90 {
		return Yellow
	} else {
		return Red
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

// Driver for the program
func main() {
	configFile := getConfigLocation()
	address := getAddress(configFile)

	weatherMap := getWeather(address)
	printWeather(weatherMap)
}
