# video-encoder-golang


```
go clean -testcache
```
## Configurations
Create a ".env" file in the project root and set theses env valiables:
```
LOCAL_STORAGE_PATH="<path_to_save_downloaded_files_from_CGP>"
```
```
GOOGLE_APPLICATION_CREDENTIALS="<path_to_credentials_file_from_CGP>"
```
## Running tests
> "video_service_test" needs a GCP account with a bucket storage configured

```
go test -v ./...
```
