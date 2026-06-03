const fs = require('fs');
const path = './gateway/krakend.json';

let data = JSON.parse(fs.readFileSync(path, 'utf8'));

const endpoints = [
  'domains', 'inbound-routes', 'outbound-routes', 'voicemails',
  'video-rooms', 'cameras', 'chats', 'sms', 'credentials',
  'webhooks', 'throttling', 'firewall', 'fraud'
];

for (const ep of endpoints) {
  // GET
  data.endpoints.push({
    endpoint: `/api/v1/${ep}`,
    method: 'GET',
    input_headers: ["Authorization", "Content-Type", "Accept"],
    backend: [{
      url_pattern: `/api/${ep}`,
      host: ["http://pbx_backend:8081"]
    }]
  });
  // POST
  data.endpoints.push({
    endpoint: `/api/v1/${ep}`,
    method: 'POST',
    input_headers: ["Authorization", "Content-Type", "Accept"],
    backend: [{
      url_pattern: `/api/${ep}`,
      method: 'POST',
      host: ["http://pbx_backend:8081"]
    }]
  });
}

fs.writeFileSync(path, JSON.stringify(data, null, 2));
console.log("krakend.json updated.");
