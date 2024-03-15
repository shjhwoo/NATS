// // ES6 modules can be natively used by set the script type
// // to "module". Now we can use native imports.
// import {
//     AckPolicy,
//     connect,
//     StringCodec,
//     } from "https://cdn.jsdelivr.net/npm/nats.ws@1.10.0/esm/nats.js";

// // Initialize a string codec for encoding and decoding message data.
// // This is because NATS message data are just byte arrays, so proper
// // encoding and decoding needs to be performed when using the working
// // with the message data.
// // See also JSONCodec.c
// let sc = new StringCodec();
// let natClient
// let jetStreamManager
// const streamName = "myChatRoom"
// const jetStreamConfig = {

//     name: streamName,
//     subjects: ["msg.*.*"],
// }

// let jetStream
// let jetStreamConsumer

// let myId = Math.floor(Math.random() * 1000)

// $(function () {
//     initialize()

//     $(document).on("click", "button", function (e) {
//         e.preventDefault()
//         const message = $("textarea").val()
//         console.log(message)
//         sendMessage(message)
//          $("textarea").empty()
//     })

//     $(document).on("keyup", "textarea", function (e) {
//         if(e.keyCode && e.keyCode == 13){
//             const message = $(this).val()
//             console.log(message)
//             sendMessage(message)
//              $("textarea").empty()
//         }
//     })
// });

// async function initialize(){
//     try {
//     natClient = await connect({
//     servers: ["ws://localhost:8080"],
//     });

//     //제트스트림 클라이언트를 만든다
//     jetStreamManager = await natClient.jetstreamManager();
//     await jetStreamManager.streams.add(jetStreamConfig);
//     jetStream = natClient.jetstream();
//     console.log("created jetStream.",jetStream)

//     await jetStream.publish(`msg.${myId}.connected`);
//     console.log("published connected message");

//     let info = await jetStreamManager.streams.info(jetStreamConfig.name);
//     console.log("stream info", info);

//     } catch (err) {
//         console.log(err.message);
//         if (err.message == 'stream name already in use with a different configuration'){
//             jetStream = natClient.jetstream();
//             await jetStream.publish(`msg.${myId}.connected`);
//             console.log("published connected message");
        
//             let info = await jetStreamManager.streams.info(jetStreamConfig.name);
//             console.log("stream info", info);
//         }
//     }finally{
//         console.log('메세지 듣는 로직 시작')
//         const jetStreamConsumer = await jetStreamManager.consumers.add(streamName,{ ack_policy: AckPolicy.Explicit, durable_name: "A", })
//         console.log("created consumer: ", jetStreamConsumer)

//         console.log("after adding consumer:", jetStream)

//         const c = await jetStream.consumers.get(streamName, jetStreamConsumer.name);
//         console.log("received: ", c)

//         // Finally drain the connection which will handle any outstanding
//         // messages before closing the connection.
//         //natClient.drain(); 쓰지마

//     }
// }

// async function sendMessage(message){
//     await jetStream.publish(`msg.${myId}.${message.trimEnd()}`);
// }