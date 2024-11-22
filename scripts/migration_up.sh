#!/bin/bash

migrate -database $env:POSTGRESQL_URL -path db/migrations up