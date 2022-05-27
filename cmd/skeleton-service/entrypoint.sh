#!/bin/bash

set -euo pipefail

/app/migrate -source=file:///app/config/db/migrations/ -database=configdb:///app/config/config.yaml up
/app/skeleton-service