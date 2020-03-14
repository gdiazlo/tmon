package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"strconv"
	"strings"

	tm "github.com/buger/goterm"
	"github.com/fatih/color"
	"github.com/tarm/serial"
)

/*
Arduino side

Three VMA320 sensors connected to pins A0, A1 and A2
Arduino UNO board

--
#include <math.h>

char data[100];

double Thermistor(int RawADC) {
  double Temp;
  Temp = log(10000.0 / (1024.0 / RawADC - 1)); // for pull-up configuration
  Temp = 1 / (0.001129148 + (0.000234125 + (0.0000000876741 * Temp * Temp )) * Temp );
  Temp = Temp - 273.15;            // Convert Kelvin to Celcius
  return Temp;
}

void setup() {
  Serial.begin(9600);
}

void loop() {
  sprintf(data, "%d %d %d", int(Thermistor(analogRead(A0))), int(Thermistor(analogRead(A1))), int(Thermistor(analogRead(A2))));
  Serial.println(data);
  delay(1000);
}

*/

const w, h = 960, 540

func main() {

	port := flag.String("port", "COM3", "Either COMX or /dev/ttySX")
	flag.Parse()

	c := &serial.Config{Name: *port, Baud: 9600}
	serialPort, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	ch := readData(serialPort)

	tm.Output = bufio.NewWriter(color.Output)
	chart := tm.NewLineChart(100, 20)

	data := new(tm.DataTable)
	data.AddColumn("t")
	data.AddColumn("a0")
	data.AddColumn("a1")
	data.AddColumn("a2")

	i := 0.0
	for v := range ch {
		if len(v) < 3 {
			continue
		}
		tm.MoveCursor(1, 1)
		values := append([]float64{i}, v...)
		tm.Println(tm.Bold("Temperature Sensor:"), tm.Bold("[T "), values[0], tm.Bold("] [A0 "), values[1], tm.Bold("] [A1 "), values[2], tm.Bold("] [A2 "), values[3], tm.Bold("]"))
		data.AddRow(values...)
		tm.Println(chart.Draw(data))
		tm.Flush()
		i++
	}

}

func readData(port io.Reader) chan []float64 {
	ch := make(chan []float64)
	scanner := bufio.NewScanner(port)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			records := strings.Split(line, " ")
			d := make([]float64, 0, 4)

			for _, rec := range records {
				f, _ := strconv.ParseFloat(rec, 64)
				d = append(d, f)
			}
			ch <- d
		}
	}()
	return ch
}
