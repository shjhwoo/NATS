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
let myId = Math.floor(Math.random() * 1000)


$(function () {
    initialize()

    $(document).on("click", "button", function (e) {
        e.preventDefault()
        const message = $("textarea").val()
        console.log(message)
        sendMessage(message)
         $("textarea").val('')
    })

    $(document).on("keyup", "textarea", function (e) {
        if(e.keyCode && e.keyCode == 13){
            const message = $(this).val()
            console.log(message)
            sendMessage(message)
             $("textarea").val('')
        }
    })
});

async function initialize(){
    natClient = await connect({
    servers: ["ws://localhost:8080/ws"],
    });

    console.log("Connected to " + natClient.getServer());

    const sub = natClient.subscribe("msg.>");
    (async () => {
        for await (const msg of sub) handle(msg)
    })();

    const handle = (msg) => {
        console.log(sc.decode(msg.data))
        const message = JSON.parse(sc.decode(msg.data))
        if (message.user != myId){
            console.log(message);
            $("ul").append(`<li class="youMsg">${message.text == 'connected'? `${message.user}님이 입장하셨습니다`:`${message.user}: ${message.text}`}</li>`)
        }else{
            if (message.text != 'connected'){
                $("ul").append(`<li class="meMsg">${message.text}</li>`)  
            }
        }
    }

    natClient.publish(`msg.${myId}`, sc.encode(`{"user":${myId},"text":"connected"}`));

    // Finally drain the connection which will handle any outstanding
    // messages before closing the connection.
    //natClient.drain(); 쓰지마
}

function sendMessage(message){
    natClient.publish(`msg.${myId}`, sc.encode(`{"user":"${myId}","text":"${message.trimEnd()}"}`));
}