# Overview

This project is currently incomplete and a work in progress.

This will be microservice and several supporting applications to synchronize and send weather alerts from NWS.
The weather alerts are obtained from the [NWS API](https://www.weather.gov/documentation/services-web-api).

While this application utilizes data from NWS, **this application and it's developer are in no way associated with the NWS.**

I make no claims about the validity of the data collected by this application nor should it's data be in a way to make any personal or official decisions based on weather events.
Use this application at your own risk.

There are a few applications to support this.

| Name                | Description                                                                  |
|---------------------|------------------------------------------------------------------------------|
| weather-alert-sync  | This will query the NWS API and add any new alerts to the mongodb database   |
| weather-zone-import | This will import shapes from the NWS shape files into the current mongodb    |
| weather-alert-svc      | TODO: This will be a microservice which returns data for imported alerts.     |
| weather-alert-publisher | TODO: This will be an application which will publish new alerts over a queue. |

Any new alerts imported to mongodb will also be sent via a message queue.
This allows other applications to use the alerts as they see fit.

## Alert Shapes

For alerts that don't provide a specific geometry, the geometries are looked up from the local mongodb.
While there is an API that allows for looking up alert geometries, I did not want to issue too many requests.

When the alert is stored in MongoDB, a query is first done to find the corresponding geometry.
When one is found, it's stored on the alert to facilitate easier geographic data.

To import shapes, they need to be downloaded from the [NWS Shape Files](https://www.weather.gov/gis/AWIPSShapefiles)
The weather-zone-import CLI program must be run.

So far, only importing the counties is supported.

## Running

```
go run .\cmd\weather-zone-import\weather-zone-import.go 
```

## Rational

I really wanted to just build a simple application in Go that I could use to further learn the language and expand my skills in k8s.
This was something that I felt was non-trivial enough for me to further understand the language.
A small project to advance my understanding of GO.

## Known Issues
* Alert import not complete, need to add logic in to only update on change and only import new alerts since previous request.
* The Marine Zones have one very large zone that contains a lot of shapes, it currently can't be imported in a single MongoDB.
