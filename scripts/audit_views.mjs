import fs from 'fs';
import path from 'path';

function walk(dir) {
  let results = [];
  const list = fs.readdirSync(dir);
  for (const file of list) {
    if (file === 'node_modules' || file === 'dist' || file === '.git') continue;
    const full = path.join(dir, file);
    const stat = fs.statSync(full);
    if (stat.isDirectory()) results = results.concat(walk(full));
    else results.push(full);
  }
  return results;
}

const views = walk('frontend/src').filter(f => f.endsWith('.svelte'));
console.log('--- AUDITING SVELTE FILES (' + views.length + ') ---');

let issueCount = 0;

for (const v of views) {
  const content = fs.readFileSync(v, 'utf-8');
  
  // 1. TODO or placeholder checks
  const todos = content.match(/\/\/\s*TODO[^\n]*/gi);
  if (todos) {
    console.log(`[TODO Found] in ${v}:`, todos);
    issueCount += todos.length;
  }

  // 2. alert() or confirm() placeholders
  const alerts = content.match(/alert\([^)]+\)/g);
  if (alerts) {
    console.log(`[Alert Placeholder] in ${v}:`, alerts);
    issueCount += alerts.length;
  }

  // 3. Buttons without handlers
  const btnMatches = content.match(/<button[\s\S]*?>/g) || [];
  for (const b of btnMatches) {
    const clean = b.replace(/\s+/g, ' ');
    if (!clean.includes('on:click') && !clean.includes('type="submit"') && !clean.includes("type='submit'") && !clean.includes('on:mousedown') && !clean.includes('on:keydown')) {
      // Check if button is inside a form with on:submit
      console.log(`[Button without Handler] in ${v}:`, clean.slice(0, 100));
      issueCount++;
    }
  }

  // 4. Broken anchor tags (href="#" or href="javascript:")
  const hrefs = content.match(/href=["'](#[^"']*|javascript:[^"']*)["']/g);
  if (hrefs) {
    console.log(`[Empty/Dead Link] in ${v}:`, hrefs);
    issueCount += hrefs.length;
  }
}

console.log(`--- TOTAL POTENTIAL ISSUES FOUND: ${issueCount} ---`);
