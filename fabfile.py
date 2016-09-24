from fabric.api import *  # NOQA


def b():
    local('go build -o build/storm_$(uname -s) main/main.go')
