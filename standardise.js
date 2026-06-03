const fs = require('fs');
const path = require('path');

const dir = path.join(__dirname, 'frontend/src/views');
const files = fs.readdirSync(dir).filter(f => f.endsWith('.svelte'));

const targetClass = 'space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-700 ease-out max-w-7xl mx-auto w-full';

for (const file of files) {
    const filePath = path.join(dir, file);
    let content = fs.readFileSync(filePath, 'utf8');
    
    // Find the first <div class="...animate-in...">
    const regex = /<div class="[^"]*animate-in[^"]*">/;
    if (regex.test(content)) {
        content = content.replace(regex, `<div class="${targetClass}">`);
        fs.writeFileSync(filePath, content, 'utf8');
        console.log(`Updated ${file}`);
    } else {
        // Fallback for files that might not have animate-in yet (like the older ones)
        // Look for the first <div class="..."> after the script tag
        const scriptEndIndex = content.indexOf('</script>');
        if (scriptEndIndex !== -1) {
            const firstDivRegex = /<div class="([^"]*)">/;
            const match = firstDivRegex.exec(content.substring(scriptEndIndex));
            if (match) {
                const globalIndex = scriptEndIndex + match.index;
                // Only replace if it looks like a wrapper
                if (match[1].includes('space-y-') || match[1].includes('p-') || match[1].includes('max-w-')) {
                    content = content.substring(0, globalIndex) + `<div class="${targetClass}">` + content.substring(globalIndex + match[0].length);
                    fs.writeFileSync(filePath, content, 'utf8');
                    console.log(`Updated ${file} (fallback)`);
                }
            }
        }
    }
}
