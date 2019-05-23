# whatgotdone

[![CircleCI](https://circleci.com/gh/mtlynch/whatgotdone.svg?style=svg&circle-token=180495ad17cc0343547e430e81d28b66ff87e9f4)](https://circleci.com/gh/mtlynch/whatgotdone)

## To build

### Build frontend

```
cd web/frontend
npm run build
```

### Run backend

```
cmd /c "go build --tags dev -o main.exe web\main.go && main.exe"
```