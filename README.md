# kv-be-go-gin-dataapi-collections
Golang backend for KillrVideo which uses Gin and the Data API.

## Work in progress!

## Overview
This repo demonstrates modern API best-practices with:

* Restful, typed request/response models
* Role-based JWT auth
* Micro-service friendly layout

## Prerequisites
1. **Go 1.25.5** or later.
2. A **DataStax Astra DB** serverless database – [grab a free account](https://astra.datastax.com).
3. Astra DB Go Data API client

## Setup & Configuration
```bash
# clone
git clone git@github.com:KillrVideo/kv-be-go-gin-dataapi-collections.git
cd kv-be-go-gin-dataapi-collections

# build and install dependencies
go get github.com/golang-jwt/jwt/v5
go get github.com/datastax/astra-db-go
go get -u github.com/gin-gonic/gin
```

### Database collectiosn:
1. Create a new keyspace named `killrvideo_dataapi`.
2. Create the following non-vector-enabled collections:
 - `comments`
 - `content_moderation`
 - `users`
 - `ratings`
 - `video_ratings`
3. Create the following vector-enabled collection:
 - `videos` (with a 384-dimensional vector)

### Load data:
AstraDB - [KillrVideo collections loader](https://github.com/KillrVideo/killrvideo-data/blob/master/loaders/astra-collections/README.md)


### Environment variables (via `export`):
| Variable | Description |
|----------|-------------|
| `ASTRA_DB_APPLICATION_TOKEN` | The token created in the Astra UI |
| `ASTRA_DB_API_ENDPOINT` | The API endpoint for your Astra database |
| `ASTRA_DB_KEYSPACE` | `killrvideo_dataapi` |
| `JWT_KEY` | A random, 64-byte secret key used to sign the JSON Web Token |
| `YOUTUBE_API_KEY` | Required for pulling new video info from YouTube API |
| `HF_API_KEY` | HuggingFace key used to hit a HuggingFace Space to create an embedding |
