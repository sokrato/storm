from fabric.api import *  # NOQA


def b():
    local('go build -o build/storm_linux main/main.go')
