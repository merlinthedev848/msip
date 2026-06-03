import json

path = './gateway/krakend.json'

with open(path, 'r') as f:
    data = json.load(f)

endpoints = [
  'domains', 'inbound-routes', 'outbound-routes', 'voicemails',
  'video-rooms', 'cameras', 'chats', 'sms', 'credentials',
  'webhooks', 'throttling', 'firewall', 'fraud'
]

for ep in endpoints:
    data['endpoints'].append({
        'endpoint': f'/api/v1/{ep}',
        'method': 'GET',
        'input_headers': ["Authorization", "Content-Type", "Accept"],
        'backend': [{
            'url_pattern': f'/api/{ep}',
            'host': ["http://pbx_backend:8081"]
        }]
    })
    data['endpoints'].append({
        'endpoint': f'/api/v1/{ep}',
        'method': 'POST',
        'input_headers': ["Authorization", "Content-Type", "Accept"],
        'backend': [{
            'url_pattern': f'/api/{ep}',
            'method': 'POST',
            'host': ["http://pbx_backend:8081"]
        }]
    })

with open(path, 'w') as f:
    json.dump(data, f, indent=2)

print("krakend.json updated.")
