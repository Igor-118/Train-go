#!/bin/sh


goose -dir=/app/migrations up

if [ $? -ne 0 ]; then
    echo "Migration failed!"
    exit 1
fi


exec ./app