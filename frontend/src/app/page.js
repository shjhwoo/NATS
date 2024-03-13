import styles from "./page.module.css";
import { connect } from "nats";
// ("use client");
import { MyMsg } from "./message/mymsg.js";
import { YouMsg } from "./message/youmsg.js";

export default async function Home() {
  const nc = await connectNATS();

  console.log(nc.getServer(), "연결된 서버 주소입니다");

  return (
    <main className={styles.main}>
      <div className={styles.chatbox}>
        <div className={styles.messageList}>
          <ul></ul>
        </div>
        <div className={styles.enterChatArea}>
          <textarea
            className={styles.writetext}
            placeholder="write message here..."
          ></textarea>
          <button className={styles.sendBtn} onClick={sendMessage}>
            Send
          </button>
        </div>
      </div>
    </main>
  );
}

async function connectNATS() {
  const servers = "nats://localhost:4222";
  const nc = await connect({
    servers: servers.split(","),
  });

  console.log(nc);

  return nc;
}

function sendMessage() {}
