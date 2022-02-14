const fs = require("fs");

const parsedFileEN = JSON.parse(fs.readFileSync("./lang/en.json"));
const parsedFileVI = JSON.parse(fs.readFileSync("./lang/vi.json"));

let keysAdded = 0;

Object.keys(parsedFileEN).forEach((element) => {
  if (!parsedFileVI.hasOwnProperty(element)) {
    parsedFileVI[element] = parsedFileEN[element];
    keysAdded++;
  }
});

fs.writeFileSync("./lang/vi.json", JSON.stringify(parsedFileVI, null, 2));

console.log(`Added ${keysAdded} key(s)`);
