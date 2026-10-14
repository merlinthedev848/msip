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

const frontendApis = new Set();
const reApi = /api\.(?:get|post|put|del)\s*(?:<[^>]+>)?\s*\(\s*[`'"]([^`'"?]+)/g;

const frontendFiles = walk('frontend/src').filter(f => /\.(svelte|ts|js)$/.test(f));
for (const file of frontendFiles) {
  const content = fs.readFileSync(file, 'utf-8');
  let match;
  while ((match = reApi.exec(content)) !== null) {
    let p = match[1].split('${')[0].replace(/\/+$/, '');
    if (!p.startsWith('/')) p = '/' + p;
    frontendApis.add(p);
  }
}

console.log(`\n=== FRONTEND API CALLS (${frontendApis.size}) ===`);
for (const a of Array.from(frontendApis).sort()) {
  console.log('  ' + a);
}

// KrakenD endpoints
const krakendContent = JSON.parse(fs.readFileSync('gateway/krakend.json', 'utf-8'));
const krakendEndpoints = new Set();
const krakendMap = new Map();
for (const ep of krakendContent.endpoints || []) {
  krakendEndpoints.add(ep.endpoint);
  krakendMap.set(ep.endpoint, ep);
}

console.log(`\n=== KRAKEND ENDPOINTS (${krakendEndpoints.size}) ===`);

// Check which frontend APIs (prefixed with /api/v1) exist in KrakenD
console.log('\n=== WIRING VERIFICATION: FRONTEND -> KRAKEND ===');
let missingInGateway = 0;
for (const a of Array.from(frontendApis).sort()) {
  const fullEndpoint = '/api/v1' + a;
  // Account for wildcard / parameter endpoints like /provisioning/{code} or /routing/ring-groups/{id}
  let found = krakendEndpoints.has(fullEndpoint);
  if (!found) {
    // Check if there is a parameterized endpoint match
    for (const ep of krakendEndpoints) {
      const regexPattern = '^' + ep.replace(/{[^}]+}/g, '[^/]+') + '$';
      if (new RegExp(regexPattern).test(fullEndpoint) || fullEndpoint.startsWith(ep.split('{')[0])) {
        found = true;
        break;
      }
    }
  }
  if (!found) {
    console.log(`  [MISSING IN GATEWAY] ${a} -> ${fullEndpoint}`);
    missingInGateway++;
  } else {
    // console.log(`  [OK] ${a}`);
  }
}
console.log(`Summary: ${missingInGateway} frontend calls missing in KrakenD`);

// Inspect Backend routes
const backendFiles = walk('backend').filter(f => f.endsWith('.go'));
const backendRoutes = new Set();
const reGin = /(?:r|api|auth|v1|public|adminGroup)\.(GET|POST|PUT|DELETE)\s*\(\s*["']([^"']+)["']/g;
for (const file of backendFiles) {
  const content = fs.readFileSync(file, 'utf-8');
  let match;
  while ((match = reGin.exec(content)) !== null) {
    backendRoutes.add(`${match[1]} ${match[2]}`);
  }
}

console.log(`\n=== BACKEND ROUTES (${backendRoutes.size}) ===`);
for (const r of Array.from(backendRoutes).sort()) {
  console.log('  ' + r);
}
