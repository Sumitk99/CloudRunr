#!/bin/bash
echo "Git url : $GIT_REPOSITORY_URL"

# Exit on any error
set -e
# Ensure GIT_REPOSITORY_URL is set
if [ -z "$GIT_REPOSITORY_URL" ]; then
  echo "GIT_REPOSITORY_URL environment variable not set"
  exit 1
fi

# Clone the repository
git clone "$GIT_REPOSITORY_URL" /home/app/output

# Run the pre-built Go binary
/home/app/app

exit


