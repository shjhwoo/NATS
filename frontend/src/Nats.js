import { connect } from "nats";

const urlObject = {servers: "nats://localhost:4222"};

export async function connectNats() {
    const nc = await connect(urlObject);
    console.log("Connected to " + nc.getServer());
    return nc
}