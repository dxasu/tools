# PS C:> $profile
# C:\Users\Apple\Documents\WindowsPowerShell\Microsoft.VSCode_profile.ps1
# PS C:> . $profile  #使profile生效

Set-Alias grep findstr
Set-Alias vi Notepad
Set-Alias touch New-Item

function listtag {
    param (
        $tagname
    )
    git tag --sort=taggerdate | findstr $tagname
}

function pushtag {
    param (
        $fulltagname
    )
    git tag -a $fulltagname -m "add tag $fulltagname"
    git push --follow-tags
}

function deltag {
    param (
        $fulltagname
    )
    git tag -d $fulltagname
    git push --delete origin $fulltagname
}

function pushall {
    param (
        $notes
    )
    if(! $notes){
        $notes="push all"
    }
    git add .
    git commit -am $notes
    git push
}


# Add Get-AllHistory function for powershell
function GetAll
{
	<#
       .SYNOPSIS
	   # Get-AllHistory
	   Get all history of powershell
	   # Get-AllHistory n
	   Show the last n history records
	   .DESCRIPTION
	   The function add by Kody
    #>
	param (
        $Count
    )
    if($Count){
		$his = Get-Content (Get-PSReadLineOption).HistorySavePath -tail $Count
	}
	else{
		$his = Get-Content (Get-PSReadLineOption).HistorySavePath
	}
    $n = $his.Length
    $out = @()
    for($i=1;$i -le $n;$i++)
    {
        $out = $out + "$i $($his[$i-1])"
    }
    return $out
}
