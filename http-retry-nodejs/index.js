const express = require("express");
const fs = require("fs");
const morgan = require("morgan");

const app = express();
const port = 5000;

// Create a write stream (in append mode) for logging
const logStream = fs.createWriteStream("error.log", { flags: "a" });

// Setup morgan to log requests to both the console and the log file
app.use(morgan("combined", { stream: logStream }));
app.use(morgan("dev"));

app.use("/500", (req, res) => {
  // Log the request manually (optional, as morgan already logs it)
  console.log(`Request received: ${req.method} ${req.path} ${req.ip}`);

  // Return a 500 status code
  res.status(500).send("Internal Server Error");
});

app.listen(port, () => {
  console.log(`App listening at http://localhost:${port}`);
});
