const fs = require('fs');
const path = require('path');

const viewsDir = path.join(__dirname, 'src', 'views');

const replacements = [
  { search: /text-white/g, replace: 'text-slate-900' },
  { search: /text-gray-100/g, replace: 'text-slate-900' },
  { search: /text-gray-200/g, replace: 'text-slate-800' },
  { search: /text-gray-300/g, replace: 'text-slate-700' },
  { search: /text-gray-400/g, replace: 'text-slate-500' },
  { search: /text-gray-500/g, replace: 'text-slate-500' },
  
  { search: /bg-gray-950/g, replace: 'bg-slate-50' },
  { search: /bg-gray-900\/50/g, replace: 'bg-white' },
  { search: /bg-gray-900/g, replace: 'bg-white' },
  { search: /bg-gray-800\/50/g, replace: 'bg-slate-50' },
  { search: /bg-gray-800\/40/g, replace: 'bg-slate-50' },
  { search: /bg-gray-800/g, replace: 'bg-slate-100' },
  { search: /bg-gray-700/g, replace: 'bg-slate-200' },
  
  { search: /border-gray-800/g, replace: 'border-slate-200' },
  { search: /border-gray-700\/50/g, replace: 'border-slate-200' },
  { search: /border-gray-700/g, replace: 'border-slate-200' },
  { search: /border-gray-600/g, replace: 'border-slate-300' },
  
  // Strip out glow effects
  { search: /shadow-\[0_0_15px_rgba[^\]]+\]/g, replace: 'shadow-sm' },
  { search: /backdrop-blur-md/g, replace: '' },
  { search: /backdrop-blur-xl/g, replace: '' },
  
  // Fix specific neon colors in raw tables / SVGs
  { search: /text-indigo-400/g, replace: 'text-blue-600' },
  { search: /bg-indigo-500\/10/g, replace: 'bg-blue-50' },
  { search: /border-indigo-500\/20/g, replace: 'border-blue-100' },
  { search: /bg-emerald-500\/10/g, replace: 'bg-emerald-50' }
];

function processDirectory(dir) {
  const files = fs.readdirSync(dir);
  for (const file of files) {
    const fullPath = path.join(dir, file);
    const stat = fs.statSync(fullPath);
    if (stat.isDirectory()) {
      processDirectory(fullPath);
    } else if (file.endsWith('.svelte')) {
      let content = fs.readFileSync(fullPath, 'utf8');
      
      for (const rule of replacements) {
        content = content.replace(rule.search, rule.replace);
      }
      
      fs.writeFileSync(fullPath, content, 'utf8');
      console.log(`Refactored ${file}`);
    }
  }
}

processDirectory(viewsDir);
console.log("Refactoring complete.");
