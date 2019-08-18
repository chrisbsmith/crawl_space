# Crawl Space Temps

I have a problem.  My wife and I collect wines from our trips abroad and bring them home to cellar.  Rather than buy a bigger wine cooler, I wanted to know if I could use my crawl space as a cellar.  

Rather than just stick an old school thermometer under my house and crawl under to read the temps, I decided to leverage a few raspberry pis I had laying around.  

I purchased a few [temperature probes](https://www.amazon.com/Waterproof-Temperature-Stainless-Thermistor-Transimitter/dp/B07G7R3KZC/ref=asc_df_B07G7R3KZC/?tag=hyprod-20&linkCode=df0&hvadid=343868117984&hvpos=1o1&hvnetw=g&hvrand=10433247359245742323&hvpone=&hvptwo=&hvqmt=&hvdev=c&hvdvcmdl=&hvlocint=&hvlocphy=9008183&hvtargid=pla-766035023649&psc=1&tag=&ref=&adgrpid=68767445426&hvpone=&hvptwo=&hvadid=343868117984&hvpos=1o1&hvnetw=g&hvrand=10433247359245742323&hvqmt=&hvdev=c&hvdvcmdl=&hvlocint=&hvlocphy=9008183&hvtargid=pla-766035023649) that are compatable with the RPIs.  Then I got to work on my first Golang program.

Build I also wanted to know what the outside temperatures were to compare what I was seeing in my crawl space.

## Dependencies
### OpenWeather API Token

You need an API token from [Open Weather](https://openweathermap.org)

### InfluxDB

Data is written to an influx DB.  I run my on the same RPI but you could ship the data anywhere.  You'll set a username and password that are needed by this app

## Settings

These settings are stored in settings.toml file, located in the root or at `/app`

The file should look like this, but with your data

```
# config file
MyDB     = "temperatures"
username = "user"
password = "password"
OwmAPIKey = "yourapitoken"
dbPort = "http://localhost:8086"
debug = false
```

These settings can all be overridden with environment variables, prefixed with `CS_`.  For instance, to change your username, you would use `CS_USERNAME`

## Building

### From a Mac
```bash
env GOOS=linux GOARCH=arm GOARM=5 go build temperatures.go
```

## Results
Turns out I can use my crawl space for a wine cellar.  Throughout the year, temperatures hover around 55-60F

## Enhancements
I also wanted to know what my central HVAC was doing when I was measuring temperatures.  I found t[his application](https://github.com/peckrob/nest-watch) did exactly what I needed

I visualized all of this data with Grafana
