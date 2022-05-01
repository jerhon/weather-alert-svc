# Downloads the current shapefiles from NWS
# These are very big files, ~50mb in size

# Counties
# Sourced from https://www.weather.gov/gis/Counties
Invoke-WebRequest -Uri https://www.weather.gov/source/gis/Shapefiles/County/c_22mr22.zip -OutFile c_22mr22.zip

# Marine Zones
# Sourced from https://www.weather.gov/gis/MarineZones
# Costal
Invoke-WebRequest -Uri https://www.weather.gov/source/gis/Shapefiles/WSOM/mz22mr22.zip -OutFile mz22mr22.zip
# Offshore
Invoke-WebRequest -Uri https://www.weather.gov/source/gis/Shapefiles/WSOM/oz22mr22.zip -OutFile oz22mr22.zip
# High Seas
Invoke-WebRequest -Uri https://www.weather.gov/source/gis/Shapefiles/WSOM/hz30jn17.zip -OutFile hz30jn17.zip