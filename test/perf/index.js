const http = require('http')
const fs = require('fs');
const { promisify } = require('util')
const sleep = promisify(setTimeout)

var options = {
  'method': 'POST',
  'hostname': '0.0.0.0',
  'port': 8090,
  'path': '/collect',
  'headers': {
    'Content-Type': 'application/json'
  },
  'maxRedirects': 20
};

var req = http.request(options, function (res) {
  var chunks = [];

  res.on("data", function (chunk) {
    chunks.push(chunk);
  });

  res.on("end", function (chunk) {
    var body = Buffer.concat(chunks);
    console.log(body.toString());
  });

  res.on("error", function (error) {
    console.error(error);
  });
});


async function main() {
    maxCnt = 10000000
    perCnt = 500000
    pdone = false
    done = false
    count = 0
    to = 1000/perCnt
    while(!done) {
        const postData = JSON.stringify({
          "Merchant": "macdonlds",
          "Count": count,
          "Sum": 169.89,
          "Project": "shop",
          "SendCurrency": "ru",
          "Method": "get",
          "Name": "some name",
          "CardName": "visa",
          "CardNumber": "12555667788",
          "ExpireDate": "11/22",
          "SecurityCode": "123",
          "ReceiveCurrency": "usd",
          "Rate": 65,
          "TransactionStatus": "done",
          "TransactionTime": "2012-04-23T18:25:43Z"
        });

        postData["Count"] = count
        await req.write(postData);
        //await req.end();
        console.log(count)

        await sleep(to)
        count++
        if(count == maxCnt) { done = true }
    }
    process.exit(0)
}

main()
