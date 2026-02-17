param(
    [ValidateSet("all", "up", "migrate", "test", "run", "down")]
    [string]$Action = "all"
)

$ErrorActionPreference = "Stop"

$ProjectId = "test-project"
$InstanceId = "test-instance"
$DatabaseId = "test-db"
$SpannerDatabase = "projects/$ProjectId/instances/$InstanceId/databases/$DatabaseId"

function Invoke-Step {
    param(
        [string]$Name,
        [scriptblock]$Script
    )

    Write-Host "==> $Name" -ForegroundColor Cyan
    & $Script
}

function Assert-Command {
    param(
        [string]$Name,
        [string]$InstallHint
    )

    if (Get-Command $Name -ErrorAction SilentlyContinue) {
        return
    }

    throw "$Name is not installed or not in PATH. $InstallHint"
}

function Resolve-GcloudCommand {
    $cmd = Get-Command "gcloud" -ErrorAction SilentlyContinue
    if ($cmd) {
        return "gcloud"
    }

    $candidates = @(
        "$env:LOCALAPPDATA\Google\Cloud SDK\google-cloud-sdk\bin\gcloud.cmd",
        "$env:ProgramFiles\Google\Cloud SDK\google-cloud-sdk\bin\gcloud.cmd",
        "$env:ProgramFiles(x86)\Google\Cloud SDK\google-cloud-sdk\bin\gcloud.cmd"
    )

    foreach ($candidate in $candidates) {
        if (Test-Path $candidate) {
            return $candidate
        }
    }

    throw "gcloud is not installed or not in PATH. Install Google Cloud CLI, for example: winget install Google.CloudSDK"
}

function Start-Emulator {
    Assert-Command -Name "docker" -InstallHint "Install Docker Desktop and ensure 'docker' is available in PATH."
    docker compose up -d
}

function Stop-Emulator {
    docker compose down
}

function Invoke-Migrate {
    $gcloud = Resolve-GcloudCommand
    $env:SPANNER_EMULATOR_HOST = "localhost:9010"

    try {
        & $gcloud config configurations create emulator --no-activate *>$null
    } catch {}
    $global:LASTEXITCODE = 0

    try {
        & $gcloud spanner instances create $InstanceId `
            --config=emulator-config `
            --description="Test Instance" `
            --nodes=1 `
            --project=$ProjectId *>$null
    } catch {}
    $global:LASTEXITCODE = 0

    try {
        & $gcloud spanner databases create $DatabaseId `
            --instance=$InstanceId `
            --project=$ProjectId `
            --ddl-file=migrations/001_initial_schema.sql *>$null
    } catch {}
    $global:LASTEXITCODE = 0

    Write-Host "Migration complete." -ForegroundColor Green
}

function Invoke-Tests {
    Assert-Command -Name "go" -InstallHint "Install Go 1.21+ and ensure 'go' is available in PATH."
    $env:SPANNER_EMULATOR_HOST = "localhost:9010"
    go test ./... -v -count=1
}

function Ensure-PortAvailable {
    param(
        [string]$Port
    )

    $listener = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue | Select-Object -First 1
    if (-not $listener) {
        return
    }

    $ownerPid = $listener.OwningProcess
    $proc = Get-Process -Id $ownerPid -ErrorAction SilentlyContinue
    if (-not $proc) {
        return
    }

    if ($proc.ProcessName -eq "server") {
        Write-Host "Port $Port is already in use by previous server process (PID: $ownerPid). Stopping it..." -ForegroundColor Yellow
        Stop-Process -Id $ownerPid -Force
        Start-Sleep -Milliseconds 300
        return
    }

    throw "Port $Port is in use by process '$($proc.ProcessName)' (PID: $ownerPid). Stop that process or set PORT to another value before running."
}

function Start-Server {
    Assert-Command -Name "go" -InstallHint "Install Go 1.21+ and ensure 'go' is available in PATH."
    go build -o bin/server ./cmd/server
    $port = if ($env:PORT) { $env:PORT } else { "50051" }

    Ensure-PortAvailable -Port $port

    $env:SPANNER_EMULATOR_HOST = "localhost:9010"
    $env:SPANNER_DATABASE = $SpannerDatabase

    if (Test-Path ".\\bin\\server.exe") {
        .\\bin\\server.exe
        return
    }

    .\\bin\\server
}

switch ($Action) {
    "all" {
        Invoke-Step -Name "Start Spanner emulator" -Script ${function:Start-Emulator}
        Invoke-Step -Name "Run migrations" -Script ${function:Invoke-Migrate}
        Invoke-Step -Name "Run tests" -Script ${function:Invoke-Tests}
        Invoke-Step -Name "Start server" -Script ${function:Start-Server}
    }
    "up" {
        Invoke-Step -Name "Start Spanner emulator" -Script ${function:Start-Emulator}
    }
    "migrate" {
        Invoke-Step -Name "Run migrations" -Script ${function:Invoke-Migrate}
    }
    "test" {
        Invoke-Step -Name "Run tests" -Script ${function:Invoke-Tests}
    }
    "run" {
        Invoke-Step -Name "Start server" -Script ${function:Start-Server}
    }
    "down" {
        Invoke-Step -Name "Stop Spanner emulator" -Script ${function:Stop-Emulator}
    }
}