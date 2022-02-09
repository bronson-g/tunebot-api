# tunebot api
The API that handles manages users and their playlists for the [tunebot app](TBD) and [tunebot controller](https://github.com/cbeimers113/tunebot-controller).

## Run Locally (port 8080)
```bash
    data/startup_windows.sh
    cd src/
    go run main.go
```

## API Endpoints

### POST /api/register/
#### request body
```json
    {
        "username": "my_username",
        "password": "my_encrypted_password"
    }
```

#### success response body
```json
    {
        "id": "5927c9ff-57fe-4d26-a840-88ed21451feb",
        "username": "my_username",
        "playlists": [],
        "blacklist": {
            "id": "",
            "enabled": false,
            "songs": null
        }
    }
```

#### error response body
```json
    {
        "error": "error message"
    }
```

### POST /api/login/
#### request body
```json
    {
        "username": "my_username",
        "password": "my_encrypted_password"
    }
```

#### success response body
```json
    {
        "id": "5927c9ff-57fe-4d26-a840-88ed21451feb",
        "username": "my_username",
        "playlists": [],
        "blacklist": {
            "id": "",
            "enabled": false,
            "songs": null
        }
    }
```

#### error response body
```json
    {
        "error": "error message"
    }
```