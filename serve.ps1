param(
  [int]$Port = 8000,
  [string]$Dir = ".",
  [string]$Prefix = "/EpicFreeGamesList",
  [string]$LogFile = "serve.log"
)

Set-StrictMode -Version Latest

function Get-MimeType {
  param([string]$ext)
  switch ($ext.ToLower()) {
    '.html' { 'text/html' }
    '.htm'  { 'text/html' }
    '.js'   { 'application/javascript' }
    '.css'  { 'text/css' }
    '.json' { 'application/json' }
    '.png'  { 'image/png' }
    '.jpg'  { 'image/jpeg' }
    '.jpeg' { 'image/jpeg' }
    '.svg'  { 'image/svg+xml' }
    '.ico'  { 'image/x-icon' }
    default { 'application/octet-stream' }
  }
}

Push-Location -LiteralPath $Dir
try {
  $listener = New-Object System.Net.HttpListener
  # Bind to localhost to avoid requiring URL ACL reservations
  $prefixUrl = "http://localhost:$Port/"
  $listener.Prefixes.Add($prefixUrl)
  $listener.Start()
  $logDirectory = Split-Path -Path $LogFile -Parent
  if (-not [string]::IsNullOrEmpty($logDirectory) -and -not (Test-Path -LiteralPath $logDirectory)) {
    New-Item -ItemType Directory -Path $logDirectory -Force | Out-Null
  }
  function Write-Log {
    param([string]$Message)
    $line = "$(Get-Date -Format o) $Message"
    Add-Content -Path $LogFile -Value $line
  }
  Write-Log "Serving '$((Get-Location).ProviderPath)' on http://0.0.0.0:$Port (prefix=$Prefix)"
  try {
    Start-Process "http://localhost:$Port$Prefix/" -ErrorAction SilentlyContinue
  } catch {
    Write-Log "Browser launch failed: $($_.Exception.Message)"
  }

  while ($listener.IsListening) {
    $context = $listener.GetContext()
    try {
      $req = $context.Request
      $res = $context.Response

      $path = $req.Url.AbsolutePath
      if ($path.StartsWith($Prefix)) { $path = $path.Substring($Prefix.Length) }
      if ([string]::IsNullOrEmpty($path) -or $path -eq '/') { $path = '/index.html' }

      $rel = $path.TrimStart('/')
      $base = (Get-Item -LiteralPath .).FullName
      $full = [System.IO.Path]::GetFullPath((Join-Path $base $rel))

      if (-not $full.StartsWith($base, [System.StringComparison]::OrdinalIgnoreCase)) {
        $res.StatusCode = 403
        $buf = [System.Text.Encoding]::UTF8.GetBytes('403 Forbidden')
        $res.ContentType = 'text/plain'
        $res.ContentLength64 = $buf.Length
        $res.OutputStream.Write($buf, 0, $buf.Length)
        continue
      }

      if (Test-Path -LiteralPath $full -PathType Leaf) {
        $bytes = [System.IO.File]::ReadAllBytes($full)
        $ext = [System.IO.Path]::GetExtension($full)
        $res.ContentType = Get-MimeType -ext $ext
        # Disable browser caching so changes appear immediately on refresh
        $res.Headers['Cache-Control'] = 'no-cache, no-store, must-revalidate'
        $res.Headers['Pragma'] = 'no-cache'
        $res.Headers['Expires'] = '0'
        $res.Headers['Last-Modified'] = (Get-Item -LiteralPath $full).LastWriteTimeUtc.ToString('R')
        $res.ContentLength64 = $bytes.Length
        $res.StatusCode = 200
        $res.OutputStream.Write($bytes, 0, $bytes.Length)
      }
      else {
        $res.StatusCode = 404
        $msg = '404 Not Found'
        $buf = [System.Text.Encoding]::UTF8.GetBytes($msg)
        $res.ContentType = 'text/plain'
        $res.ContentLength64 = $buf.Length
        $res.OutputStream.Write($buf, 0, $buf.Length)
      }
    } catch {
      Write-Log "Request error: $($_.Exception.Message)"
    } finally {
      try { $context.Response.OutputStream.Close() } catch {}
      try { $context.Response.Close() } catch {}
    }
  }
} finally {
  try { $listener.Stop() } catch {}
  Pop-Location
}
