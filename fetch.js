import { writeFileSync } from 'fs';

const BANGS_URL = "https://duckduckgo.com/bang.js";
const OUTPUT_FILE = "bangs.json";

const response = await fetch(BANGS_URL);
const jsContent = await response.text();
const match = jsContent.match(/\[[\s\S]*\]/);
const bangs = JSON.parse(match[0]);
const filtered = bangs.map(bang => ({
    t: bang.t,
    u: bang.u
}));
writeFileSync(OUTPUT_FILE, JSON.stringify(filtered, null, 2));
