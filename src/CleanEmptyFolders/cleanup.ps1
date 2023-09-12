function Remove-EmptyDirectories {
    param (
        [Parameter(Mandatory=$true)]
        [string]$Path
    )

    # Get all directories recursively
    $directories = Get-ChildItem -Path $Path -Directory -Recurse

    # Loop over directories and remove if empty
    foreach ($dir in $directories) {
        $items = Get-ChildItem -Path $dir.FullName -Recurse
        if (!$items) {
            Write-Host "Removing empty directory: $($dir.FullName)"
            Remove-Item -Path $dir.FullName
        }
    }
}

# Usage
Remove-EmptyDirectories -Path "C:\path\to\directory"
