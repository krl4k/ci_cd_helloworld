#!/bin/bash

# Install GitLab Runner
curl -L "https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.deb.sh" | sudo bash
sudo apt install gitlab-runner

# Register the runner with your GitLab project
sudo gitlab-runner register

# Start and enable the runner
sudo systemctl start gitlab-runner
sudo systemctl enable gitlab-runner 
