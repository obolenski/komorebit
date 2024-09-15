function New-IcoIcon {
    param(
        [Parameter(Mandatory = $true)]
        [string]$label,
        [Parameter(Mandatory = $true)]
        [string]$icoFile
    )
    Write-Host "Converting label $label to $icoFile"
    magick `
        -background none `
        -fill white `
        -font "JetBrainsMonoNerdFont-Regular.ttf" `
        -pointsize 200 `
        -size 256x256 `
        -gravity center `
        label:"$label" `
        -extent 256x256 `
        $icoFile
}

function New-GoIcon {
    param(
        [Parameter(Mandatory = $true)]
        [string]$fileName
    )
    $icoFile = Join-Path $icoFolder "$fileName.ico"
    Write-Host "Generating $fileName.go from $icoFile"
    
    $outputFile = Join-Path $goOutFolder "$fileName.go"
    
    "//+build windows" | Set-Content $outputFile
    "" | Add-Content $outputFile
    
    $fileNameTitleCase = $fileName.Substring(0, 1).ToUpper() + $fileName.Substring(1)
    $bytes = [System.IO.File]::ReadAllBytes($icoFile)
    $result = $bytes | 2goarray $fileNameTitleCase generated
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Unable to create output file"
        exit
    }
    
    $result | Add-Content $outputFile
    Write-Host "Done"
}

$icoFolder = Join-Path $PSScriptRoot "assets\icons\ico"
$goOutFolder = Join-Path $PSScriptRoot "internal\icons\generated"

if (Test-Path $icoFolder) {
    Remove-Item $icoFolder -Recurse
}

if (Test-Path $goOutFolder) {
    Remove-Item $goOutFolder -Recurse
}

New-Item -ItemType Directory -Path $icoFolder | Out-Null
New-Item -ItemType Directory -Path $goOutFolder | Out-Null

[PsCustomObject[]]$icons = @(
    [PsCustomObject]@{ display = "1"; name = "One" },
    [PsCustomObject]@{ display = "2"; name = "Two" },
    [PsCustomObject]@{ display = "3"; name = "Three" },
    [PsCustomObject]@{ display = "4"; name = "Four" },
    [PsCustomObject]@{ display = "5"; name = "Five" },
    [PsCustomObject]@{ display = "6"; name = "Six" },
    [PsCustomObject]@{ display = "7"; name = "Seven" },
    [PsCustomObject]@{ display = "8"; name = "Eight" },
    [PsCustomObject]@{ display = "9"; name = "Nine" },
    [PsCustomObject]@{ display = "10"; name = "Ten" },
    [PsCustomObject]@{ display = "`u{db80}`u{dfe4}"; name = "Pause" },
    [PsCustomObject]@{ display = ":("; name = "Sad" },
    [PsCustomObject]@{ display = "`u{db81}`u{df25}"; name = "Tilde" }
)

foreach ($icon in $icons) {
    Write-Host "Generating icon for $icon"
    $label = $icon.display
    $fileName = $icon.name
    $icoPath = Join-Path $icoFolder "$fileName.ico"
    New-IcoIcon $label $icoPath
    New-GoIcon $icon.name
}
    
