# Queries

## GetCapabilities
```bash
curl 'http://localhost:8000/ows?access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjN9.5z5GA65oR__Viieu_1it2Bjr2Ycj-DUIwDISMGfnXIQ&SERVICE=WFS&REQUEST=DescribeFeatureType&VERSION=2.0.0&TYPENAMES=breweries_denver&TYPENAME=breweries_denver'
```

## DescribeFeatureType
```bash
curl 'http://localhost:8000/ows?access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjN9.5z5GA65oR__Viieu_1it2Bjr2Ycj-DUIwDISMGfnXIQ&SERVICE=WFS&REQUEST=DescribeFeatureType&VERSION=2.0.0&TYPENAMES=breweries_denver&TYPENAME=breweries_denver'
```

## GetFeature

```bash
curl 'http://localhost:8000/ows?access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjN9.5z5GA65oR__Viieu_1it2Bjr2Ycj-DUIwDISMGfnXIQ&SERVICE=WFS&REQUEST=GetFeature&VERSION=2.0.0&TYPENAMES=breweries_denver&TYPENAME=breweries_denver' > sample_responses/GetFeature.xml
```
