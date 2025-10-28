import csvParser from "csv-parser";
import * as fs from "fs";

let results = [];

const args = process.argv.slice(2);

const csvFilePath = args[0];

const stream = fs.createReadStream(csvFilePath).pipe(csvParser());

function errorCallback(err) {
  if (err) {
    console.error("Error writing JSON file:", err);
  } else {
    console.log("JSON file created successfully");
  }
}

function existInResult(target) {
  for (const item of results) {
    if (item === target) {
      return true;
    }
  }
}

// https://stackoverflow.com/questions/32589197/how-can-i-capitalize-the-first-letter-of-each-word-in-a-string-using-javascript
// capitalize every word
function titleCase(str) {
  var splitStr = str.toLowerCase().split(" ");
  for (var i = 0; i < splitStr.length; i++) {
    if (splitStr[i] === "and" || splitStr[i] === "of") continue;

    splitStr[i] =
      splitStr[i].charAt(0).toUpperCase() + splitStr[i].substring(1);
  }
  return splitStr.join(" ");
}

if (csvFilePath === "majors.csv") {
  stream
    .on("data", (data) => {
      const item = titleCase(Object.values(data)[1]);

      if (!existInResult(item)) {
        // results.push({
        //   id: item,
        //   name: item,
        // });
        results.push(item);
      }
    })
    .on("end", () => {
      results = results;
      fs.writeFile("majors.json", JSON.stringify(results), errorCallback);
    });
}

if (csvFilePath === "schools.csv") {
  stream
    .on("data", (data) => {
      const item = titleCase(Object.values(data)[0]);

      if (!existInResult(item)) {
        // results.push({
        //   id: item,
        //   name: item,
        // });
        results.push(item);
      }
    })
    .on("end", () => {
      fs.writeFile("schools.json", JSON.stringify(results), errorCallback);
    });
}

if (csvFilePath === "countries.csv") {
  stream
    .on("data", (data) => {
      const item = titleCase(Object.values(data)[0]);

      if (!existInResult(item)) {
        results.push(item);
      }
    })
    .on("end", () => {
      fs.writeFile("countries.json", JSON.stringify(results), errorCallback);
    });
}
