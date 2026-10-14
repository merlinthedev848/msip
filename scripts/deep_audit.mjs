import fs from 'fs';
import path from 'path';

function walk(dir) {
  let results = [];
  const list = fs.readdirSync(dir);
  for (const file of list) {
    if (file === 'vendor' || file === 'node_modules' || file === '.git' || file === 'dist') continue;
    const full = path.join(dir, file);
    const stat = fs.statSync(full);
    if (stat.isDirectory()) {
      results = results.concat(walk(full));
    } else {
      results.push(full);
    }
  }
  return results;
}

// 1. Backend routes
const backendFiles = walk('backend').filter(f => f.endsWith('.go'));
const backendRoutes = new Map(); // "METHOD PATH" -> file
const reGin = /(r|api|auth|v1|public|adminGroup)\.(GET|POST|PUT|DELETE)\s*\(\s*["']([^"']+)["']/g;
for (const file of backendFiles) {
  const content = fs.readFileSync(file, 'utf-8');
  let match;
  while ((match = reGin.exec(content)) !== null) {
    const grp = match[1];
    const method = match[2];
    let p = match[3];
    if (!p.startsWith('/')) p = '/' + p;
    if (grp === 'auth') p = '/auth' + p;
    backendRoutes.set(`${method} ${p}`, file);
  }
}

// 2. KrakenD endpoints
const krakendContent = JSON.parse(fs.readFileSync('gateway/krakend.json', 'utf-8'));
const krakendEndpoints = new Map(); // "METHOD ENDPOINT" -> backend info
for (const ep of krakendContent.endpoints || []) {
  krakendEndpoints.set(`${ep.method} ${ep.endpoint}`, ep);
}

// 3. Frontend API calls
const frontendFiles = walk('frontend/src').filter(f => /\.(svelte|ts|js)$/.test(f));
const frontendCalls = [];
const reApi = /api\.(get|post|put|del)\s*(?:<[^>]+>)?\s*\(\s*[`'"]([^`'"]+)[`'"]/g;
const reApiTemplate = /api\.(get|post|put|del)\s*(?:<[^>]+>)?\s*\(\s*`([^`]+)`/g;

for (const file of frontendFiles) {
  const content = fs.readFileSync(file, 'utf-8');
  let match;
  while ((match = reApi.exec(content)) !== null) {
    const method = match[1] === 'del' ? 'DELETE' : match[1].toUpperCase();
    let endpoint = match[2].split('?')[0];
    if (!endpoint.startsWith('/')) endpoint = '/' + endpoint;
    frontendCalls.push({ file, method, endpoint, raw: match[0] });
  }
  while ((match = reApiTemplate.exec(content)) !== null) {
    const method = match[1] === 'del' ? 'DELETE' : match[1].toUpperCase();
    let template = match[2].split('?')[0];
    // Replace ${...} with :param or placeholder
    let normalized = template.replace(/\$\{[^}]+\}/g, ':id');
    if (!normalized.startsWith('/')) normalized = '/' + normalized;
    frontendCalls.push({ file, method, endpoint: normalized, raw: match[0] });
  }
}

console.log(`\n=== FRONTEND CALLS TOTAL: ${frontendCalls.length} ===`);

// Check each frontend call against KrakenD and Backend
const missingInGateway = [];
const missingInBackend = [];

for (const call of frontendCalls) {
  // Check gateway: KrakenD has /api/v1 prefix
  // In KrakenD, params are {param} rather than :param
  const krakendPattern = '/api/v1' + call.endpoint.replace(/:([a-zA-Z0-9_]+)/g, '{$1}');
  
  let foundGateway = false;
  for (const [key, ep] of krakendEndpoints.entries()) {
    const [method, path] = key.split(' ');
    if (method !== call.method) continue;
    // Check match
    const regexPattern = '^' + path.replace(/{[^}]+}/g, '[^/]+') + '$';
    if (new RegExp(regexPattern).test('/api/v1' + call.endpoint)) {
      foundGateway = true;
      break;
    }
  }
  if (!foundGateway) {
    missingInGateway.push(call);
  }

  // Check Backend
  let foundBackend = false;
  for (const [key, file] of backendRoutes.entries()) {
    const [method, path] = key.split(' ');
    if (method !== call.method) continue;
    const regexPattern = '^' + path.replace(/:[a-zA-Z0-9_]+/g, '[^/]+') + '$';
    if (new RegExp(regexPattern).test(call.endpoint)) {
      foundBackend = true;
      break;
    }
  }
  if (!foundBackend) {
    missingInBackend.push(call);
  }
}

console.log(`\n=== MISSING IN GATEWAY (${missingInGateway.length}) ===`);
for (const m of missingInGateway) {
  console.log(`${m.method.padEnd(7)} ${m.endpoint.padEnd(35)} (in ${m.file})`);
}

console.log(`\n=== MISSING IN BACKEND (${missingInBackend.length}) ===`);
for (const m of missingInBackend) {
  console.log(`${m.method.padEnd(7)} ${m.endpoint.padEnd(35)} (in ${m.file})`);
}

// 4. Check all backend routes against KrakenD
const backendOnlyMissingInKrakenD = [];
for (const [key, file] of backendRoutes.entries()) {
  if (file.endsWith('_test.go')) continue;
  const [method, path] = key.split(' ');
  if (['/health', '/api/health', '/fs', '/fs/curl', '/cdr'].includes(path)) continue;
  if (path.startsWith('/ai/stream') || path.startsWith('/wallboard/stream')) continue; // ws handled in nginx

  let kPath = '/api/v1' + path.replace(/:([a-zA-Z0-9_]+)/g, '{$1}');
  if (path === '/login') kPath = '/api/v1/auth/login';

  let found = false;
  for (const [gKey, ep] of krakendEndpoints.entries()) {
    const [gMethod, gPath] = gKey.split(' ');
    if (gMethod !== method) continue;
    const regexPattern = '^' + gPath.replace(/{[^}]+}/g, '[^/]+') + '$';
    if (new RegExp(regexPattern).test(kPath)) {
      found = true;
      break;
    }
  }
  if (!found) {
    backendOnlyMissingInKrakenD.push({ method, path, kPath, file });
  }
}

console.log(`\n=== BACKEND ROUTES MISSING IN KRAKEND (${backendOnlyMissingInKrakenD.length}) ===`);
for (const b of backendOnlyMissingInKrakenD) {
  console.log(`${b.method.padEnd(7)} ${b.path.padEnd(35)} -> ${b.kPath} (${b.file})`);
}

