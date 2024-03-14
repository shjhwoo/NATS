// ES6 modules can be natively used by set the script type
// to "module". Now we can use native imports.
import {
connect,
StringCodec,
} from "https://cdn.jsdelivr.net/npm/nats.ws@1.10.0/esm/nats.js";

// Initialize a string codec for encoding and decoding message data.
// This is because NATS message data are just byte arrays, so proper
// encoding and decoding needs to be performed when using the working
// with the message data.
// See also JSONCodec.c
let sc = new StringCodec();
let natClient


$(function () {
    initialize()

    $(document).on("click", "button", function (e) {
        e.preventDefault()
        const message = $("textarea").val()
        console.log(message)
        sendMessage(message)
         $("textarea").empty()
    })

    $(document).on("keyup", "textarea", function (e) {
        if(e.keyCode && e.keyCode == 13){
            const message = $(this).val()
            console.log(message)
            sendMessage(message)
             $("textarea").empty()
        }
    })
});

async function initialize(){
// Establish a connection to the NATS demo server. This uses the
// native WebSocket support built into the NATS server.
natClient = await connect({
servers: ["ws://localhost:8080"],
});

console.log("client:", natClient)

// Subscribe to the "echo" subject and define a message
// handler that will respond to the requester.
const sub = natClient.subscribe(">");
const handle = (msg) => {
console.log(msg); 
console.log(`Received a request: ${sc.decode(msg.data)}`);
const message = sc.decode(msg.data)
$("ul").append(`<li>${message}</li>`)
msg.respond(msg.data);
}

// Wait to receive messages from the subscription and handle them
// asynchronously..
(async () => {
for await (const msg of sub) handle(msg)
})();

// Now we can send a couple requests to that subject. Note how we
// are encoding the string data on request and decoding the reply
// message data.
let rep = await natClient.request("echo", sc.encode("Hello!"));
console.log(`Received a reply: ${sc.decode(rep.data)}`);

rep = await natClient.request("echo", sc.encode("World!"));
console.log(`Received a reply: ${sc.decode(rep.data)}`);

// Finally drain the connection which will handle any outstanding
// messages before closing the connection.
//natClient.drain(); 쓰지마
}

function sendMessage(message){
    natClient.publish("echo", sc.encode( message));
}