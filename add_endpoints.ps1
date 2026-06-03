$path = ".\gateway\krakend.json"
$json = Get-Content $path -Raw | ConvertFrom-Json

$endpoints = @(
  'domains', 'inbound-routes', 'outbound-routes', 'voicemails',
  'video-rooms', 'cameras', 'chats', 'sms', 'credentials',
  'webhooks', 'throttling', 'firewall', 'fraud'
)

foreach ($ep in $endpoints) {
    # GET
    $getEndpoint = @{
        endpoint = "/api/v1/$ep"
        method = "GET"
        input_headers = @("Authorization", "Content-Type", "Accept")
        backend = @(
            @{
                url_pattern = "/api/$ep"
                host = @("http://pbx_backend:8081")
            }
        )
    }
    $json.endpoints += $getEndpoint

    # POST
    $postEndpoint = @{
        endpoint = "/api/v1/$ep"
        method = "POST"
        input_headers = @("Authorization", "Content-Type", "Accept")
        backend = @(
            @{
                url_pattern = "/api/$ep"
                method = "POST"
                host = @("http://pbx_backend:8081")
            }
        )
    }
    $json.endpoints += $postEndpoint
}

$json | ConvertTo-Json -Depth 10 | Set-Content $path
Write-Host "krakend.json updated."
