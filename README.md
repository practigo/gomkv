# Go MKV

Try to parse [MKV]((https://en.wikipedia.org/wiki/Matroska)) for learning :)

(Code is similar to https://github.com/practigo/gomp4)

## Intro

Matroska/MKV is a Document Type of [EBML](https://en.wikipedia.org/wiki/Extensible_Binary_Meta_Language).
- EBML [spec](https://matroska-org.github.io/libebml/specs.html)
- [WebM](https://en.wikipedia.org/wiki/WebM) is a limited subset so it is supported too.

### Layout

```text
+-------------+
| EBML Header |
+---------------------------+
| Segment     | SeekHead    | index of Top-Level Elements locations (RECOMMENDED)
|             |-------------|
|             | Info        | vital information for identifying the whole Segment
|             |-------------|
|             | Tracks      | technical details for each track (decode the data)
|             |-------------|
|             | Chapters    | lists all of the chapters (points to jump)
|             |-------------|
|             | Cluster     | content for each track (SHOULD contain at least one)
|             |-------------|
|             | Cues        | temporal index for seeking (SHOULD contain at least one)
|             |-------------|
|             | Attachments | for attaching files (pictures/fonts...)
|             |-------------|
|             | Tags        | metadata for Segment/Tracks/Chapters...
+---------------------------+
```

See https://www.matroska.org/technical/diagram.html for more details.

## Refs

- https://www.matroska.org/technical/basics.html
- https://github.com/remko/go-mkvparse
