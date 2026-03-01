Write-Host "=> Fetching dependencies..."
go mod tidy

Write-Host "=> Building FeatherTrailMD..."
make build

Write-Host "=> Installing to GOPATH..."
make install

Write-Host "=> Installation complete! Run 'ft help' to get started."
