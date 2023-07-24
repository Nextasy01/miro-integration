# Miro Integration

## Selenium setup

In order to use Selenium. The web driver and selenium standalone server needs to be included.

There are 2 options to start selenium server:

1. Create "webdrivers" folder inside the project. Download suitable chrome driver according to your OS from [here](https://chromedriver.storage.googleapis.com/index.html?path=114.0.5735.90/) and [server standalone](https://selenium-release.storage.googleapis.com/3.5/selenium-server-standalone-3.5.3.jar). Extract chrome from chrome driver zip archive and paste it inside folder. Paste server standalone .jar archive in the same folder as well.

2. If you are familiar with Docker and have it in your system then it will be easier. You can run selenium official container, like this:

```
docker pull selenium/standalone-chrome 
docker run -d -p 4444:4444 --shm-size=2g selenium/standalone-chrome:latest
```

The port should be exactly 4444.

## Usage

If you set up selenium server using 1st option then execute the following command:
```
go run main.go -option 1 -os windows|linux|mac
```

If you ran docker container then just run:
```
go run main.go -option 2
```