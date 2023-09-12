function Get-MikeScript-CleanUpEmptyFolders {
    [CmdletBinding(SupportsShouldProcess=$true, ConfirmImpact='Medium')]
    param (
        [Parameter(Mandatory=$true)]
        [string]$Path
    )

    process {
        # Get all directories recursively
        $directories = Get-ChildItem -Path $Path -Directory -Recurse

        # Loop over directories and remove if empty
        foreach ($dir in $directories) {
            $items = Get-ChildItem -Path $dir.FullName
            if (!$items) {
                if ($PSCmdlet.ShouldProcess($dir.FullName, "Removing empty directory")) {
                    Remove-Item -Path $dir.FullName
                }
            }
        }
    }
}

# Usage with WhatIf
Get-MikeScript-CleanUpEmptyFolders -Path "C:\path\to\directory" -WhatIf

# Actual usage with confirmation
Get-MikeScript-CleanUpEmptyFolders -Path "C:\path\to\directory" -Confirm
